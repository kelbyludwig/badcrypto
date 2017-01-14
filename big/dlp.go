package big

import (
	"fmt"
	"math/big"
	"math/rand"
)

var IndexNotRecoveredErr error = fmt.Errorf("index not found")

//PohligHellmanOnline implements subgroup confinement against a "online"
//oracle. The oracle function should take in a group element `g` return `h =
//g^x (mod modulus)`. PohligHellmanOnline will generate small order groups and
//use the oracle to recover as much of the private key `x` as it can.
func PohligHellmanOnline(modulus *big.Int, oracle func(g *big.Int) (h *big.Int)) (index *big.Int, err error) {
	mmo := new(big.Int).Sub(modulus, one)
	factors, _ := Factor(mmo, 65536)
	indices := make([]*big.Int, 0)
	moduli := make([]*big.Int, 0)
	for factor, _ := range factors {
		primeFactor := big.NewInt(int64(factor))
		sGen := RandomSubgroupElement(modulus, primeFactor)
		sElem := oracle(sGen)
		ind, err := ComputeIndexWithinRange(sElem, sGen, modulus, zero, primeFactor)
		if err != nil {
			continue
		}
		indices = append(indices, ind)
		moduli = append(moduli, primeFactor)
	}

	if len(indices) == 0 || len(moduli) == 0 {
		err = IndexNotRecoveredErr
		return nil, err
	}

	index, err = CRT(indices, moduli)
	if err != nil {
		err = IndexNotRecoveredErr
		return nil, err
	}

	return index, nil
}

//PohligHellman will solve for the index of `elem` using the generator `gen`
//for the group of order `order`.
func PohligHellman(elem, gen, modulus, order *big.Int) (index *big.Int, err error) {

	ord := new(big.Int).SetBytes(order.Bytes())
	//TODO(kkl): Loop over Factor (I will need a FactorRange function first) and the `rest` return value here to
	//           intelligently only factor what we need to recover the index.
	factors, _ := Factor(ord, 65536)
	indices := make([]*big.Int, 0)
	moduli := make([]*big.Int, 0)

	//mmo := new(big.Int).Sub(modulus, one)
	//TODO(kkl): Update this to cover repeat prime factors.
	for factor, _ := range factors {
		primeFactor := big.NewInt(int64(factor))
		exp := new(big.Int).Div(ord, primeFactor)
		sElem := new(big.Int).Exp(elem, exp, modulus)
		sGen := new(big.Int).Exp(gen, exp, modulus)
		ind, err := ComputeIndexWithinRange(sElem, sGen, modulus, zero, primeFactor)
		if err != nil {
			continue
		}
		indices = append(indices, ind)
		moduli = append(moduli, primeFactor)
	}

	if len(indices) == 0 || len(moduli) == 0 {
		err = IndexNotRecoveredErr
		return nil, err
	}

	index, err = CRT(indices, moduli)
	if err != nil {
		err = IndexNotRecoveredErr
		return nil, err
	}

	return index, nil
}

//ComputeIndexWithinRange will solve for x in the equation gen^x = elem = (mod
//modulus). If the index does not fall within the specified range, this
//function will return an error.
func ComputeIndexWithinRange(elem, gen, modulus, min, max *big.Int) (index *big.Int, err error) {

	el := big.NewInt(1)
	index = new(big.Int).SetBytes(min.Bytes())
	for index.Cmp(max) != 1 {
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

//ComputeIndex will solve for x in the quation gen^x = elem (mod order). If
//the index is not found, this function will return an error.
func ComputeIndex(elem, gen, modulus *big.Int) (index *big.Int, err error) {
	two := big.NewInt(2)
	return ComputeIndexWithinRange(elem, gen, modulus, two, modulus)
}

//CRT will return the result of the chinese remainder thereom applied to the
//supplied residues and respective moduli.
func CRT(a, moduli []*big.Int) (*big.Int, error) {
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

// SqrtBig returns floor(sqrt(n)). It panics on n < 0.
// Source: https://github.com/cznic/mathutil/blob/master/mathutil.go#L151
func SqrtBig(n *big.Int) (x *big.Int) {
	switch n.Sign() {
	case -1:
		panic(-1)
	case 0:
		return big.NewInt(0)
	}

	var px, nx big.Int
	x = big.NewInt(0)
	x.SetBit(x, n.BitLen()/2+1, 1)
	for {
		nx.Rsh(nx.Add(x, nx.Div(n, x)), 1)
		if nx.Cmp(x) == 0 || nx.Cmp(&px) == 0 {
			break
		}
		px.Set(x)
		x.Set(&nx)
	}
	return

}

//BSGS uses Shank's Baby-Step Giant-Step algorithm to compute the discrete log
//of `elem`.
func BSGS(elem, gen, modulus *big.Int) (index *big.Int, err error) {

	m := SqrtBig(modulus)
	m = m.Add(m, big.NewInt(1))
	lookup := make(map[string]*big.Int)

	i := big.NewInt(1)
	res := big.NewInt(1)

	lookup["0"] = big.NewInt(1)
	for i.Cmp(m) != 1 {
		res = res.Mul(res, gen)
		res = res.Mod(res, modulus)
		if res.Cmp(zero) == 0 || res.Cmp(one) == 0 {
			break
		}
		lookup[res.String()] = new(big.Int).Set(i)
		i = i.Add(i, big.NewInt(1))
	}
	ginv := new(big.Int).ModInverse(gen, modulus)
	ginv = ginv.Exp(ginv, m, modulus)
	h := new(big.Int).Set(elem)
	i = big.NewInt(0)

	for i.Cmp(m) < 1 {

		j, ok := lookup[h.String()]
		if ok {
			index = new(big.Int).Set(i)
			index = index.Mul(index, m)
			index = index.Add(index, j)
			return
		}
		h = h.Mul(h, ginv)
		h = h.Mod(h, modulus)
		i = i.Add(i, big.NewInt(1))

	}

	return nil, IndexNotRecoveredErr

}
