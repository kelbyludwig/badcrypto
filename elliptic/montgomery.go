package elliptic

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	bbig "github.com/kelbyludwig/badcrypto/big"
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

	//I am not proud of this code...
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

func (curve montgomeryCurve) randomPoint() (x *big.Int) {
	buf := make([]byte, len(curve.P.Bytes()))
	rand.Read(buf)
	x = new(big.Int).SetBytes(buf)
	x = x.Mod(x, curve.P)

	x3 := new(big.Int).Exp(x, three, curve.P)

	ax2 := new(big.Int).Exp(x, two, curve.P)
	ax2 = ax2.Mul(curve.A, ax2)
	ax2 = ax2.Mod(ax2, curve.P)

	rhs := new(big.Int).Add(x3, ax2)
	rhs = rhs.Add(rhs, x)
	rhs = rhs.Mod(rhs, curve.P)
	return rhs
}

func (curve montgomeryCurve) isZeroPoint(x *big.Int) bool {
	zero := big.NewInt(0)
	xt := new(big.Int).Mod(x, curve.P)
	return curve.PointEquals(xt, zero)
}

func (curve montgomeryCurve) twistPointWithSpecifiedOrder(twistOrder, primeFactor *big.Int) *big.Int {
	for {
		x := curve.randomPoint()
		xsqrt := new(big.Int).ModSqrt(x, curve.P)
		if xsqrt != nil {
			continue
		}
		fmt.Printf("%d is not square mod %d\n", x, curve.P)
		nr := new(big.Int).Div(twistOrder, primeFactor)
		a := curve.ScalarMult(x, nr.Bytes())
		if !curve.isZeroPoint(a) {
			return a
		}
	}
}

func (curve montgomeryCurve) PointEquals(x1, x2 *big.Int) bool {
	//ensure the points are reduced mod P
	x1r := new(big.Int).Mod(x1, curve.P)
	x2r := new(big.Int).Mod(x2, curve.P)

	if x1r.Cmp(x2r) == 0 {
		return true
	}
	return false
}

func (curve montgomeryCurve) ComputeIndexWithinRange(a, x, min, max *big.Int) (index *big.Int, err error) {
	index = new(big.Int).SetBytes(min.Bytes())
	for index.Cmp(max) != 1 {
		aa := curve.ScalarMult(x, index.Bytes())
		if curve.PointEquals(a, aa) {
			return
		}
		index = index.Add(index, one)
	}
	return nil, IndexNotRecoveredErr

}

func (curve montgomeryCurve) PohligHellmanOnline(oracle func(*big.Int) *big.Int) (index, newmod *big.Int, err error) {

	indices := make([]*big.Int, 0)
	moduli := make([]*big.Int, 0)

	//Count points on the twist using knowledge of the order of the curve.
	curveOrder := curve.N
	fmt.Printf("curve order %d\n", curveOrder)
	totalPoints := new(big.Int).Mul(big.NewInt(2), curve.P)
	totalPoints = totalPoints.Add(totalPoints, big.NewInt(2))
	fmt.Printf("total points %d\n", totalPoints)
	twistOrder := new(big.Int).Sub(totalPoints, curveOrder)
	fmt.Printf("twist order %d\n", twistOrder)

	factors, _ := bbig.Factor(twistOrder, 65536)

	for factor, _ := range factors {
		primeFactor := big.NewInt(int64(factor))
		if primeFactor.Cmp(two) == 0 {
			continue
		}

		ind := new(big.Int)
		//using an outer loop because twistPointWithSpecifiedOrder sometimes returns non-twist points and i'm not sure why.
		for {
			fmt.Printf("new factor %d\n", primeFactor)
			x := mcurve.twistPointWithSpecifiedOrder(twistOrder, primeFactor)
			y := oracle(x)

			ind, err = mcurve.ComputeIndexWithinRange(y, x, zero, primeFactor)
			if err != nil {
				fmt.Printf("failed to recover index...\n")
				continue
			}
			break
		}
		fmt.Printf("recovered index. new residue x mod %d = %d\n", primeFactor, ind)
		indices = append(indices, ind)
		moduli = append(moduli, primeFactor)
	}

	return bbig.CRT(indices, moduli)

}
