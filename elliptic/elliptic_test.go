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
var order *big.Int
var hund *big.Int
var hundX *big.Int
var hundY *big.Int
var curve shortWeierstrassCurve

func init() {
	a, _ = new(big.Int).SetString("-95051", 10)
	b, _ = new(big.Int).SetString("11279326", 10)
	p, _ = new(big.Int).SetString("233970423115425145524320034830162017933", 10)
	gx, _ = new(big.Int).SetString("182", 10)
	gy, _ = new(big.Int).SetString("85518893674295321206118380980485522083", 10)
	order, _ = new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	hund, _ = new(big.Int).SetString("100", 10)
	hundX, _ = new(big.Int).SetString("12246423879899346038895890356990169239", 10)
	hundY, _ = new(big.Int).SetString("58231960761567435246734586214813749649", 10)
	//TODO(kkl): using zero for curve order until SEA is implmented.
	curve = NewCurve(a, b, p, zero, gx, gy)
}

func TestCryptopals59(t *testing.T) {

	b1 := big.NewInt(210)
	o1, _ := new(big.Int).SetString("233970423115425145550826547352470124412", 10)
	b2 := big.NewInt(504)
	o2, _ := new(big.Int).SetString("233970423115425145544350131142039591210", 10)
	b3 := big.NewInt(727)
	o3, _ := new(big.Int).SetString("233970423115425145545378039958152057148", 10)
	priv := big.NewInt(705485)

	curve1 := NewCurve(a, b1, p, o1, gx, gy)
	curve2 := NewCurve(a, b2, p, o2, gx, gy)
	curve3 := NewCurve(a, b3, p, o3, gx, gy)
	smallOrderCurves := []shortWeierstrassCurve{
		curve1, curve2, curve3,
	}

	oracle := func(x, y *big.Int) (*big.Int, *big.Int) {
		a, b := curve.ScalarMult(x, y, priv.Bytes())
		return a, b
	}

	ind, _, err := curve.pohligHellmanOnline(smallOrderCurves, oracle)
	if err != nil {
		t.Errorf("%v", err)
	}
	if ind.Cmp(priv) != 0 {
		t.Errorf("failed to recover private key")
	}
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

	zx, zy := curve.Add(x0, y0, gx, gy)

	if !curve.PointEquals(zx, zy, zero, one) {
		t.Errorf("adding the inverse did not result in 0")
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

func TestScalarBaseMult(t *testing.T) {

	gx2, gy2 := curve.ScalarBaseMult(two.Bytes())
	gxd, gyd := curve.Double(gx, gy)

	if !curve.PointEquals(gx2, gy2, gxd, gyd) {
		t.Errorf("scalar base multiplication by 2 did not result in the same point as doubling")
		return
	}

}

func TestIsOnCurve(t *testing.T) {

	if !curve.IsOnCurve(gx, gy) {
		t.Errorf("IsOnCurve failed to identify the base point on the curve")
		return
	}

	gx1 := new(big.Int).Add(gx, one)

	if curve.IsOnCurve(gx1, gy) {
		t.Errorf("IsOnCurve failed to identify an off-curve point")
		return
	}
}

func TestOrderOfBasePoint(t *testing.T) {

	x, y := curve.ScalarBaseMult(order.Bytes())
	if !curve.PointEquals(x, y, zero, one) {
		t.Errorf("multiplying our base point by the order of the curve did not return 0")
		t.Errorf("result: (%d, %d)\n", x, y)
	}
}

func TestPreComputedScalarMult(t *testing.T) {

	gxh, gyh := curve.ScalarBaseMult(hund.Bytes())

	if !curve.PointEquals(gxh, gyh, hundX, hundY) {
		t.Errorf("point did not match the pre-computed point")
		t.Logf("result: (%d, %d)\n", gxh, gyh)
		return
	}

	gxo, gyo := curve.ScalarBaseMult(order.Bytes())
	if !curve.PointEquals(gxo, gyo, zero, one) {
		t.Errorf("scalar multiplication by the order did not return 0")
		t.Logf("result: (%d, %d)\n", gxo, gyo)
		return
	}
}

func TestECDH(t *testing.T) {

	bobK, _ := new(big.Int).SetString("15", 10)
	aliceK, _ := new(big.Int).SetString("16", 10)

	bobPubX, bobPubY := curve.ScalarBaseMult(bobK.Bytes())
	alicePubX, alicePubY := curve.ScalarBaseMult(aliceK.Bytes())

	aliceSharedX, aliceSharedY := curve.ScalarMult(bobPubX, bobPubY, aliceK.Bytes())
	bobSharedX, bobSharedY := curve.ScalarMult(alicePubX, alicePubY, bobK.Bytes())

	if !curve.PointEquals(aliceSharedX, aliceSharedY, bobSharedX, bobSharedY) {
		t.Errorf("shared secrets did not match")
		return
	}

}

func TestRandomPoint(t *testing.T) {

	for i := 0; i < 100; i++ {
		rx, ry := curve.randomPoint()
		if !curve.IsOnCurve(rx, ry) {
			t.Errorf("generated random point was not on the curve")
			return
		}
	}
}
