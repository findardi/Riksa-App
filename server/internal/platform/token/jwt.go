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
	Email    string
	Status   string
	Typ      string
}

func (g *Generator) CreateToken(claims JwtClaims, tokenType TokenType) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       claims.ID,
		"username": claims.Username,
		"status":   claims.Status,
		"email":    claims.Email,
		"typ":      string(tokenType),
		"exp":      time.Now().Add(TokenTTL[tokenType]).Unix(),
	})

	tokenString, err := token.SignedString(g.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (g *Generator) VerifyToken(tokenString string) (*JwtClaims, error) {
	t, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { // cegah alg-confusion
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return g.secret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	id, _ := claims["id"].(string)
	username, _ := claims["username"].(string)
	status, _ := claims["status"].(string)
	email, _ := claims["email"].(string)
	typ, _ := claims["typ"].(string)
	return &JwtClaims{
		ID:       id,
		Username: username,
		Email:    email,
		Status:   status,
		Typ:      typ,
	}, nil
}
