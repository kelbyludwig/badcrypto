package rsa

import (
	"log"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	priv, err := GenerateKey(512)
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	log.Printf("priv.D %v\n", priv.D)
	log.Printf("priv.Primes %v\n", priv.Primes)
	log.Printf("priv.Public %v\n", priv.PublicKey)
}

func TestEncryptDecrypt(t *testing.T) {
	priv, err := GenerateKey(512)
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	message := []byte("Cannnnnnnn do.")
	ciphertext := EncryptNoPadding(message, priv.PublicKey)
	log.Printf("ciphertext %v\n", ciphertext)
	plaintext := DecryptNoPadding(ciphertext, priv)
	log.Printf("plaintext %v\n", plaintext)
	if string(plaintext) != string(message) {
		t.Errorf("decrypted message did not match plaintext")
		return
	}
}
