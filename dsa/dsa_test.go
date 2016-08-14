package dsa

import (
	"crypto/sha1"
	"fmt"
	"math/big"
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

//TestRecoverPrivateKeyFromSubKey is a test for Cryptopals Set 6 Challenge 43
func TestRecoverPrivateKeyFromSubKey(t *testing.T) {
	msg := []byte("For those that envy a MC it can be hazardous to your health\n" +
		"So be friendly, a matter of life and death, just like a etch-a-sketch\n")

	msgDigest := sha1.Sum(msg)
	msgDigestHex := fmt.Sprintf("%x", msgDigest)

	if msgDigestHex != "d2d0714f014a9784047eaeccf956520045c45265" {
		t.Errorf("message digests did not match")
		t.Logf("actual: %s\n", msgDigestHex)
		t.Logf("expect: d2d0714f014a9784047eaeccf956520045c45265\n")
		return
	}

	priv, err := GenerateKey()

	if err != nil {
		t.Errorf("error generating keypair")
		return
	}

	//Don't need a valid private key for this exercise so lets
	//just overwrite this keypair's Y. Its simpler than moving
	//all the DSA parameter generation/setup code here.
	priv.PublicKey.Y = new(big.Int).SetBytes(msgDigest[:])

	r, ok1 := new(big.Int).SetString("548099063082341131477253921760299949438196259240", 10)
	s, ok2 := new(big.Int).SetString("857042759984254168557880549501802188789837994940", 10)

	if !ok1 || !ok2 {
		t.Errorf("failed to generate signature")
		return
	}

	for i := 0; i <= 65536; i++ {
		k := big.NewInt(int64(i))
		x := RecoverPrivateKeyFromSubKey(msg, r, s, k, priv.PublicKey)
		xHex := fmt.Sprintf("%x", x)
		xDigest := sha1.Sum([]byte(xHex))
		xHex = fmt.Sprintf("%x", xDigest)
		if xHex == "0954edd5e0afe5542a4adf012611a91912a3ec16" {
			t.Logf("Recovered private key: %x\n", x)
			return
		}
	}
	t.Errorf("did not find matching private key")

}
