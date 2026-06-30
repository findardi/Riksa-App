package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/findardi/Wadi/server/internal/platform/token"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func reqWithClaims(param, val, userID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rctx := chi.NewRouteContext()
	if param != "" {
		rctx.URLParams.Add(param, val)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)

	if userID != "" {
		ctx = context.WithValue(ctx, claimsKey, &token.JwtClaims{ID: userID, Typ: token.TokenLogin})
	}

	return req.WithContext(ctx)
}

func spyHandler(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

func TestRequirePermission(t *testing.T) {
	m := New(nil, nil, nil)

	t.Run("grants when permission present", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), membershipKey, &Membership{
			Role:        "guest",
			Permissions: []string{"member:view", "role:view"},
		}))

		rec := httptest.NewRecorder()
		m.RequirePermission("member:view")(spyHandler(&called)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, called)
	})

	t.Run("forbids when permission absent", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), membershipKey, &Membership{
			Role:        "guest",
			Permissions: []string{"member:view"},
		}))

		rec := httptest.NewRecorder()
		m.RequirePermission("member:delete")(spyHandler(&called)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.False(t, called)
	})

	t.Run("fails closed when membership missing in context", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()
		m.RequirePermission("member:view")(spyHandler(&called)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.False(t, called)
	})
}

func TestRequireMember(t *testing.T) {
	m := New(nil, nil, nil)

	t.Run("unauthorized when claims missing", func(t *testing.T) {
		called := false
		resolver := func(ctx context.Context, workspaceID, userID string) (*Membership, error) {
			return &Membership{}, nil
		}
		req := reqWithClaims("workspaceID", "ws1", "")

		rec := httptest.NewRecorder()
		m.RequireMember("workspaceID", resolver)(spyHandler(&called)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.False(t, called)
	})

	t.Run("forbids non-member (ErrResourceNotFound)", func(t *testing.T) {
		called := false
		resolver := func(ctx context.Context, workspaceID, userID string) (*Membership, error) {
			return nil, ErrResourceNotFound
		}
		req := reqWithClaims("workspaceID", "ws1", "user1")

		rec := httptest.NewRecorder()
		m.RequireMember("workspaceID", resolver)(spyHandler(&called)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.False(t, called)
	})

	t.Run("internal error on generic resolver failure", func(t *testing.T) {
		called := false
		resolver := func(ctx context.Context, workspaceID, userID string) (*Membership, error) {
			return nil, errors.New("db down")
		}
		req := reqWithClaims("workspaceID", "ws1", "user1")

		rec := httptest.NewRecorder()
		m.RequireMember("workspaceID", resolver)(spyHandler(&called)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.False(t, called)
	})

	t.Run("passes and stashes membership in context", func(t *testing.T) {
		want := &Membership{Role: "owner", Permissions: []string{"member:delete"}, Status: "active"}
		resolver := func(ctx context.Context, workspaceID, userID string) (*Membership, error) {
			assert.Equal(t, "ws1", workspaceID)
			assert.Equal(t, "user1", userID)
			return want, nil
		}

		var got *Membership
		var ok bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got, ok = MembershipFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})
		req := reqWithClaims("workspaceID", "ws1", "user1")

		rec := httptest.NewRecorder()
		m.RequireMember("workspaceID", resolver)(next).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})
}
