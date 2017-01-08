package big

import (
	"math/big"
	"testing"
)

type factorTestcase struct {
	numToFactor *big.Int
	factors     Factors
	rest        *big.Int
	max         int64
}

func TestFactor(t *testing.T) {

	testCases := []factorTestcase{
		{big.NewInt(24), map[int64]int{2: 3, 3: 1}, big.NewInt(1), 64},
		{big.NewInt(7), map[int64]int{7: 1}, big.NewInt(1), 64},
		{big.NewInt(49), map[int64]int{7: 2}, big.NewInt(1), 64},
		{big.NewInt(469), map[int64]int{7: 1}, big.NewInt(67), 64},
	}

	for _, testCase := range testCases {
		factorResults, restResult := Factor(testCase.numToFactor, testCase.max)
		if restResult.Cmp(testCase.rest) != 0 {
			t.Errorf("rest results did not match expected rest")
			return
		}

		if len(factorResults) != len(testCase.factors) {
			t.Errorf("factor results did not have the same number of factors")
			return
		}
		for expectedKey, expectedVal := range testCase.factors {
			resultVal, ok := factorResults[expectedKey]
			if !ok {
				t.Errorf("factors did not include expected factor\n")
				return
			}
			if resultVal != expectedVal {
				t.Errorf("factor did not have expected power\n")
				return
			}
		}
	}

}
