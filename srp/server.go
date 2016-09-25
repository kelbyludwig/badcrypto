//srp implements the SRP PAKE partially based on RFC5054.
package srp

import (
	"crypto/rand"
	"math/big"
)

type Server struct {
	username     string
	salt         []byte
	verifier     []byte
	ephemPublic  []byte
	ephemPrivate []byte
	scramble     []byte
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
