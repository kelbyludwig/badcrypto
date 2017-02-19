package elliptic

import (
	"fmt"
	"math/big"
	"testing"
)

var ma *big.Int
var mb *big.Int
var mp *big.Int
var mgx *big.Int
var mgy *big.Int
var morder *big.Int
var mcurve montgomeryCurve

func init() {
	ma, _ = new(big.Int).SetString("534", 10)
	mb, _ = new(big.Int).SetString("1", 10)
	mp, _ = new(big.Int).SetString("233970423115425145524320034830162017933", 10)
	mgx, _ = new(big.Int).SetString("4", 10)
	mgy, _ = new(big.Int).SetString("85518893674295321206118380980485522083", 10)
	morder, _ = new(big.Int).SetString("233970423115425145498902418297807005944", 10)
	//TODO(kkl): using zero for curve order until SEA is implmented.
	mcurve = NewMontgomeryCurve(ma, mb, mp, morder, mgx, mgy)
}

func TestLadder(t *testing.T) {

	w := mcurve.ScalarMult(mgx, morder.Bytes())
	if w.Cmp(zero) != 0 {
		t.Errorf("scalarmult by order did not equal zero")
		return
	}
}

func TestCryptopals60(t *testing.T) {

	priv := big.NewInt(705485)
	oracle := func(x *big.Int) *big.Int {
		fmt.Printf("oracle: %d*%d = ", priv, x)
		a := mcurve.ScalarMult(x, priv.Bytes())
		fmt.Printf("%d\n", a)
		return a
	}
	ind, _, _ := mcurve.PohligHellmanOnline(oracle)
	if ind.Cmp(priv) != 0 {
		t.Errorf("failed to recover private key")
	}
}
