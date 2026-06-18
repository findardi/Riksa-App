package access

import (
	"context"

	"github.com/findardi/Wadi/server/internal/access/handler"
	"github.com/findardi/Wadi/server/internal/access/repository"
	"github.com/findardi/Wadi/server/internal/access/service"
	auth "github.com/findardi/Wadi/server/internal/auth/repository"
	"github.com/findardi/Wadi/server/internal/platform/middleware"
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
	handler *handler.AccessHandler
	mw      *middleware.Middleware
}

func NewModule(pool *pgxpool.Pool, verifier middleware.TokenVerifier, mail service.MailService) *Module {
	r := repository.New(pool)
	s := service.NewAccessService(r, mail)
	h := handler.NewAccessHandler(s)

	mw := middleware.New(verifier, userStatusReader{repo: auth.New(pool)}, nil)
	return &Module{
		handler: h,
		mw:      mw,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/access", func(r chi.Router) {
		r.Use(m.mw.RequireAuth)
		r.Use(m.mw.RequireActive)

		r.Route("/role", func(r chi.Router) {
			r.Post("/{workspaceID}", m.handler.CreateRole)
			r.Get("/{workspaceID}", m.handler.GetRoles)
			r.Get("/{workspaceID}/{roleID}", m.handler.GetRole)
			r.Put("/{workspaceID}/{roleID}", m.handler.UpdateRole)
			r.Delete("/{workspaceID}/{roleID}", m.handler.DeleteRole)
		})
	})
}
