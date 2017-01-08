package big

import (
	"fmt"
	"math/big"
	"math/rand"
)

var IndexNotRecoveredErr error = fmt.Errorf("index not found")

func PohligHellman(elem, gen, modulus *big.Int) (index *big.Int, err error) {

	ord := new(big.Int).Sub(modulus, big.NewInt(1))
	bigPow := new(big.Int)
	subElem := new(big.Int)
	subGen := new(big.Int)

	indices := make([]*big.Int, 0)
	moduli := make([]*big.Int, 0)

	factors, _ := Factor(ord, Int64Max)

	for factor, pow := range factors {
		primePower := new(big.Int).Exp(big.NewInt(factor), big.NewInt(int64(pow)), nil)
		bigPow = bigPow.Div(ord, primePower)
		subElem = subElem.Exp(elem, bigPow, modulus)
		subGen = subGen.Exp(gen, bigPow, modulus)
		index, err := ComputeIndex(subElem, subGen, modulus)
		if err != nil {
			return index, err
		}
		indices = append(indices, index)
		moduli = append(moduli, primePower)
	}

	index, err = CRT(indices, moduli)
	if err != nil {
		err = IndexNotRecoveredErr
		return
	}
	return index, err
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

//ComputeIndex will solve for x in the quation gen^x = elem (mod modudlus). If
//the index is not found, this function will return an error.
func ComputeIndex(elem, gen, modulus *big.Int) (index *big.Int, err error) {
	one := big.NewInt(1)
	return ComputeIndexWithinRange(elem, gen, modulus, one, modulus)
}

//CRT will return the result of the chinese remainder thereom applied to the
//supplied residues and respective moduli.
func CRT(a, moduli []*big.Int) (*big.Int, error) {
	var one = big.NewInt(1)
	p := new(big.Int).Set(moduli[0])
	for _, n1 := range moduli[1:] {
		p.Mul(p, n1)
	}
	var x, q, s, z big.Int
	for i, n1 := range moduli {
		q.Div(p, n1)
		z.GCD(nil, &s, n1, &q)
		if z.Cmp(one) != 0 {
			return nil, fmt.Errorf("%d not coprime", n1)
		}
		x.Add(&x, s.Mul(a[i], s.Mul(&s, &q)))
	}
	return x.Mod(&x, p), nil
}
