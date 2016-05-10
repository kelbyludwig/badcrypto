## RSA Timing Attack

Attack described in [Remote Timing Attacks Are
Practical](https://crypto.stanford.edu/~dabo/papers/ssl-timing.pdf).
Takes advantage of several OpenSSL optimizations of RSA which expose timing
side-channels. These optimizations include:

* Chinese Remainder Theorem (CRT) Optimized Decryption

* Sliding Window Exponentiation (An optimization of square-and-multiply)

* Montgomery Multiplication and Reductions

* Karatsuba Multiplication

### Timing Differences: CRT 

* TODO Not explained yet. I think it suggested there is not timing difference
  for CRT optimizations.

### Timing Differences: Sliding Window Exponentiation

* Sliding Window exponentiation is a optimization of square-and-multiply 
  that processes groups of bits at once. It is used to decrypt ciphertexts
  (i.e. g^d (mod q) )

* g is the attacker-controlled ciphertext. Sliding window exponentation 
  computes several multiplications of g. This can cause timing differences
  based on g's relationship with q. q is one of the factors of the private
  key.

### Timing Differences: Montgomery Reductions

* Montgomery reductions are optimized modular reductions. Instead of relying
  on expensive division operations to perform modular reduction, the number
  being reduced is put in "Montgomery form." Even though putting a number
  in Montgomery form does incur some minor costs, the reduction can then be
  done by bitshifts instead of divisions.

* Montgomery reductions sometimes have an extra reduction step to keep the 
  reduced number below the modulus. This step is optional and has a property
  where a number that is the multiple of the modulus is not likely to need a
  final reduction (and is much faster).

### Timing Differences: Karatsuba Multiplication

* Karatsuba multiplication is a optimization of same-size bignums. If two
  bignums are equal in size, there is a timing difference.
