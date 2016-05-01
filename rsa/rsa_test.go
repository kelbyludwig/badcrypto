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
