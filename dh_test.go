package gopenid

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	params = Params{
		P: new(big.Int).SetBytes([]byte{0xdc, 0xf9, 0x3a, 0x0b, 0x88, 0x39, 0x72, 0xec, 0x0e, 0x19, 0x98, 0x9a, 0xc5, 0xa2, 0xce, 0x31, 0x0e, 0x1d, 0x37, 0x71, 0x7e, 0x8d, 0x95, 0x71, 0xbb, 0x76, 0x23, 0x73, 0x18, 0x66, 0xe6, 0x1e, 0xf7, 0x5a, 0x2e, 0x27, 0x89, 0x8b, 0x05, 0x7f, 0x98, 0x91, 0xc2, 0xe2, 0x7a, 0x63, 0x9c, 0x3f, 0x29, 0xb6, 0x08, 0x14, 0x58, 0x1c, 0xd3, 0xb2, 0xca, 0x39, 0x86, 0xd2, 0x68, 0x37, 0x05, 0x57, 0x7d, 0x45, 0xc2, 0xe7, 0xe5, 0x2d, 0xc8, 0x1c, 0x7a, 0x17, 0x18, 0x76, 0xe5, 0xce, 0xa7, 0x4b, 0x14, 0x48, 0xbf, 0xdf, 0xaf, 0x18, 0x82, 0x8e, 0xfd, 0x25, 0x19, 0xf1, 0x4e, 0x45, 0xe3, 0x82, 0x66, 0x34, 0xaf, 0x19, 0x49, 0xe5, 0xb5, 0x35, 0xcc, 0x82, 0x9a, 0x48, 0x3b, 0x8a, 0x76, 0x22, 0x3e, 0x5d, 0x49, 0x0a, 0x25, 0x7f, 0x05, 0xbd, 0xff, 0x16, 0xf2, 0xfb, 0x22, 0xc5, 0x83, 0xab}),
		G: new(big.Int).SetInt64(2),
	}
)

func TestSharedSecret(t *testing.T) {
	a, err := GenerateKey(rand.Reader, 1024, params)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	b, err := GenerateKey(rand.Reader, 1024, params)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, a.SharedSecret(b.PublicKey), b.SharedSecret(a.PublicKey))
}
