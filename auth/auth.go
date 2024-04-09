package auth

import (
	"crypto/ed25519"
	"github.com/zorotocol/zone/errorutils"
	"github.com/zorotocol/zone/pb"
	"google.golang.org/protobuf/proto"
	"io"
)

func Generate(seed io.Reader) ([]byte, error) {
	_, key, err := ed25519.GenerateKey(seed)
	return key, err
}
func Derive(key []byte) []byte {
	return key[32:]
}
func Sign(key []byte, token *pb.Token) {
	token.Signature = nil
	token.Signature = ed25519.Sign(key, errorutils.Must(proto.Marshal(token)))
}
func Verify(token *pb.Token, pubs [][]byte) int {
	for i, pub := range pubs {
		if verify(token, pub) {
			return i
		}
	}
	return -1
}

func verify(token *pb.Token, pub []byte) bool {
	sig := token.GetSignature()
	token.Signature = nil
	ok := ed25519.Verify(pub, errorutils.Must(proto.Marshal(token)), sig)
	token.Signature = sig
	return ok
}
