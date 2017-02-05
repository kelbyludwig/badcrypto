package elliptic

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	bbig "github.com/kelbyludwig/badcrypto/big"
)

var zero *big.Int = big.NewInt(0)
var one *big.Int = big.NewInt(1)
var two *big.Int = big.NewInt(2)
var three *big.Int = big.NewInt(3)

var IndexNotRecoveredErr error = fmt.Errorf("index not found")

type scalarMultOracle func(x, y *big.Int) (*big.Int, *big.Int)

type shortWeierstrassCurve struct {
	*elliptic.CurveParams
	A *big.Int
}

//NewCurve creates a new curve that implements the `elliptic.Curve` interface
//with custom parameters. The curve represented by shortWeierstrassCurve has
//the following structure:y^2 = x^3 + a*x + b.  p is the order of the
//underlying field and n is the order of the base point (gx,gy).
func NewCurve(a, b, p, n, gx, gy *big.Int) (curve shortWeierstrassCurve) {
	curveParams := elliptic.CurveParams{
		Name:    "short weierstrass curve",
		BitSize: 0,
		N:       n,
		Gx:      gx,
		Gy:      gy,
		B:       b,
		P:       p,
	}
	curve = shortWeierstrassCurve{
		CurveParams: &curveParams,
		A:           a,
	}
	return curve
}

//Params returns the curves CurveParams struct.
func (curve shortWeierstrassCurve) Params() *elliptic.CurveParams {
	return curve.CurveParams
}

//IsOnCurve will return true if the supplied point (x, y) is a valid point
//for the curve and false otherwise.
func (curve shortWeierstrassCurve) IsOnCurve(x, y *big.Int) bool {
	//y^2 = x^3 + a*x + b
	lhs := new(big.Int).Exp(y, two, curve.P)
	rhs := new(big.Int).Exp(x, three, curve.P)
	ax := new(big.Int).Mul(curve.A, x)
	rhs = rhs.Add(rhs, ax)
	rhs = rhs.Add(rhs, curve.B)
	rhs = rhs.Mod(rhs, curve.P)

	if lhs.Cmp(rhs) == 0 {
		return true
	}
	return false
}

//Add implements generic Short Weierstrass curve addition.
func (curve shortWeierstrassCurve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	if curve.isZeroPoint(x1, y1) {
		x = new(big.Int).SetBytes(x2.Bytes())
		y = new(big.Int).SetBytes(y2.Bytes())
		return
	}

	if curve.isZeroPoint(x2, y2) {
		x = new(big.Int).SetBytes(x1.Bytes())
		y = new(big.Int).SetBytes(y1.Bytes())
		return
	}

	x2i, y2i := curve.invertPoint(x2, y2)
	if curve.PointEquals(x1, y1, x2i, y2i) {
		x = new(big.Int).SetBytes(zero.Bytes())
		y = new(big.Int).SetBytes(one.Bytes())
		return
	}

	m := new(big.Int)
	if curve.PointEquals(x1, y1, x2, y2) {
		//m = (3*x1^2 + a) / 2*y1
		m = m.Exp(x1, two, curve.P)
		m = m.Mul(m, three)
		m = m.Add(m, curve.A)
		bot := new(big.Int).Mul(y1, two)
		bot = bot.Mod(bot, curve.P)
		bot = bot.ModInverse(bot, curve.P)
		m = m.Mul(m, bot)
		m = m.Mod(m, curve.P)
	} else {
		//m = (y2 - y1) / (x2 - x1)
		m = m.Sub(y2, y1)
		bot := new(big.Int).Sub(x2, x1)
		bot = bot.Mod(bot, curve.P)
		bot = bot.ModInverse(bot, curve.P)
		m = m.Mul(m, bot)
		m = m.Mod(m, curve.P)
	}

	x = new(big.Int).Exp(m, two, curve.P)
	x = x.Sub(x, x1)
	x = x.Sub(x, x2)
	x = x.Mod(x, curve.P)

	y = new(big.Int).Sub(x1, x)
	y = y.Mul(m, y)
	y = y.Sub(y, y1)
	y = y.Mod(y, curve.P)
	return
}

//Double returns the supplied point doubled.
func (curve shortWeierstrassCurve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	return curve.Add(x1, y1, x1, y1)
}

//ScalarMult returns k*(x1, y1).
func (curve shortWeierstrassCurve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	K := new(big.Int).SetBytes(k)
	Qx := big.NewInt(0)
	Qy := big.NewInt(1)

	for i := K.BitLen(); i >= 0; i-- {
		bit := K.Bit(i)
		Qx, Qy = curve.Double(Qx, Qy)
		if bit == 1 {
			Qx, Qy = curve.Add(Qx, Qy, x1, y1)
		}
	}
	return Qx, Qy
}

//ScalarBaseMult returns k*(x1, y1) where (x1, y1) is the base point for the
//supplied curve.
func (curve shortWeierstrassCurve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return curve.ScalarMult(curve.Gx, curve.Gy, k)
}

//isZeroPoint will return true if the supplied (x,y) values are the "zero"
//value on a curve.
func (curve shortWeierstrassCurve) isZeroPoint(x, y *big.Int) bool {
	one := big.NewInt(1)
	zero := big.NewInt(0)
	return curve.PointEquals(x, y, zero, one)
}

//PointEquals will return true if the supplied points (x1, y1) and (x2, y2) are
//equal on the given curve.
func (curve shortWeierstrassCurve) PointEquals(x1, y1, x2, y2 *big.Int) bool {
	//ensure the points are reduced mod P
	x1r := new(big.Int).Mod(x1, curve.P)
	y1r := new(big.Int).Mod(y1, curve.P)
	x2r := new(big.Int).Mod(x2, curve.P)
	y2r := new(big.Int).Mod(y2, curve.P)

	if x1r.Cmp(x2r) == 0 && y1r.Cmp(y2r) == 0 {
		return true
	}
	return false
}

//invertPoint inverts the point (x, y) on the given curve.
func (curve shortWeierstrassCurve) invertPoint(x, y *big.Int) (xi, yi *big.Int) {
	xi = new(big.Int).SetBytes(x.Bytes())
	yi = new(big.Int).Sub(curve.P, y)
	return
}

//randomPoint generates a random point on the given curve.
func (curve shortWeierstrassCurve) randomPoint() (x, y *big.Int) {
	for {
		buf := make([]byte, len(curve.P.Bytes()))
		rand.Read(buf)
		x = new(big.Int).SetBytes(buf)
		x = x.Mod(x, curve.P)
		rhs := new(big.Int).Exp(x, three, curve.P)
		ax := new(big.Int).Mul(curve.A, x)
		rhs = rhs.Add(rhs, ax)
		rhs = rhs.Add(rhs, curve.B)
		rhs = rhs.Mod(rhs, curve.P)
		y = new(big.Int).ModSqrt(rhs, curve.P)
		if y != nil {
			y = y.Mod(y, curve.P)
			break
		}
	}
	return
}

func (curve shortWeierstrassCurve) pointWithSpecifiedOrder(r *big.Int) (*big.Int, *big.Int) {

	for {
		x, y := curve.randomPoint()
		nr := new(big.Int).Div(curve.N, r)
		a, b := curve.ScalarMult(x, y, nr.Bytes())
		if !curve.isZeroPoint(a, b) {
			return a, b
		}
	}

}

//ComputeIndexWithinRange will solve for s in the equation s*(x,y) = (a,b) where P
//is curve's base point. If the index does not fall within the specified range,
//this function will return an error.
func (curve shortWeierstrassCurve) ComputeIndexWithinRange(a, b, x, y, min, max *big.Int) (index *big.Int, err error) {
	index = new(big.Int).SetBytes(min.Bytes())
	for index.Cmp(max) != 1 {
		aa, bb := curve.ScalarMult(x, y, index.Bytes())
		if curve.PointEquals(a, b, aa, bb) {
			return
		}
		index = index.Add(index, one)
	}
	return nil, IndexNotRecoveredErr

}

//pohligHellmanOnline implements the invalid curve attack against a specified
//curve `curve` and an oracle function `oracle` that computes scalarmults on
//the input point. This method takes pre-generated small-order curves as input
//because I have not implemented point counting yet.
func (curve shortWeierstrassCurve) pohligHellmanOnline(smallOrderCurves []shortWeierstrassCurve, oracle scalarMultOracle) (index, newmod *big.Int, err error) {

	for _, soc := range smallOrderCurves {
		if soc.N != nil && soc.N.Cmp(zero) == 0 {
			return nil, nil, fmt.Errorf("order of smallOrderCurve not supplied")
		}
	}

	indices := make([]*big.Int, 0)
	moduli := make([]*big.Int, 0)

	for _, soc := range smallOrderCurves {
		mmo := new(big.Int).SetBytes(soc.N.Bytes())
		factors, _ := bbig.Factor(mmo, 1048576)

	NewFactor:
		for factor, _ := range factors {
			primeFactor := big.NewInt(int64(factor))
			if primeFactor.Cmp(two) == 0 {
				continue
			}
			for _, rec := range moduli {
				gcd := new(big.Int).GCD(nil, nil, rec, primeFactor)
				if gcd.Cmp(one) != 0 {
					continue NewFactor
				}
			}
			x, y := soc.pointWithSpecifiedOrder(primeFactor)
			xx, yy := oracle(x, y)
			ind, err := soc.ComputeIndexWithinRange(xx, yy, x, y, zero, primeFactor)
			if err != nil {
				continue
			}
			indices = append(indices, ind)
			moduli = append(moduli, primeFactor)
		}
	}
	return bbig.CRT(indices, moduli)

}
