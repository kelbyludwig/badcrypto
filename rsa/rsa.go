package rsa

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

//PublicKey represents the public half of an RSA keypair.
type PublicKey struct {
	N *big.Int // Modulus
	E int64    //Public exponent
}

//PrivateKey represents the private half of an RSA keypair.
type PrivateKey struct {
	PublicKey *PublicKey
	D         *big.Int
	Primes    []*big.Int
}

//EncryptNoPadding encrypts the supplied plaintext byte slice using the supplied public key.
//EncryptNoPadding does not pad the plaintext prior to encryption.
func EncryptNoPadding(plaintext []byte, publicKey *PublicKey) (ciphertext []byte) {
	num := new(big.Int).SetBytes(plaintext)
	ct := new(big.Int).Exp(num, big.NewInt(publicKey.E), publicKey.N)
	ciphertext = ct.Bytes()
	return
}

//DecryptNoPadding decrypts the supplied ciphertext using the supplied PrivateKey.
//DecryptNoPadding does not validate or strip off any form of padding.
func DecryptNoPadding(ciphertext []byte, privateKey *PrivateKey) (plaintext []byte) {
	num := new(big.Int).SetBytes(ciphertext)
	N := privateKey.PublicKey.N
	pt := new(big.Int).Exp(num, privateKey.D, N)
	plaintext = pt.Bytes()
	return
}

//GenerateKey generates an RSA private key (and corresponding public key)
//given the size of a modulus in bits.
func GenerateKey(bits int) (priv *PrivateKey, err error) {

	if bits%2 != 0 {
		err = fmt.Errorf("bits must be a multiple of 2")
		return
	}

	pub := new(PublicKey)
	pub.E = 3

	priv = new(PrivateKey)
	priv.Primes = make([]*big.Int, 2)

	for {
		p, err1 := rand.Prime(rand.Reader, bits)
		q, err2 := rand.Prime(rand.Reader, bits)

		if err1 != nil || err2 != nil {
			err = fmt.Errorf("unable to generate prime numbers")
			return
		}

		priv.Primes[0] = p
		priv.Primes[1] = q

		totient1 := new(big.Int).Sub(p, big.NewInt(1))
		totient2 := new(big.Int).Sub(q, big.NewInt(1))
		totient := new(big.Int).Mul(totient1, totient2)

		gcd := new(big.Int).GCD(nil, nil, totient, big.NewInt((pub.E)))
		if gcd.Cmp(big.NewInt(1)) == 0 {
			pub.N = new(big.Int).Mul(p, q)
			priv.D = new(big.Int).ModInverse(big.NewInt(int64(pub.E)), totient)
			break
		}
	}

	priv.PublicKey = pub
	return

}

//ChineseRemainderTheorem solves a set of congruences of the form:
//  x = a1 (mod m1)
//  x = a2 (mod m2)
//  ...
//  x = an (mod mn)
//ChineseRemainderTheorem takes the set of ai and mi as input and returns x.
func ChineseRemainderTheorem(as, ms []*big.Int) (*big.Int, error) {

	if len(as) != len(ms) {
		return big.NewInt(0), fmt.Errorf("lists provided were unequal in lenght")
	}

	M := big.NewInt(1)
	for _, m := range ms {
		gcd := new(big.Int).GCD(nil, nil, M, m)
		if gcd.Cmp(big.NewInt(1)) != 0 {
			return big.NewInt(0), fmt.Errorf("moduli were not comprime")
		}
		M = M.Mul(M, m)
	}

	result := big.NewInt(0)
	for i, a := range as {
		b := new(big.Int).Div(M, ms[i])
		bi := new(big.Int).ModInverse(b, ms[i])
		mul := new(big.Int).Mul(b, bi)
		mul = mul.Mod(mul, M)
		mul = mul.Mul(mul, a)
		mul = mul.Mod(mul, M)
		result = result.Add(result, mul)
	}
	result = result.Mod(result, M)
	return result, nil
}
