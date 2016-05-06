package big

import (
	"math/big"
)

//MontgomeryReduction computes the montgomery reduction
//of t modulo m with respect to R = 2^n where n is the
//bit length of m.
//This algorithm is based off the implementation in
//"Handbook of Applied Cryptography".
func MontgomeryReduction(m, t *big.Int) *big.Int {

	b := big.NewInt(2)
	n := m.BitLen()

	R := new(big.Int).Exp(b, big.NewInt(int64(n)), nil)

	nm := new(big.Int).Neg(m)
	nm = nm.Mod(nm, R)
	mp := new(big.Int).ModInverse(nm, R)

	A := new(big.Int).Set(t)

	ui := new(big.Int).Mod(mp, b)
	ui = ui.Mul(ui, m)
	bi := new(big.Int)
	for i := 0; i < n; i++ {
		ai := A.Bit(i)
		if ai == 1 {
			bi = bi.Exp(b, big.NewInt(int64(i)), nil)
			ad := new(big.Int).Mul(ui, bi)
			A = A.Add(A, ad)
		}
	}

	A = A.Rsh(A, uint(n))
	if cmp := A.Cmp(m); cmp == 0 || cmp == 1 {
		A = A.Sub(A, m)
	}

	return A

}
