package big

import (
	"math/big"
	"testing"
)

func TestMontgomeryReduction(t *testing.T) {
	m := big.NewInt(187)
	T := big.NewInt(563)

	answer := big.NewInt(19)
	result := MontgomeryReduction(m, T)

	if answer.Cmp(result) != 0 {
		t.Errorf("montgomery reduction test failed\n")
		return
	}
}
