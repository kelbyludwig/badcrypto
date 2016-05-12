package big

import (
	"math/big"
)

//montgomeryReduction computes the montgomery reduction
//of t modulo m with respect to R = 2^n where n is the
//bit length of m.
//This algorithm is based off the implementation in
//"Handbook of Applied Cryptography".
func montgomeryReduction(m, t *big.Int) *big.Int {

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

//montgomeryMul returns xy (mod m).
//montgomeryMul does not return the result in the
//montgomery domain so it is non-optimized. This version
//is just used to verify algorithm correctness.
//This expects x and y to already be reduced mod m.
func montgomeryMul(x, y, m *big.Int) *big.Int {

	b := big.NewInt(2)
	n := m.BitLen()

	R := new(big.Int).Exp(b, big.NewInt(int64(n)), nil)

	nm := new(big.Int).Neg(m)
	nm = nm.Mod(nm, R)
	mp := new(big.Int).ModInverse(nm, R)

	A := big.NewInt(0)

	ui := new(big.Int).Mod(mp, b)
	ui = ui.Mul(ui, m)

	y0 := big.NewInt(int64(y.Bit(0)))
	for i := 0; i < n; i++ {

		a0 := big.NewInt(int64(A.Bit(0)))
		xi := big.NewInt(int64(x.Bit(i)))
		ui := new(big.Int).Mul(xi, y0)
		ui = ui.Add(ui, a0)
		ui = ui.Mul(ui, mp)
		ui = ui.Mod(ui, b)
		ui = ui.Mul(ui, m)

		ad := new(big.Int).Mul(xi, y)
		ad = ad.Add(ad, ui)
		ad = ad.Add(ad, A)
		A = A.Rsh(ad, uint(1))
	}

	if cmp := A.Cmp(m); cmp == 0 || cmp == 1 {
		A = A.Sub(A, m)
	}

	A = A.Mul(A, R)
	A = A.Mod(A, m)
	return A

}

//MontgomeryMul computes xyR^-1 (mod m) using the
//Montgomery multiplication technique.
//The int that is returned is used for error handling and
//assisting in side-channel research. It will return -1
//when there was an error, 0 if no "extra reduction" was
//performed, and 1 if an "extra reducation" was performed.
func MontgomeryMul(x, y, m *big.Int) (*big.Int, int) {

	if x.Cmp(m) == 1 || y.Cmp(m) == 1 {
		return x, -1
	}

	b := big.NewInt(2)
	n := m.BitLen()

	R := new(big.Int).Exp(b, big.NewInt(int64(n)), nil)

	nm := new(big.Int).Neg(m)
	nm = nm.Mod(nm, R)
	mp := new(big.Int).ModInverse(nm, R)

	A := big.NewInt(0)

	ui := new(big.Int).Mod(mp, b)
	ui = ui.Mul(ui, m)

	y0 := big.NewInt(int64(y.Bit(0)))
	for i := 0; i < n; i++ {

		a0 := big.NewInt(int64(A.Bit(0)))
		xi := big.NewInt(int64(x.Bit(i)))
		ui := new(big.Int).Mul(xi, y0)
		ui = ui.Add(ui, a0)
		ui = ui.Mul(ui, mp)
		ui = ui.Mod(ui, b)
		ui = ui.Mul(ui, m)

		ad := new(big.Int).Mul(xi, y)
		ad = ad.Add(ad, ui)
		ad = ad.Add(ad, A)
		A = A.Rsh(ad, uint(1))
	}

	extra := 0
	//The big bad "extra reduction" step :)
	if cmp := A.Cmp(m); cmp == 0 || cmp == 1 {
		A = A.Sub(A, m)
		extra = 1
	}

	return A, extra

}

//ExpMont computes x^e (mod m) using
//the Montgomery multiplication algorithm.
//The int that is returned is used for error handling and
//assisting in side-channel research. It will return -1
//when there was an error, 0 if no "extra reduction" was
//performed, and 1 if an "extra reducation" was performed.
func MontgomeryExp(x, e, m *big.Int) (*big.Int, int) {

	if x.Cmp(m) == 1 {
		return x, -1
	}
	b := big.NewInt(2)
	n := m.BitLen()
	R := new(big.Int).Exp(b, big.NewInt(int64(n)), nil)
	R2 := new(big.Int).Exp(R, b, m)
	A := R.Mod(R, m)
	xs, _ := MontgomeryMul(x, R2, m)

	//extras tracks how many "extra reductions" where performed
	//over the course of the exponentiation.
	extras := 0
	extra := 0
	for i := e.BitLen(); i >= 0; i-- {
		A, extra = MontgomeryMul(A, A, m)
		extras += extra
		if e.Bit(i) == 1 {
			A, extra = MontgomeryMul(A, xs, m)
			extras += extra
		}
		extra = 0
	}
	A, extra = MontgomeryMul(A, big.NewInt(1), m)
	extras += extra
	return A, extras
}
