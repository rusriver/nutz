This is a very small and simple interface wrapper over the zerolog logger,
including the reference implementation of the wrapper. It is production ready,
and used since 2021 privately, until open sourced in 2024.

Goals of making such a wrapper around the zerolog logger:

    1) Ability to make an implementation which will allow logging to extra destinations,
        e.g. to queues or to various logging collection systems, without interfering with the
        user codebase.

            Clarification: In fact, you don't need a wrapper for that, moreover, this wrapper
            won't help with that (if we're talking about the zerolog underneath). That's because
            the buffer with serialized data is private to zerolog, and inaccessible (and imho
            that was a mistake of them). However, you can do so, by intercepting the output
            of zerolog logger. The framing is natural, because it call single Write() for each
            logline.

            Tip: If you make own implementation, you can also use Msg(string) instead of Send()
            internally in Send(), and so pass that string as argument to the zerolog hook.
            Well, that'd be weird though.

            Resume: currently, you can only control whether to do send or not, from the Send()
            hook; the send operation itself is detached in a writer, and you can't link the two,
            but to fork and fix the original zerolog.

    2) Technical ability to add filters, including those activated and managed at runtime.

    3) Ability to set logging ratelims for individual codes/tags or groups of codes/tags,
        without interfering into the main user codebase.

    4) Ability to switch off some points of logging, without removing them from the source code.
        Like a filter, but hardcoded.

    5) Ability to transparently add metrics reporting - by levels and by codes/tags.

    6) Ability to build a histogram in Grafana of what is logging, and so determine
        any parasitic logging, see what can be decreased or throttled or levels corrected.

    7) A histogram in Grafana can also be used as a natural live registry of what is logged,
        so an operator can quickly figure out what codes to look after in the logging system.
        That is, all codes/tags and titles would be accounted for in there, which is very convenient
        and important to have.

    8) In case of erroneous flooding of the logging system, the histogram will help quickly
        diagnose what is the root cause of such a flooding, and assist in quickest possible
        remedy.

    9) Ability to add alerts, both in logging and in metrics systems, for specific codes/tags,
        their groups or sets. It is very important functionality, it allows to track situations of
        repeated restarts of sub-systems, too often refreshes, or the lack of thereof, and a lot of another
        kinds of anomalies, by accounting for growth or fall of specific trends of logging rates in
        general, as well as by specific sub-systems, points, or groups of codes.

    10) Ability to make graphs in grafana of specific codes or groups of them, and visualize
        the key states of the service, such as start, restart, etc. For more complex services,
        it will allow to visualize its current states as state diagrams.

        You can also see when the service panics, without the need to look after it in logs.

    11) Ability to use Msgtag struct for production-level readiness sensor, which can be hard-coded
        for specific logging codes, and keep an account what subsystems are running already.

    12) Возможность использовать коды логирования для heartbeat sensor подсистем, по аналогии с
        readiness sensor.

    12) Ability to use logging codes for a production heartbeat sensor systems, analogously to
        the readiness sensor systems.

    13) If you use this interface wrapper, you can always make your own implementation easily,
        and employ different logger underneath.

    14) Ability to use reported metrics for ML systems and other pattern finding algorithms,
        for determining the health and runtime state and modes of operation of services.

--

Some technical details: why wasn't standard hooks mechanism of zerolog wasn't enough?

    Well. We can use them. However, you need to have some logging attributes in them as scalar objects,
    and the point of zerolog is rather to do immediate serialization. So, to get those, you'd need to
    deserialize it back, which is absurd. With this wrapper, we are free to put and keep any attributes
    we want as scalars, and then use in our own sending hook.

    Currently, in the reference implementation, the msgtag array, the title, level, and subsystem are
    all available as scalars.

--
