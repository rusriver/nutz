Let's document this thing.

Implementation of Circuit Breaker, simple and efficient.

1) Prior art

    Existing alternatives I've found so far:
        https://github.com/mercari/go-circuitbreaker/blob/master/breaker.go
        https://github.com/sony/gobreaker
        https://github.com/go-kit/kit   -- uses sony's one

2) Terminology

    I don't like using analogy from electrical engineering, where the contact can be closed or open,
    because for some people and in some contexts it sounds as the reverse - it effect making it
    ambiguous to understand correctly. For this reson, I intentionally will avoid using words
    "open" and "closed".

    Instead, we will operate in these terms:

        OFF, DENIED
        ON, ALLOWED

    (In electrical engineering, it would correspond as OFF == open contact (no pass), and ON ==
    closed contact (the current can pass), but I promised to not use these terms here).

3) States

    The circuit breaker has two separate states variables, one is more general, another is for more
    detailed monitoring of what exactly happens there (and it also the main one internally):

        CbState
            The general state variable, can have two states - on and off. It affects the allowance
            of operations to execute, if it's on the op is executed, if it's off the operation
            in mostly (but not completely) blocked from execution.

        GState
            The internal state variable, the actual state of the internal state machine, can also be
            used for monitoring purposes, to see what exactly happens inside. It can have three
            states: DENY_RESTART, DENY_WAIT, ALLOW.

    The CbState is set from the GState state machine, by exact matching:

            GState = DENY_RESTART       CbState = OFF_DENIED
            GState = ALLOW              CbState = ON_ALLOWED
            GState = DENY_WAIT          CbState = OFF_DENIED

4) How does it work

    The circuit breaker struct has state int, chReport chan, chReingite chan fields.

    It has three methods:

            Report(r bool) -- Report the result of executed operation.

                    SWITCH CbState {
                        CASE OFF_DENIED {
                            Send the result to the chReport non-blockingly - this means
                            the channel capacity saturation would never be a problem, and
                            losing some result isn't a problem either.
                        }
                        CASE ON_ALLOWED {
                            The same, but send only negative results, because sending positives
                            is a useless waste of resources in this state.
                        }
                    }

            IsOn() (ok bool) -- check if we can execute the next operation

                    SWITCH CbState {
                        CASE ON_ALLOWED {
                            Return true.
                        }
                        CASE OFF_DENIED {
                            Read if there's a signal present on chReignite channel - if yes then
                            return true, else return false. This is used for (re-)start-up.
                        }
                    }

            WantToDo(f func() (success bool)) (result bool, executed bool) -- convenience wrapper, think Try()

                    if cb.IsOn() {
                        result = f()
                        cb.Report(result)
                        return result, true
                    } else {
                        return false, false
                    }

    The constructor launches a parallel G, with context, which reads the chReport, and manages its
    state internally. The result is some events (1 or several) written to chReignite, and changing the state.

    The state is read-only by IsOn, and writable from the G.

    The helper method Report() writes commands to chReport non-blockingly, like this:

        cb.Report(true)
        cb.Report(false)

    If CbState==OFF, both true and false are reported;
    If CbState==ON, only falses are reported;
    This is to avoid excessive and redundant load on channel, when everything works well, under high load.

    Because reporting is non-blocking, the channel capacity contention should never be a bottleneck.
    (But, the contention to use the channel still present, obviously).

    The G does (the core of the algorithm):

        Stop on Context.Done;

        Internal state (GState) is:

            DENY-RESTART (init state):
                CbState = OFF_DENIED
                Write N tokens to ChReignite, every T time, non-blockingly;
                        This is to (re-)start-up the process, after DENY-WAIT;
                        Should N == cap(ChReignite)? yes in typical case;
                        In constructor, we must fill the ChReignite, to make it on on the startup;
                        Typically this must be small, like just few, e.g. 3;
                        (i) Using ticker here makes useful side effect (feature) - randomization of reignition moment;
                Keep counter of fails/successes: +1 on ok, -1 on fail, init at 0;
                If counter is >X, switch to ALLOW;
                If counter is <Y, switch to DENY-WAIT (again);

            ALLOW:
                CbState = ON_ALLOWED
                Keep counter of fails/successes: +1 on ok, -1 on fail, init at 0;
                    Because only fails are reported in this state, we must simulate increment +N/sec, by elapsed time.
                If counter is <Y, switch to DENY-WAIT;

            DENY-WAIT:
                CbState = OFF_DENIED
                Wait some time;
                Switch to DENY-RESTART;

    See also the attached picture.

5) Usage

    User does:

        Init the CB, the G is started;

        if cb.IsOn() == true:
            do the operation
            report true/false vis the cb.Report()
        else:
            skip the operation

    OR, using the convenience wrapper:

        cb.WantToDo(func() bool {
            do the operation
            return if success is true or false
        })

6) Introspection

    If you want to report metrics or log state changes, please use the StateHookFunc callback.
    External direct access to states of the circuit breaker isn't possible.


===

