package rsa

import (
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
	ciphertext := EncryptNoPaddingMontgomery(message, priv.PublicKey)
	ciphertextCorrect := EncryptNoPadding(message, priv.PublicKey)

	if string(ciphertext) != string(ciphertextCorrect) {
		t.Errorf("montgomery exp returns different result")
		return
	}
	log.Printf("private %+v\n", priv)
	log.Printf("public  %+v\n", priv.PublicKey)
	log.Printf("ciphertext %v\n", ciphertext)
	plaintext := DecryptNoPaddingMontgomery(ciphertext, priv)
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
