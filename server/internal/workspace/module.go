package workspace

import (
	"context"

	auth "github.com/findardi/Wadi/server/internal/auth/repository"
	"github.com/findardi/Wadi/server/internal/platform/middleware"
	"github.com/findardi/Wadi/server/internal/workspace/handler"
	"github.com/findardi/Wadi/server/internal/workspace/repository"
	"github.com/findardi/Wadi/server/internal/workspace/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userStatusReader struct {
	repo *auth.Repository
}

func (s userStatusReader) UserStatus(ctx context.Context, userID string) (string, error) {
	var uid pgtype.UUID
	if err := uid.Scan(userID); err != nil {
		return "", err
	}

	user, err := s.repo.GetUserById(ctx, uid)
	if err != nil {
		return "", err
	}
	return user.Status, nil
}

type Module struct {
	handler *handler.WorkspaceHandler
	mw      *middleware.Middleware
}

func NewModule(pool *pgxpool.Pool, verifier middleware.TokenVerifier) *Module {
	r := repository.New(pool)
	s := service.NewWorkspaceService(r)
	h := handler.NewWorkspaceHandler(s)

	// RateLimit is unused for workspace routes, so no limiter is wired.
	mw := middleware.New(verifier, userStatusReader{repo: auth.New(pool)}, nil)

	return &Module{
		handler: h,
		mw:      mw,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/workspaces", func(r chi.Router) {
		r.Use(m.mw.RequireAuth)
		r.Use(m.mw.RequireActive)

		r.Post("/", m.handler.Create)
		r.Get("/", m.handler.GetWorkspaces)
		r.Post("/detail", m.handler.GetWorkspace)
		r.Put("/", m.handler.UpdateWorkspace)
		r.Patch("/status", m.handler.UpdateStatusWorkspace)
		r.Delete("/", m.handler.DeleteWorkspace)
	})
}
