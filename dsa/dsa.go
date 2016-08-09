package dsa

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"math/big"
)

type PublicKey struct {
	P *big.Int
	Q *big.Int
	G *big.Int
	Y *big.Int
}

type PrivateKey struct {
	PublicKey *PublicKey
	X         *big.Int
}

//p,q, and g are hardcoded DSA parameters that
//are set by init().
var p *big.Int
var q *big.Int
var g *big.Int

func init() {
	var ok1, ok2, ok3 bool
	p, ok1 = new(big.Int).SetString("800000000000000089e1855218a0e7dac38136"+
		"ffafa72eda7859f2171e25e65eac698c170257"+
		"8b07dc2a1076da241c76c62d374d8389ea5aef"+
		"fd3226a0530cc565f3bf6b50929139ebeac04f"+
		"48c3c84afb796d61e5a4f9a8fda812ab594942"+
		"32c7d2b4deb50aa18ee9e132bfa85ac4374d7f"+
		"9091abc3d015efc871a584471bb1", 16)
	q, ok2 = new(big.Int).SetString("f4f47f05794b256174bba6e9b396a7707e563c5b", 16)
	g, ok3 = new(big.Int).SetString("5958c9d3898b224b12672c0b98e06c60df923cb8bc999d119"+
		"458fef538b8fa4046c8db53039db620c094c9fa077ef389b5"+
		"322a559946a71903f990f1f7e0e025e2d7f7cf494aff1a047"+
		"0f5b64c36b625a097f1651fe775323556fe00b3608c887892"+
		"878480e99041be601a62166ca6894bdd41a7054ec89f756ba"+
		"9fc95302291", 16)

	if !ok1 || !ok2 || !ok3 {
		panic("failed to generate hardcoded dsa params")
	}
}

//GenerateKey generates a new DSA signing keypair.
//GenerateKey uses a fixed set of parameters for simplicity.
func GenerateKey() (priv *PrivateKey, err error) {
	priv = new(PrivateKey)
	priv.PublicKey = new(PublicKey)
	priv.PublicKey.P = p
	priv.PublicKey.G = g
	priv.PublicKey.Q = q

	x := make([]byte, len(q.Bytes()))
	_, err = rand.Read(x)
	if err != nil {
		return
	}
	priv.X = new(big.Int).SetBytes(x)
	priv.X = priv.X.Mod(priv.X, q)
	priv.PublicKey.Y = new(big.Int).Exp(g, priv.X, p)
	return
}

//Sign DSA signs message using the supplied private key.
func Sign(message []byte, privateKey *PrivateKey) (r, s *big.Int, err error) {

	var k *big.Int
	kbuf := make([]byte, len(privateKey.PublicKey.Q.Bytes()))
	for {
		//Generate per-message key
		_, err = rand.Read(kbuf)
		if err != nil {
			return
		}
		k = new(big.Int).SetBytes(kbuf)
		k = k.Mod(k, privateKey.PublicKey.Q)

		r = new(big.Int).Exp(privateKey.PublicKey.G, k, privateKey.PublicKey.P)
		r = r.Mod(r, privateKey.PublicKey.Q)

		if r.Cmp(big.NewInt(0)) != 0 {
			break
		}
	}
	kinv := new(big.Int).ModInverse(k, privateKey.PublicKey.Q)
	xr := new(big.Int).Mul(privateKey.X, r)
	digest := sha1.Sum(message)
	digestNum := new(big.Int).SetBytes(digest[:])
	s = new(big.Int).Add(digestNum, xr)
	s = s.Mod(s, privateKey.PublicKey.Q)
	s = s.Mul(s, kinv)
	s = s.Mod(s, privateKey.PublicKey.Q)
	return
}

//Verify verifies a signature (r,s) for message under the supplied publicKey.
//Returns a non-nil error on signature validation failure.
func Verify(message []byte, r, s *big.Int, publicKey *PublicKey) error {

	//Reject the signature if 0<r<q or 0<s<q is not satisfied.
	zero := big.NewInt(0)
	if r.Cmp(zero) <= 0 ||
		s.Cmp(zero) <= 0 ||
		r.Cmp(publicKey.Q) == 1 ||
		s.Cmp(publicKey.Q) == 1 {
		fmt.Printf("1\n")
		return fmt.Errorf("invalid signature")
	}

	w := new(big.Int).ModInverse(s, publicKey.Q)
	digest := sha1.Sum(message)
	digestNum := new(big.Int).SetBytes(digest[:])
	u1 := new(big.Int).Mul(digestNum, w)
	u1 = u1.Mod(u1, publicKey.Q)
	u2 := new(big.Int).Mul(r, w)
	u2 = u2.Mod(u2, publicKey.Q)
	gu1 := new(big.Int).Exp(publicKey.G, u1, publicKey.P)
	gu2 := new(big.Int).Exp(publicKey.Y, u2, publicKey.P)
	gu := new(big.Int).Mul(gu1, gu2)
	gu = gu.Mod(gu, publicKey.P)
	v := new(big.Int).Mod(gu, publicKey.Q)

	if v.Cmp(r) != 0 {
		fmt.Printf("2\n")
		fmt.Printf("v: %x\n", v)
		fmt.Printf("r: %x\n", r)
		return fmt.Errorf("invalid signature")
	} else {
		return nil
	}

}
