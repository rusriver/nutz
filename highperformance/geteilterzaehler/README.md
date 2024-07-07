### High-performance universal counter, optimized for heavy concurrent use

Uses a racing shuttle, with user-injected extra entropy, to minimize the mutex contention of individual shards.

When used properly, it is usually N times faster than a simple mutex protected counter, where N is the number
of real CPU cores, when using simple float64 counter. If the counting logic is more advanced, the gain can be more substantial.

Yields precise results, but doesn't guarantee specific time of measurement. See Heisenberg's indeterminacy principle,
https://en.wikipedia.org/wiki/Uncertainty_principle.

### References

Computer Networks Volume 244, May 2024, 110343
TbDd: A new trust-based, DRL-driven framework for blockchain sharding in IoT
Zixu Zhang, Guangsheng Yu, Caijun Sun, Xu Wang, Ying Wang, Ming Zhang, Wei Ni, Ren Ping Liu, Andrew Reevesf, Nektarios Georgalasf

