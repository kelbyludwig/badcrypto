package dsa

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"math/big"
	"os"
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

//TestDSANonceReuse is a test for Cryptopals Set 6 Challenge 44
func TestDSANonceReuse(t *testing.T) {

	priv, err := GenerateKey()

	if err != nil {
		t.Errorf("error generating keypair")
		return
	}

	file, err := os.Open("44.txt")
	defer file.Close()
	if err != nil {
		t.Errorf("failed to open 44.txt")
		return
	}

	scanner := bufio.NewScanner(file)
	i := 0
	var msgs [][]byte
	var ss, rs, ms []*big.Int
	hex2bn := func(s string) *big.Int {
		b, _ := new(big.Int).SetString(s, 16)
		return b
	}
	dec2bn := func(s string) *big.Int {
		b, _ := new(big.Int).SetString(s, 10)
		return b
	}
	for scanner.Scan() {
		line := scanner.Text()
		i += 1
		switch i % 4 {
		case 1:
			msg := line[5:]
			msgs = append(msgs, []byte(msg))
		case 2:
			s := line[3:]
			ss = append(ss, dec2bn(s))
		case 3:
			r := line[3:]
			rs = append(rs, dec2bn(r))
		case 0:
			m := line[3:]
			ms = append(ms, hex2bn(m))
		}
	}

	if len(ss) != len(rs) || len(rs) != len(ms) {
		t.Errorf("lengths of signatures read from file were wrong")
		return
	}

	for i, x := range ms {
		for j, y := range ms {
			//messages were the same. lets skip.
			if x.Cmp(y) == 0 {
				continue
			}
			sinv := new(big.Int).Sub(ss[i], ss[j])
			sinv = sinv.ModInverse(sinv, priv.PublicKey.Q)
			k := new(big.Int).Sub(ms[i], ms[j])
			k = k.Mul(k, sinv)
			k = k.Mod(k, priv.PublicKey.Q)
			x := RecoverPrivateKeyFromSubKey(msgs[i], rs[i], ss[i], k, priv.PublicKey)
			xHex := fmt.Sprintf("%x", x)
			xDigest := sha1.Sum([]byte(xHex))
			xHex = fmt.Sprintf("%x", xDigest)
			if xHex == "ca8f6f7c66fa362d40760d135b763eb8527d3d52" {
				t.Logf("Recovered private key: %x\n", x)
				return
			}
		}
	}
}
