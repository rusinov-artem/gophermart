package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

const hashLength = 50

type TokenGenerator struct {
}

func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

func (t *TokenGenerator) Generate() string {
	b := make([]byte, hashLength)
	_, _ = rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
