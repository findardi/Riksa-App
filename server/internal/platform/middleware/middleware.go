package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/findardi/Riksa-App/server/internal/platform/response"
	"github.com/findardi/Riksa-App/server/internal/platform/token"
	"github.com/go-chi/chi/v5"
)

type TokenVerifier interface {
	VerifyToken(tokenString string) (*token.JwtClaims, error)
}

type StatusReader interface {
	UserStatus(ctx context.Context, userID string) (string, error)
}

// ErrResourceNotFound is returned by an OwnerResolver when the resource does
// not exist; RequireOwner maps it to 404.
var ErrResourceNotFound = errors.New("resource not found")

// OwnerResolver resolves the owner (creator) user id of the resource identified
// by id. Each domain (workspace, folder, ...) supplies its own resolver so the
// middleware stays decoupled from any domain. Return ErrResourceNotFound when
// the resource is absent.
type OwnerResolver func(ctx context.Context, id string) (ownerID string, err error)

// Membership is the caller's resolved standing in a workspace: their role name,
// the flattened permission set of that role, and the member status. RequireMember
// loads it into the request context; RequirePermission reads it.
type Membership struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	Status      string   `json:"status"`
}

type MemberResolver func(ctx context.Context, workspaceID string, userID string) (*Membership, error)

type RateStore interface {
	Allow(key string, limit int, window time.Duration) (allowed bool, retryAfter time.Duration)
}

type KeyFunc func(r *http.Request) string

type RateConfig struct {
	Name   string
	Limit  int
	Window time.Duration
	Key    KeyFunc
}

type ctxKey string

const claimsKey ctxKey = "auth_claims"

const membershipKey ctxKey = "auth_membership"

type Middleware struct {
	verifier TokenVerifier
	status   StatusReader
	limiter  RateStore
}

func New(verifier TokenVerifier, status StatusReader, limiter RateStore) *Middleware {
	return &Middleware{
		verifier: verifier,
		status:   status,
		limiter:  limiter,
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

		// only access tokens may pass; reject anything not minted as token_login
		if claims.Typ != token.TokenLogin {
			response.Error(w, http.StatusUnauthorized, "invalid token type", nil)
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

func (m *Middleware) RequireOwner(param string, resolve OwnerResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}

			ownerID, err := resolve(r.Context(), chi.URLParam(r, param))
			switch {
			case errors.Is(err, ErrResourceNotFound):
				response.Error(w, http.StatusNotFound, "not found", nil)
			case err != nil:
				log.Printf("require owner internal error: %v", err)
				response.Error(w, http.StatusInternalServerError, "internal server error", nil)
			case ownerID != claims.ID:
				response.Error(w, http.StatusForbidden, "forbidden", nil)
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
}

func (m *Middleware) RequireMember(param string, resolver MemberResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}

			ms, err := resolver(r.Context(), chi.URLParam(r, param), claims.ID)
			switch {
			case errors.Is(err, ErrResourceNotFound):
				response.Error(w, http.StatusForbidden, "forbidden", nil)
			case err != nil:
				log.Printf("require member internal error: %v", err)
				response.Error(w, http.StatusInternalServerError, "internal server error", nil)
			default:
				ctx := context.WithValue(r.Context(), membershipKey, ms)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

func (m *Middleware) RequirePermission(perm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ms, ok := MembershipFromContext(r.Context())
			if !ok {
				log.Printf("require permission: membership missing in context (RequireMember not applied?)")
				response.Error(w, http.StatusInternalServerError, "internal server error", nil)
				return
			}

			for _, p := range ms.Permissions {
				if p == perm {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Error(w, http.StatusForbidden, "forbidden", nil)
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*token.JwtClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(*token.JwtClaims)
	return claims, ok
}

func MembershipFromContext(ctx context.Context) (*Membership, bool) {
	ms, ok := ctx.Value(membershipKey).(*Membership)
	return ms, ok
}

func (m *Middleware) RateLimit(cfg RateConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := ""
			if cfg.Key != nil {
				id = cfg.Key(r)
			}

			key := cfg.Name + "|" + clientIP(r) + "|" + id

			allowed, retryAfter := m.limiter.Allow(key, cfg.Limit, cfg.Window)
			if !allowed {
				secs := int(math.Ceil(retryAfter.Seconds()))
				if secs < 1 {
					secs = 1
				}
				w.Header().Set("Retry-After", strconv.Itoa(secs))
				response.Error(w, http.StatusTooManyRequests, "too many requests, please try again later", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func KeyFromClaims(r *http.Request) string {
	if claims, ok := ClaimsFromContext(r.Context()); ok {
		return strings.ToLower(claims.Email)
	}
	return ""
}

func KeyFromJSONField(field string) KeyFunc {
	return func(r *http.Request) string {
		if r.Body == nil {
			return ""
		}
		buf, err := io.ReadAll(io.LimitReader(r.Body, MaxBodyBytesPeek))
		r.Body = io.NopCloser(bytes.NewReader(buf))
		if err != nil {
			return ""
		}
		var body map[string]any
		if err := json.Unmarshal(buf, &body); err != nil {
			return ""
		}
		if v, ok := body[field].(string); ok {
			return strings.ToLower(strings.TrimSpace(v))
		}
		return ""
	}
}

const MaxBodyBytesPeek = 1 << 20

func ClientIP(r *http.Request) string {
	return clientIP(r)
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return strings.TrimSpace(xr)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
