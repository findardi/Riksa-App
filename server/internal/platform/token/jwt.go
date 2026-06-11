package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenLogin   string = "token_login"
	TokenRefresh string = "token_refresh"
)

var TokenTTL = map[TokenType]time.Duration{
	TokenType(TokenLogin):   15 * time.Minute,
	TokenType(TokenRefresh): 24 * time.Hour,
}

type Generator struct {
	secret []byte
}

func New(secret string) *Generator {
	return &Generator{
		secret: []byte(secret),
	}
}

type JwtClaims struct {
	ID       string
	Username string
	Status   string
}

func (g *Generator) CreateToken(claims JwtClaims, tokenType TokenType) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       claims.ID,
		"username": claims.Username,
		"status":   claims.Status,
		"exp":      time.Now().Add(TokenTTL[tokenType]).Unix(),
	})

	tokenString, err := token.SignedString(g.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (g *Generator) VerifyingToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return g.secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
