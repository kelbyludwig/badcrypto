package main

import (
	"crypto/elliptic"
	"math/big"
)

type genericCurve struct {
	*CurveParams
}

func NewCurve(name string, bitsize int, gx, gy, b, n, p *big.Int) (curve elliptic.Curve) {
	curve = genericCurve{}
	curve.Name = name
	curve.BitSize = bitsize
	curve.Gx = gx
	curve.Gy = gy
	curve.B = b
	curve.N = n
	curve.P = p
	return curve
}

func (curve genericCurve) Params() *elliptic.CurveParams {
	return curve.CurveParams
}

func (curve genericCurve) IsOnCurve(x, y *big.Int) bool {

}

func (curve genericCurve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {

}
func (curve genericCurve) Double(x1, y1 *big.Int) (x, y *big.Int) {

}

func (curve genericCurve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {

}

func (curve genericCurve) ScalarBaseMult(k []byte) (x, y *big.Int) {

}
