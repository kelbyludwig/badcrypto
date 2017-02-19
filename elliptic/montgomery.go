package elliptic

import (
	"crypto/elliptic"
	"math/big"
)

//montgomeryCurve represents a montgomery curve with the following formula:
//B*y^2=x^3+A*x^2+x
type montgomeryCurve struct {
	*elliptic.CurveParams
	A *big.Int
}

func NewMontgomeryCurve(a, b, p, n, gx, gy *big.Int) (curve montgomeryCurve) {
	curveParams := elliptic.CurveParams{
		Name:    "montgomery curve",
		BitSize: 0,
		N:       n,
		Gx:      gx,
		Gy:      gy,
		B:       b,
		P:       p,
	}
	curve = montgomeryCurve{
		CurveParams: &curveParams,
		A:           a,
	}
	return curve
}

func (curve montgomeryCurve) ScalarMult(x1 *big.Int, k []byte) (w *big.Int) {

	cswap := func(x, y *big.Int, b uint) (*big.Int, *big.Int) {
		sx := new(big.Int).SetBytes(x.Bytes())
		sy := new(big.Int).SetBytes(y.Bytes())
		if b == 1 {
			return sy, sx
		} else {
			return sx, sy
		}
	}

	u2 := big.NewInt(1)
	w2 := big.NewInt(0)
	u3 := new(big.Int).SetBytes(x1.Bytes())
	w3 := big.NewInt(1)

	lp := curve.P.BitLen()

	kb := new(big.Int).SetBytes(k)
	for i := lp - 1; i >= 0; i-- {
		bit := kb.Bit(i)
		u2, u3 = cswap(u2, u3, bit)
		w2, w3 = cswap(w2, w3, bit)

		//u3 = (u2*u3-w2*w3)^2
		ou3 := new(big.Int).SetBytes(u3.Bytes())
		ul := new(big.Int).Mul(u2, u3)
		ul = ul.Mod(ul, curve.P)
		wr := new(big.Int).Mul(w2, w3)
		wr = wr.Mod(wr, curve.P)
		ul = ul.Sub(ul, wr)
		ul = ul.Mod(ul, curve.P)
		u3 = u3.Exp(ul, big.NewInt(2), curve.P)

		//w3 = x1*(u2*w3-w2*u3)^2
		u2w3 := new(big.Int).Mul(u2, w3)
		u2w3 = u2w3.Mod(u2w3, curve.P)
		w2u3 := new(big.Int).Mul(w2, ou3)
		w2u3 = w2u3.Mod(w2u3, curve.P)
		s := new(big.Int).Sub(u2w3, w2u3)
		s = s.Mod(s, curve.P)
		w3 = w3.Exp(s, big.NewInt(2), curve.P)
		w3 = w3.Mul(w3, x1)
		w3 = w3.Mod(w3, curve.P)

		//u2 = (u2^2-w2^2)^2
		ou2 := new(big.Int).SetBytes(u2.Bytes())
		ul = new(big.Int).Exp(u2, big.NewInt(2), curve.P)
		ur := new(big.Int).Exp(w2, big.NewInt(2), curve.P)
		u2 = u2.Sub(ul, ur)
		u2 = u2.Mod(u2, curve.P)
		u2 = u2.Exp(u2, big.NewInt(2), curve.P)

		//w2 = 4*u2*w2*(u2^2 + self.a*u2*w2 + w2^2)
		ul = new(big.Int).Exp(ou2, big.NewInt(2), curve.P) //u2^2
		ur = new(big.Int).Mul(curve.A, ou2)
		ur = ur.Mod(ur, curve.P)
		ur = ur.Mul(ur, w2)
		ur = ur.Mod(ur, curve.P)
		ul = ul.Add(ul, ur)
		ul = ul.Mod(ul, curve.P)
		ur = ur.Exp(w2, big.NewInt(2), curve.P)
		ul = ul.Add(ul, ur)
		ul = ul.Mod(ul, curve.P)
		ul = ul.Mul(ul, w2)
		ul = ul.Mod(ul, curve.P)
		ul = ul.Mul(ul, ou2)
		ul = ul.Mod(ul, curve.P)
		w2 = ul.Mul(ul, big.NewInt(4))
		w2 = w2.Mod(w2, curve.P)

		u2, u3 = cswap(u2, u3, bit)
		w2, w3 = cswap(w2, w3, bit)
	}
	p2 := new(big.Int).Sub(curve.P, big.NewInt(2))
	w2 = w2.Exp(w2, p2, curve.P)
	w2 = w2.Mul(w2, u2)
	w2 = w2.Mod(w2, curve.P)
	return w2
}
