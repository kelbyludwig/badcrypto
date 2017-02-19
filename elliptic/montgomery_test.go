package elliptic

import (
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
	morder, _ = new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	//TODO(kkl): using zero for curve order until SEA is implmented.
	mcurve = NewMontgomeryCurve(ma, mb, mp, zero, mgx, mgy)
}

func TestLadder(t *testing.T) {

	w := mcurve.ScalarMult(mgx, morder.Bytes())
	if w.Cmp(zero) != 0 {
		t.Errorf("scalarmult by order did not equal zero")
		return
	}
}
