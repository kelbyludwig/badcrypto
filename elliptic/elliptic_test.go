package elliptic

import (
	"math/big"
	"testing"
)

var a *big.Int
var b *big.Int
var p *big.Int
var gx *big.Int
var gy *big.Int
var curve shortWeierstrassCurve

func init() {
	a, _ = new(big.Int).SetString("-95051", 10)
	b, _ = new(big.Int).SetString("11279326", 10)
	p, _ = new(big.Int).SetString("233970423115425145524320034830162017933", 10)
	gx, _ = new(big.Int).SetString("182", 10)
	gy, _ = new(big.Int).SetString("85518893674295321206118380980485522083", 10)
	curve = NewCurve(a, b, p, gx, gy)
}

func TestCurveAddition(t *testing.T) {
	zX := big.NewInt(0)
	zY := big.NewInt(1)

	x1, y1 := curve.Add(zX, zY, gx, gy)
	if !curve.PointEquals(x1, y1, gx, gy) {
		t.Errorf("adding by zero did not return original point")
		return
	}

	x2, y2 := curve.Add(gx, gy, zX, zY)
	if !curve.PointEquals(gx, gy, x2, y2) {
		t.Errorf("adding by zero did not return original point")
		return
	}

}

func TestCurveInverse(t *testing.T) {

	x0, y0 := curve.invertPoint(gx, gy)
	x0i, y0i := curve.invertPoint(x0, y0)

	if !curve.PointEquals(x0i, y0i, gx, gy) {
		t.Errorf("inverting the inverse did not results in the same point")
		return
	}
}

func TestCurveDouble(t *testing.T) {

	gx2, gy2 := curve.Double(gx, gy)
	gxi, gyi := curve.invertPoint(gx, gy)
	x2, y2 := curve.Add(gx2, gy2, gxi, gyi)

	if !curve.PointEquals(x2, y2, gx, gy) {
		t.Errorf("point subtraction did not return base point")
		return
	}

	x3, y3 := curve.Add(gx, gy, gxi, gyi)
	if !curve.PointEquals(x3, y3, zero, one) {
		t.Errorf("addition of inverse did not return zero point")
		return
	}
}

func TestScalarMult(t *testing.T) {

	gx2, gy2 := curve.ScalarMult(gx, gy, two.Bytes())
	gxd, gyd := curve.Double(gx, gy)

	if !curve.PointEquals(gx2, gy2, gxd, gyd) {
		t.Errorf("scalar multiplication by 2 did not result in the same point as doubling")
		return
	}

}
