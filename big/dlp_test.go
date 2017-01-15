package big

import (
	"math/big"
	"testing"
)

func TestCryptopals58(t *testing.T) {

	g, _ := new(big.Int).SetString("622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357", 10)
	p, _ := new(big.Int).SetString("11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623", 10)
	x := big.NewInt(705485)
	q, _ := new(big.Int).SetString("335062023296420808191071248367701059461", 10)
	A, _ := new(big.Int).SetString("7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119", 10)
	min := big.NewInt(0)
	max := big.NewInt(1048576)
	result, _ := Kangaroo(A, g, p, min, max)
	if result.Cmp(x) != 0 {
		t.Errorf("incorrect index returned for smaller case")
		return
	}

	if !testing.Short() {
		b, _ := new(big.Int).SetString("57329793620964792054131044619600637163", 10)
		B, _ := new(big.Int).SetString("386104828672368663930207367458554559828042669904904646699298598861116944033336766314203586720389684643394710343432816562435399603863591968684908781766626", 10)
		oracle := func(in *big.Int) *big.Int {
			return new(big.Int).Exp(in, b, p)
		}
		partialIndex, partialModulus, err := PohligHellmanOnline(p, oracle)
		if err != nil {
			t.Errorf("unexpected error occurred in Pohlig-Hellman step")
			return
		}
		if partialIndex.Cmp(b) == 0 {
			t.Errorf("we accidently recovered the whole index....")
			return
		}

		index, err := RecoverPartialIndex(B, g, p, q, partialIndex, partialModulus)
		if err != nil {
			t.Errorf("unexpected error occurred in partial index recovery step")
			return
		}
		if index.Cmp(b) != 0 {
			t.Errorf("full index did not match expected index")
			return
		}
	}

	if !testing.Short() {
		B, _ := new(big.Int).SetString("9388897478013399550694114614498790691034187453089355259602614074132918843899833277397448144245883225611726912025846772975325932794909655215329941809013733", 10)
		max = big.NewInt(1099511627776)
		y := big.NewInt(359579674340)
		result, _ = Kangaroo(B, g, p, min, max)
		if result.Cmp(y) != 0 {
			t.Errorf("incorrect index returned for larger case")
			return
		}
	}
}

func TestCryptopals57(t *testing.T) {

	g, _ := new(big.Int).SetString("4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143", 10)
	p, _ := new(big.Int).SetString("7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771", 10)
	x, _ := new(big.Int).SetString("751699345113817921351405264898819572", 10)
	A, _ := new(big.Int).SetString("982459683069415175513312739702735463971070011724878329469725424964365076122094000771127132784175677820400822159509270013047484755753176549654608169804860", 10)
	q, _ := new(big.Int).SetString("236234353446506858198510045061214171961", 10)

	oracler := func(mod, x *big.Int) func(*big.Int) *big.Int {
		return func(g *big.Int) *big.Int {
			return new(big.Int).Exp(g, x, mod)
		}
	}
	tests := []struct {
		elem, gen, mod, ord, expected *big.Int
	}{
		{big.NewInt(3), big.NewInt(7), big.NewInt(11), big.NewInt(10), big.NewInt(4)},
		{big.NewInt(1572), big.NewInt(2), big.NewInt(3307), big.NewInt(3306), big.NewInt(789)},
		{big.NewInt(298403), big.NewInt(2), big.NewInt(510529), big.NewInt(510528), big.NewInt(3500)},
		{A, g, p, q, x},
	}
	for _, te := range tests {
		oracle := oracler(te.mod, te.expected)
		result, _, err := PohligHellmanOnline(te.mod, oracle)
		if err != nil {
			t.Errorf("unexpected error occurred")
		}
		if result.Cmp(te.expected) != 0 {
			t.Errorf("incorrect result returned: %d != %d^%d %% %d == %d", result, te.gen, te.expected, te.mod, te.elem)
		}
	}
}

func TestKangaroo(t *testing.T) {

	tests := []struct {
		elem, gen, mod, min, max, expected *big.Int
	}{
		{big.NewInt(3), big.NewInt(7), big.NewInt(11), big.NewInt(2), big.NewInt(6), big.NewInt(4)},
		{big.NewInt(1572), big.NewInt(2), big.NewInt(3307), big.NewInt(600), big.NewInt(800), big.NewInt(789)},
		{big.NewInt(298403), big.NewInt(2), big.NewInt(510529), big.NewInt(3000), big.NewInt(3501), big.NewInt(3500)},
	}

	for _, te := range tests {
		result, err := Kangaroo(te.elem, te.gen, te.mod, te.min, te.max)
		if err != nil {
			t.Errorf("unexpected error occurred")
		}
		if result.Cmp(te.expected) != 0 {
			t.Errorf("incorrect result returned: %d != %d^%d %% %d == %d", result, te.gen, te.expected, te.mod, te.elem)
		}
	}

}

func TestPohligHellman(t *testing.T) {

	tests := []struct {
		elem, gen, mod, ord, expected *big.Int
	}{
		{big.NewInt(3), big.NewInt(7), big.NewInt(11), big.NewInt(10), big.NewInt(4)},
		{big.NewInt(1572), big.NewInt(2), big.NewInt(3307), big.NewInt(3306), big.NewInt(789)},
		{big.NewInt(298403), big.NewInt(2), big.NewInt(510529), big.NewInt(510528), big.NewInt(3500)},
	}

	for _, te := range tests {
		result, _, err := PohligHellman(te.elem, te.gen, te.mod, te.ord)
		if err != nil {
			t.Errorf("unexpected error occurred")
		}
		if result.Cmp(te.expected) != 0 {
			t.Errorf("incorrect result returned: %d != %d^%d %% %d == %d", result, te.gen, te.expected, te.mod, te.elem)
		}
	}

}

func TestCRT(t *testing.T) {

	tests := []struct {
		residues, moduli []*big.Int
		expected         *big.Int
	}{
		{
			[]*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(2)},
			[]*big.Int{big.NewInt(3), big.NewInt(5), big.NewInt(7)},
			big.NewInt(23),
		},
		{
			[]*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(5)},
			[]*big.Int{big.NewInt(3), big.NewInt(5), big.NewInt(11)},
			big.NewInt(38),
		},
	}

	result := new(big.Int)
	var err error
	for _, te := range tests {
		result, _, err = CRT(te.residues, te.moduli)
		if err != nil {
			t.Errorf("unexpected error occurred")
			return
		}
		if result.Cmp(te.expected) != 0 {
			t.Errorf("incorrect result returned")
			return
		}
	}
}

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

func TestBSGS(t *testing.T) {
	tests := []struct {
		elem, gen, mod, expected *big.Int
	}{
		{big.NewInt(6), big.NewInt(3), big.NewInt(31), big.NewInt(25)},
		{big.NewInt(3), big.NewInt(7), big.NewInt(11), big.NewInt(4)},
		{big.NewInt(100), big.NewInt(2), big.NewInt(101), big.NewInt(50)},
		{big.NewInt(1572), big.NewInt(2), big.NewInt(3307), big.NewInt(789)},
	}
	for _, te := range tests {
		index, err := BSGS(te.elem, te.gen, te.mod)
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
