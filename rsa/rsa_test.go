package rsa

import (
	"fmt"
	"log"
	"math/big"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	priv, err := GenerateKey(64)
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	log.Printf("priv.D %v\n", priv.D)
	log.Printf("priv.Primes %v\n", priv.Primes)
	log.Printf("priv.Public %v\n", priv.PublicKey)
}

//TestEncryptDecrypt is a test for Cryptopals Set 5 Challenge 39
func TestEncryptDecrypt(t *testing.T) {
	priv, err := GenerateKey(128)
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	message := []byte("Cannnnnnnn do.")
	ciphertext := EncryptNoPadding(message, priv.PublicKey)
	log.Printf("private %+v\n", priv)
	log.Printf("public  %+v\n", priv.PublicKey)
	log.Printf("ciphertext %v\n", ciphertext)
	plaintext := DecryptNoPadding(ciphertext, priv)
	log.Printf("plaintext %v\n", plaintext)
	if string(plaintext) != string(message) {
		t.Errorf("decrypted message did not match plaintext")
		return
	}
}

//TestEncryptDecrypt is a test to verify Montgomery exponentiation works as planned.
func TestEncryptDecryptMontgomery(t *testing.T) {
	priv, err := GenerateKey(128)
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	message := []byte("Cannnnnnnn do.")
	ciphertext, _ := encryptNoPaddingMontgomery(message, priv.PublicKey)
	ciphertextCorrect := EncryptNoPadding(message, priv.PublicKey)

	if string(ciphertext) != string(ciphertextCorrect) {
		t.Errorf("montgomery exp returns different result")
		return
	}
	log.Printf("private %+v\n", priv)
	log.Printf("public  %+v\n", priv.PublicKey)
	log.Printf("ciphertext %v\n", ciphertext)
	plaintext, _ := decryptNoPaddingMontgomery(ciphertext, priv)
	log.Printf("plaintext %v\n", plaintext)
	if string(plaintext) != string(message) {
		t.Errorf("decrypted message did not match plaintext")
		return
	}
}

func TestChineseRemainderTheorem(t *testing.T) {
	//x = 2 mod 13
	//x = 3 mod 11
	//x = 4 mod 9
	//x = 652

	as := []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(4)}
	ms := []*big.Int{big.NewInt(13), big.NewInt(11), big.NewInt(9)}

	result, err := ChineseRemainderTheorem(as, ms)

	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

	if result.Cmp(big.NewInt(652)) != 0 {
		t.Errorf("crt result was incorrect. %v != 652", result)
		return
	}
}

//TestBroadcastAttack is a test for Cryptopals Set 5 Challenge 40
func TestBroadcastAttack(t *testing.T) {

	p1, err1 := GenerateKey(128)
	p2, err2 := GenerateKey(128)
	p3, err3 := GenerateKey(128)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("error in key generation\n")
		return
	}

	message := []byte("Oh yeah. Can do.")
	mb := new(big.Int).SetBytes(message)
	c1 := EncryptNoPadding(message, p1.PublicKey)
	c2 := EncryptNoPadding(message, p2.PublicKey)
	c3 := EncryptNoPadding(message, p3.PublicKey)

	cs := []*big.Int{new(big.Int).SetBytes(c1), new(big.Int).SetBytes(c2), new(big.Int).SetBytes(c3)}
	ns := []*big.Int{p1.PublicKey.N, p2.PublicKey.N, p3.PublicKey.N}

	result, err := ChineseRemainderTheorem(cs, ns)

	if err != nil {
		t.Errorf("error in crt")
		return
	}
	m := BigIntCubeRootFloor(result)
	if m.Cmp(mb) != 0 {
		t.Errorf("the cube root result was incorrect")
		return
	}
}

func TestMontgomeryExtraReductionsCount(t *testing.T) {

	priv, err := GenerateKey(512)

	if err != nil {
		t.Errorf("error in keygen\n")
		return
	}

	message := []byte("What is my purpose? Pass the butter.\n")

	ct, _ := encryptNoPaddingMontgomery(message, priv.PublicKey)
	_, ex1 := decryptNoPaddingMontgomery(ct, priv)
	t.Logf("Valid ciphertext did %v extra reductions\n", ex1)

	p := priv.Primes[0]
	q := priv.Primes[1]
	p1 := new(big.Int).Sub(p, big.NewInt(1))
	p2 := new(big.Int).Sub(p, big.NewInt(2))
	p3 := new(big.Int).Sub(p, big.NewInt(3))
	p4 := new(big.Int).Sub(p, big.NewInt(4))
	_, ex2 := decryptNoPaddingMontgomery(p.Bytes(), priv)
	_, ex3 := decryptNoPaddingMontgomery(q.Bytes(), priv)

	_, exm1 := decryptNoPaddingMontgomery(p1.Bytes(), priv)
	_, exm2 := decryptNoPaddingMontgomery(p2.Bytes(), priv)
	_, exm3 := decryptNoPaddingMontgomery(p3.Bytes(), priv)
	_, exm4 := decryptNoPaddingMontgomery(p4.Bytes(), priv)

	t.Logf("Ciphertext of p did %v extra reductions\n", ex2)
	t.Logf("Ciphertext of q did %v extra reductions\n", ex3)

	t.Logf("Ciphertext of p-4 did %v extra reductions\n", exm4)
	t.Logf("Ciphertext of p-3 did %v extra reductions\n", exm3)
	t.Logf("Ciphertext of p-2 did %v extra reductions\n", exm2)
	t.Logf("Ciphertext of p-1 did %v extra reductions\n", exm1)

}

//TestUnpaddedMessageRecovery is a test for Cryptopals Set 6 Challenge 41
func TestUnpaddedMessageRecovery(t *testing.T) {

	priv, err := GenerateKey(512)

	dupes := make(map[string]bool)
	decryptNoDupes := func(ct []byte) (pt []byte, err error) {
		if dupes[string(ct)] {
			return ct, fmt.Errorf("duplicate detected!")
		} else {
			dupes[string(ct)] = true
			return DecryptNoPadding(ct, priv), nil
		}
	}

	secretMessage := []byte("I hope no one reads this!")
	secretCiphertext := EncryptNoPadding(secretMessage, priv.PublicKey)

	sm1, err := decryptNoDupes(secretCiphertext)

	if err != nil {
		t.Errorf("failed to decrypt initial ciphertext")
		return
	}

	if string(sm1) != string(secretMessage) {
		t.Errorf("failed to properly decrypt ciphertext")
		return
	}

	//Now the attacker has the ciphertext! Oh no!
	_, err = decryptNoDupes(secretCiphertext)
	if err == nil {
		t.Errorf("the \"server\" failed to detect a dupe ciphertext")
		return
	}

	//create ((s**e mod n) c) mod n
	//where s is a random integer between 1 and n. here, s is 42.
	c := new(big.Int).SetBytes(secretCiphertext)
	s := big.NewInt(42)
	e := big.NewInt(priv.PublicKey.E)
	n := priv.PublicKey.N
	inverse := new(big.Int).ModInverse(s, n)
	s = s.Exp(s, e, n)
	s = s.Mul(s, c)
	s = s.Mod(s, n)

	sm3, err := decryptNoDupes(s.Bytes())

	if err != nil {
		t.Errorf("the server failed to decrypt our dupe ciphertext")
		return
	}

	plaintext := new(big.Int).SetBytes(sm3)
	plaintext = plaintext.Mul(plaintext, inverse)
	plaintext = plaintext.Mod(plaintext, n)

	if string(plaintext.Bytes()) != string(secretMessage) {
		t.Errorf("the decryption of our secret message was wrong")
		t.Logf("expected: %s\n", secretMessage)
		t.Logf("result:   %s\n", sm3)
		return
	}
}
