package otp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Generator struct {
	secret []byte
}

func New(secret string) *Generator {
	return &Generator{
		secret: []byte(secret),
	}
}

func (g *Generator) Generate() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%06d", n.Int64())
}

func (g *Generator) Hash(code string) string {
	mac := hmac.New(sha256.New, g.secret)
	mac.Write([]byte(code))
	return hex.EncodeToString(mac.Sum(nil))
}

func (g *Generator) Compare(hash, code string) bool {
	return hmac.Equal([]byte(hash), []byte(g.Hash(code)))
}

func (g *Generator) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
