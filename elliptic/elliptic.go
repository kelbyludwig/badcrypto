package elliptic

import (
	"crypto/elliptic"
	"math/big"
)

var zero *big.Int = big.NewInt(0)
var one *big.Int = big.NewInt(1)
var two *big.Int = big.NewInt(2)
var three *big.Int = big.NewInt(3)

type shortWeierstrassCurve struct {
	*elliptic.CurveParams
	A *big.Int
}

//NewCurve creates a new curve that implements the `elliptic.Curve` interface
//with custom parameters.
func NewCurve(a, b, p, gx, gy *big.Int) (curve shortWeierstrassCurve) {
	curveParams := elliptic.CurveParams{
		Name:    "short weierstrass curve",
		BitSize: 0,
		N:       zero,
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

func (curve shortWeierstrassCurve) Params() *elliptic.CurveParams {
	return curve.CurveParams
}

func (curve shortWeierstrassCurve) IsOnCurve(x, y *big.Int) bool {
	panic("not implemented")
	return false
}

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

func (curve shortWeierstrassCurve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	return curve.Add(x1, y1, x1, y1)
}

func (curve shortWeierstrassCurve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	K := new(big.Int).SetBytes(k)
	Nx := new(big.Int).SetBytes(x1.Bytes())
	Ny := new(big.Int).SetBytes(y1.Bytes())
	Qx := big.NewInt(0)
	Qy := big.NewInt(1)

	for i := K.BitLen(); i >= 0; i-- {
		bit := K.Bit(i)
		if bit == 1 {
			Qx, Qy = curve.Add(Qx, Qy, Nx, Ny)
		}
		Nx, Ny = curve.Double(Nx, Ny)
	}
	return Qx, Qy
}

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
