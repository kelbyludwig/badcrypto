//srp implements the SRP PAKE partially based on RFC5054.
package srp

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

//1024-bit SRP group modulus from Appendix A of RFC5054.
const primeString = "EEAF0AB9ADB38DD69C33F80AFA8FC5E86072618775FF3C0B9EA2314C9C256576D674DF7496EA81D3383B4813D692C6E0E0D5D8E250B98BE48E495C1D6089DAD15DC7D7B46154D6B6CE8EF4AD69B15D4982559B297BCF1885C529F566660E57EC68EDBC3C05726CC02FD4CBF4976EAA9AFD5138FE8376435B9FC61D2FC0EB06E3"

//The size of the modulus in bytes.
const primeByteSize = 128

var n *big.Int
var g *big.Int

var decodeError = fmt.Errorf("srp: unable to decode prime")
var randError = fmt.Errorf("srp: unable to read from urandom")
var connError = fmt.Errorf("srp: unable to connect to server")

type Client struct {
	server       *Server
	x            []byte
	ephemPrivate []byte
	ephemPublic  []byte
}

func init() {
	primeBytes, err := hex.DecodeString(primeString)
	if err != nil {
		panic(decodeError)
	}
	n = new(big.Int).SetBytes(primeBytes)
	g = big.NewInt(2)
}

//computePremasterSecret will compute the shared Premaster Secret between a
//client and server.
func (c *Client) ComputePremasterSecret(k, x, u, serverPublicKey []byte) []byte {

	kn := new(big.Int).SetBytes(k)
	xn := new(big.Int).SetBytes(x)
	un := new(big.Int).SetBytes(u)
	Bn := new(big.Int).SetBytes(serverPublicKey)
	an := new(big.Int).SetBytes(c.ephemPrivate)

	base := new(big.Int).Exp(g, xn, n)
	base = base.Mul(base, kn)
	base = base.Sub(Bn, base)
	base = base.Mod(base, n)

	exp := new(big.Int).Mul(un, xn)
	exp = exp.Add(an, exp)

	pmk := new(big.Int).Exp(base, exp, n)
	return pad(pmk.Bytes())

}

//New initializes a Client struct's public and private keys.
func NewClient() (*Client, error) {

	ephemPrivate := make([]byte, primeByteSize)
	_, err := rand.Read(ephemPrivate)
	if err != nil {
		return nil, randError
	}
	ephemPrivateNum := new(big.Int).SetBytes(ephemPrivate)
	ephemPublic := new(big.Int).Exp(g, ephemPrivateNum, n)

	c := new(Client)
	c.server, err = NewServer()
	if err != nil {
		return nil, err
	}
	c.ephemPrivate = ephemPrivate
	c.ephemPublic = ephemPublic.Bytes()
	return c, nil
}

//Register will do all necessary client-side setup and submit a registration
//request to the server. If client-side setup or server-side registration
//fails, an error will be returned.
func (c *Client) Register(username, password string) error {

	//For simplicity (avoids a RT) lets have the client generate the salt.
	salt := make([]byte, 16)
	_, err := rand.Read(salt)

	if err != nil {
		return randError
	}

	//Generate the verifier and register our identity with the server.
	h := sha1.New()
	innerHash := sha1.Sum([]byte(username + ":" + password))
	io.WriteString(h, string(salt))
	io.WriteString(h, string(innerHash[:]))
	x := h.Sum(nil)
	xN := new(big.Int).SetBytes(x[:])
	v := new(big.Int).Exp(g, xN, n)

	c.server.Register(username, salt, v.Bytes())
	return nil

}
