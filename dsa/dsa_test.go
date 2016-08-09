package dsa

import (
	"testing"
)

func TestDSASign(t *testing.T) {

	message := []byte("i'm walking here!")
	priv, err := GenerateKey()
	if err != nil {
		t.Errorf("failed to generate key")
		return
	}

	r, s, err := Sign(message, priv)

	if err != nil {
		t.Errorf("failed to sign message")
		return
	}

	err = Verify(message, r, s, priv.PublicKey)

	if err != nil {
		t.Errorf("failed to verify signature")
		return
	}

}
