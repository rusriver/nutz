This is open-sourced in 2024, but was developed long ago, and is quite aged technology.
Initially developed as a demonstrator for an idea of reverse cluster discovery, it found
its use as a counter of cluster replicas.

Isn't applicable if there are sticky sessions/connections.

The algorithm is this:

    We have two TTL caches, slow mode cache and fast mode cache.

    SLOW MODE (initial):
        - Initially we think there is one replica, and, for example, the polling period is 30s,
            the TTL and fast mode period are x2 of poll period, fast mode speed is x20.
        - We make a first call, and immediatelly discover the new replica and switch to fast mode.
        - In slow mode we can not report result at all, because it is already reported while in
            the fast mode.

    FAST MODE:
        - We have TTL x20 shorter, i.e. the 3s, and a poll period would be 1.5s initially.
        - We can report the result every 3s, and on every report recalculate the poll period.
            For example, on first recalc we'd have 2 replicas, and poll period becomes 0.75s,
            on second recalc - 3 repl and 0.5s period, then - 5 repl and 0.3s, and so on.
            The reported number of replicas would grow fast and smooth.
        - If we have had replicas already, then on their mass replacement will happen a smooth
            displacement of old ones with new ones, new ones will be added, old ones will expire
            by TTL, and the reporten number will be stable.
        - If for the fast mode period there were added no new replicas, then we switch back to slow mode.
            The slow mode cache must be already filled with new replicas.

