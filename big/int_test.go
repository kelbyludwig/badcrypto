package big

import (
	"math/big"
	"testing"
)

func TestMontgomeryReduction(t *testing.T) {
	m := big.NewInt(187)
	T := big.NewInt(563)

	answer := big.NewInt(19)
	result := montgomeryReduction(m, T)

	if answer.Cmp(result) != 0 {
		t.Errorf("montgomery reduction test failed\n")
		return
	}
}

func TestMontgomeryMult(t *testing.T) {

	tests := []struct {
		x, y, m, answer *big.Int
	}{

		{big.NewInt(70), big.NewInt(91), big.NewInt(563), big.NewInt(177)},
		{big.NewInt(91), big.NewInt(70), big.NewInt(563), big.NewInt(177)},
		{big.NewInt(456), big.NewInt(123), big.NewInt(789), big.NewInt(69)},
		{big.NewInt(123), big.NewInt(456), big.NewInt(789), big.NewInt(69)},
		{big.NewInt(1234567), big.NewInt(890123), big.NewInt(999999999), big.NewInt(916482839)},
	}

	for i, a := range tests {
		result := montgomeryMul(a.x, a.y, a.m)

		if a.answer.Cmp(result) != 0 {
			t.Errorf("montgomery multiplication test failed\n")
			t.Logf("(%v) %v != %v\n", i, a.answer, result)
			return
		}
	}
}

func TestExpMont(t *testing.T) {
	tests := []struct {
		x, y, m, answer *big.Int
	}{
		{big.NewInt(2), big.NewInt(8), big.NewInt(9), big.NewInt(4)},
		{big.NewInt(70), big.NewInt(54), big.NewInt(17), big.NewInt(13)},
	}

	for i, a := range tests {
		result := ExpMont(a.x, a.y, a.m)

		if a.answer.Cmp(result) != 0 {
			t.Errorf("exp mont test failed\n")
			t.Logf("(%v) %v != %v\n", i, a.answer, result)
			return
		}
	}

}
