package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/findardi/Wadi/server/internal/platform/response"
	"github.com/findardi/Wadi/server/internal/platform/token"
)

type TokenVerifier interface {
	VerifyToken(tokenString string) (*token.JwtClaims, error)
}

type StatusReader interface {
	UserStatus(ctx context.Context, userID string) (string, error)
}

type ctxKey string

const claimsKey ctxKey = "auth_claims"

type Middleware struct {
	verifier TokenVerifier
	status   StatusReader
}

func New(verifier TokenVerifier, status StatusReader) *Middleware {
	return &Middleware{
		verifier: verifier,
		status:   status,
	}
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(w, http.StatusUnauthorized, "missing or invalid authorization", nil)
			return
		}

		claims, err := m.verifier.VerifyToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) RequireActive(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}

		status, err := m.status.UserStatus(r.Context(), claims.ID)
		if err != nil {
			log.Printf("require active internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		if status != "active" {
			response.Error(w, http.StatusForbidden, "account not active", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ClaimsFromContext(ctx context.Context) (*token.JwtClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(*token.JwtClaims)
	return claims, ok
}
