# Koonin's keyhash algorithm

Implementation of Koonin's keyhash 8 bits (one way) and 16 bits (double way) algorithms.

Exceptionally fast, doesn't contain any modulo or division operations. It yields even distribution
of hash values, even in edge cases, when only one byte changes by 1.

This does not fit everywhere, but in certain application areas its characteristics are good.

# Latitov's keyhash algorithm

Gives 17% better distribution uniformness, and passes all the requirements of the
Koonin's one. A little bit more expensive though.

# Rendezvous hashing

TLDR: To get excellent results, use XOR double-wrap of the key, then Latitov's hash.
Only if your key is UUID+UUID, or just two long strings. If one key is incremental ID,
then read below.

Somehow, neither hash function gives even uniform distribution when doing rendezvous.
But with double-wrap, either one is good, and Latitov's yields 10 times smaller
(better) min-max delta than Koonin's, 350 vs 4800, median, accordingly, on the same test.

What about other hashes? Here are the results, smaller the better, as performed on the same test, A - plain key, B - with double-wrap:

- FNV (https://en.wikipedia.org/wiki/Fowler-Noll-Vo_hash_function):
  - A: 3100 median
  - B: 350 median

- CRC32 with IEEE table:
  - A: 12500 median
  - B: 12500 median

- SHA256:
  - A: 350 median
  - B: 350 median

- MD5:
  - A: 350 median
  - B: 350 median

Note that cryptographic hashes also substantially slower.

Note that these results are for rendezvous hashing ONLY, and if you use the hash directly,
there are different comparison criteria.

Both Koonin's and Latitov's functions are roughly as fast as CRC32, but are suitable
for cases where uniformness is required (sharding, rendezvous, etc).

(The test was `xm/xm-http/Test_001()`)

Turns out, if one key is incremental counter, then double-wrap then Latitov's hash gives
8800, which isn't good obviously. The MD5 gives 500+, and double-wrap then MD5 gives best results, the 250 median. By the way, the MD5 is faster than SHA.

Turns out, the MD5, even with double-wrap, gives uneven results on small datasets of double UUIDs.
The scheme that gives the best results was found to be ShuffleBytes(5) then Latitov's hash - it
gives the delta in 300..500, for `xm/xm-http/Test_001()` and `xm/xm-http/Test_001_1()` and `xm/xm-http/Test_002()`.

