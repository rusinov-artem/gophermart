package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

type TokenGenerator struct {
}

func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

func (t *TokenGenerator) Generate() string {
	b := make([]byte, 50)
	_, _ = rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
