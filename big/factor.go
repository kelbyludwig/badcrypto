//int is supplmental big integer package that adds some extra bignum
//functionality that is not in the stdlib.
package big

import "math/big"

type Factors map[int64]int

var zero *big.Int = big.NewInt(0)
var one *big.Int = big.NewInt(1)

const Int64Max int64 = 9223372036854775807

//TODO(kkl): Functions to add here: FactorRange(num *big.Int, min, max int64), Factor(num *big.Int).
//			 It will be useful to start factoring from a partially factored `rest` starting from `min`.

//Factor takes in a bignum and returns a map of its factors. The result is a
//map from with prime factor keys and the prime factor's power as the value.
//Factor will only extract prime factors smaller than 9223372036854775807.  Any
//other supplied max value will be truncated. The return parameter `rest` is
//used return any remaining unfactored portions of the supplied integer.
func Factor(num *big.Int, max int64) (factors Factors, rest *big.Int) {
	rest = new(big.Int).SetBytes(num.Bytes())
	factors = make(map[int64]int)
	bigFact := new(big.Int)
	modResult := new(big.Int)

	if max > Int64Max {
		max = Int64Max
	}

	var fact int64
	for fact = 2; fact <= max; {
		bigFact = big.NewInt(fact)
		if rest.Cmp(one) == 0 {
			return
		}
		if modResult = modResult.Mod(rest, bigFact); modResult.Cmp(zero) == 0 {
			pow, ok := factors[fact]
			if !ok {
				pow = 0
			}
			factors[fact] = pow + 1
			rest = rest.Div(rest, bigFact)
		} else {
			fact++
		}
	}
	return
}
