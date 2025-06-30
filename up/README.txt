Like systemd, for DI-containers. Generalizes the graph on a rigorous theoretical basis.

See also the readiness-sensor - it has similar functionality.

---

    AllDone may require specific combination of systems ready/ceded, we can't just
    say that some critical system ceded and so we're ready.

    I propose several independent states:
        NoFreshLeft - none left in status fresh, everyone reported
        ReadinessAchieved - we are ready to go forward, because it matches the readiness pattern

    The readiness pattern is specified in Start():
        map[<name>]<required_status>

    where required_status can be:
        - up
        - up_or_ceded

    Only when and if the pattern matches, the ReadinessAchieved is reported.

    The logic in user code must be something like this:

            if ReadinessAchieved {
                ...
            } else if NoFreshLeft {
                ...
            }

    or better:

            if CanReact {
                ...
            } else if NoFreshLeft {
                ...
            }

    or:

            if CanReact {
                ...
            }

---

    OTDG - one-thread dependencies graph.

        - otdg reports to reactor
        - reactor func starts new otdgs

    Simple back-recursive removal:

        - start with otdg or reactor
        - if reactor:
            - of all the systems in status up, get their otdgs, as reported in SendStatus()
            - call RecursiveWindUp() on all of them in parallel
        - if otdg:
            - remove all components in reverse order
            - get parent reactor, as reported when otdg was created
            - call RecursiveWindUp() on it

    Diamond graph topology very well supported, however, for very complex and high-reliability
    systems it is recommended to instead explicitly program a forward graph for removal of all
    components, in specific sequence and order. For simple things the recursive removal works perfect.

---

    Reactor - for joining parallel processes.
    
    Fork - to split processes into parallel. It's like a counter-reactor, the reactor in reverse.

    All parallel things must be explicitly started from the fork, so that on wind-up it
    could automatically join things. I.e., on wind-up, the fork behaves as a join, and
    a reactor behaves like a fork (split).

    OTDG - a generic container for a one-thread process.
    
    Components which are managed by one OTDG bear no specific name.

    OTDG may be started from the reactor, as well as from other OTDGs. For example,
    lauch of initial prerequisite  processes, maybe logger and config, if there's no
    parallel concurrency, can be combined in several OTDGs without a reactor.
    
    In general, the OTDG is not required to be a graph, it can be a branch, i.e. just
    a large list of launching components, sequential.

    It's convenient to place each OTDG into separate files, but is not required.

---
