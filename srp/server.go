//srp implements the SRP PAKE partially based on RFC5054.
package srp

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"math/big"
)

type Server struct {
	username     string
	salt         []byte
	verifier     []byte
	ephemPublic  []byte
	ephemPrivate []byte
	pms          []byte
}

//pad returns a byte slice that is left-padded with zeros. This function is
//described in section 2.1 of RFC5054.
func pad(in []byte) []byte {
	l := len(in)
	buf := make([]byte, primeByteSize)
	copy(buf[primeByteSize-l:], in)
	return buf
}

//Server creates and initializes a Server.
func NewServer() (*Server, error) {
	ephemPrivate := make([]byte, primeByteSize)
	_, err := rand.Read(ephemPrivate)
	if err != nil {
		return nil, randError
	}
	ephemPrivateNum := new(big.Int).SetBytes(ephemPrivate)
	ephemPublic := new(big.Int).Exp(g, ephemPrivateNum, n)
	s := new(Server)
	s.ephemPublic = ephemPublic.Bytes()
	s.ephemPrivate = ephemPrivate
	return s, nil
}

//Register will store all authentication information for a given user.
func (s *Server) Register(username string, salt, verifier []byte) {
	s.username = username
	s.salt = salt
	s.verifier = verifier
}

//clientHelloResponse returns all the public information that the client does
//not store. clientHelloResponse in this toy example does not bother with
//sending DH params. It doesn't need to "lookup" a salt because a server only
//stores one user. The response to a ClientHello can be seen in RFC 5054 in
//section 2.2.
func (s *Server) clientHelloResponse() (salt, serverEphemPublic []byte) {
	return s.salt, s.ephemPublic
}

func (s *Server) ComputePremasterSecret(clientEphemPublic, u []byte) []byte {
	vn := new(big.Int).SetBytes(s.verifier)
	un := new(big.Int).SetBytes(u)
	An := new(big.Int).SetBytes(clientEphemPublic)
	bn := new(big.Int).SetBytes(s.ephemPrivate)

	pms := new(big.Int).Exp(vn, un, n)
	pms = pms.Mul(pms, An)
	pms = pms.Mod(pms, n)
	pms = pms.Exp(pms, bn, n)
	return pad(pms.Bytes())

}

func (s *Server) clientKeyExchangeResponse(clientEphemPublic []byte) {
	h := sha1.New()
	io.WriteString(h, string(pad(clientEphemPublic)))
	io.WriteString(h, string(pad(s.ephemPublic)))
	u := h.Sum(nil)
	s.pms = s.ComputePremasterSecret(clientEphemPublic, u[:])
}

func (s *Server) finishedResponse() error {
	//Do a MAC check to verify same pms
	return fmt.Errorf("srp: not implemented")
}
