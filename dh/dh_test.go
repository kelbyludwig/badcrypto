package dh

import (
	"math/big"
	"testing"
)

func TestRandomSubgroupElement(t *testing.T) {

	tests := []struct {
		prime, factor, expected *big.Int
	}{
		{big.NewInt(7), big.NewInt(3), big.NewInt(4)},
		{big.NewInt(7), big.NewInt(2), big.NewInt(6)},
		{big.NewInt(11), big.NewInt(2), big.NewInt(10)},
		{big.NewInt(11), big.NewInt(5), big.NewInt(4)},
	}
	for i, test := range tests {
		result := RandomSubgroupElement(test.prime, test.factor)
		if test.expected.Cmp(result) != 0 {
			t.Errorf("expected subgroup generator not returned: %d != %d\n", i, test.expected, result)
			return
		}
	}
}

func TestComputeIndexWithinRange(t *testing.T) {
	tests := []struct {
		elem, gen, mod, min, max, expected *big.Int
	}{
		{big.NewInt(3), big.NewInt(7), big.NewInt(11), big.NewInt(1), big.NewInt(6), big.NewInt(4)},
		{big.NewInt(100), big.NewInt(2), big.NewInt(101), big.NewInt(1), big.NewInt(100), big.NewInt(50)},
		{big.NewInt(1572), big.NewInt(2), big.NewInt(3307), big.NewInt(700), big.NewInt(800), big.NewInt(789)},
	}
	for _, te := range tests {
		index, err := ComputeIndexWithinRange(te.elem, te.gen, te.mod, te.min, te.max)
		if err != nil {
			t.Errorf("index not recovered")
			return
		}
		if te.expected.Cmp(index) != 0 {
			t.Errorf("expected index not returned: %d != %d", index, te.expected)
			return
		}
	}
}
