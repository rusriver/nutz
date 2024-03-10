# Koonin's keyhash algorithm

Implementation of Koonin's keyhash 8 bits (one way) and 16 bits (double way) algorithms.

Exceptionally fast, doesn't contain any modulo or division operations. It yields even distribution
of hash values, even in edge cases, when only one byte changes by 1.

This does not fit everywhere, but in certain application areas its characteristics are pretty charming.

