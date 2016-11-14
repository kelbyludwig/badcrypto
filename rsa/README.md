### Montgomery Reductions

* "Montgomery multiplication is an algorithm used to perform multiple precision modular multiplication quickly by replacing division (which is a slow operation) by multiplications."

* Putting numbers in Montgomery can improve the performance of algorithms that rely on many modular multiplication operations (e.g. modular exponentiation).

* By investing some up-front costs of putting numbers in Montgomery form, the subsequent multiplication and modular reduction steps will be more efficient.

#### References

* [Understanding the Montgomery Reduction Algorithm](http://alicebob.cryptoland.net/understanding-the-montgomery-reduction-algorithm/)

* [Montgomery Multiplication](http://www.mersennewiki.org/index.php/Montgomery_multiplication)
 
* [Hacker's Delight: Montgomery Multiplication (PDF)](http://www.hackersdelight.org/MontgomeryMultiplication.pdf)
