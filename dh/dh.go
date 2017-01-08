package dh

import (
	"fmt"
	"math/big"
	"math/rand"
)

var IndexNotRecoveredErr error = fmt.Errorf("index not found")

//RandomSubgroupElement takes a prime and a factor of p-1 and returns an
//element of a subgroup of order `factor`.
func RandomSubgroupElement(prime, factor *big.Int) (elem *big.Int) {
	one := big.NewInt(1)
	rand := rand.New(rand.NewSource(99))
	for {
		elem = new(big.Int).Rand(rand, prime)
		pow := new(big.Int).Sub(prime, one)
		pow = pow.Div(pow, factor)
		elem = elem.Exp(elem, pow, prime)
		if elem.Cmp(one) != 0 {
			return
		}
	}
}

//ComputeIndexWithinRange will solve for x in the equation gen^x = elem = (mod
//modulus). If the index does not fall within the specified range, this
//function will return an error.
func ComputeIndexWithinRange(elem, gen, modulus, min, max *big.Int) (index *big.Int, err error) {
	el := big.NewInt(1)
	one := big.NewInt(1)
	index = new(big.Int).SetBytes(min.Bytes())
	for index.Cmp(max) != 0 {
		el = el.Exp(gen, index, modulus)
		if el.Cmp(elem) == 0 {
			return
		}
		index = index.Add(index, one)
	}
	return nil, IndexNotRecoveredErr
}

func ComputeIndex(elem, gen, modulus *big.Int) (index *big.Int, err error) {
	one := big.NewInt(1)
	return ComputeIndexWithinRange(elem, gen, modulus, one, modulus)
}
