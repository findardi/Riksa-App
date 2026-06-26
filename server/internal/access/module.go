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

func NewModule(pool *pgxpool.Pool, verifier middleware.TokenVerifier, mail service.MailService, asvc service.AuthService, token service.Tokenizer, webURL string) *Module {
	r := repository.New(pool)
	s := service.NewAccessService(r, mail, asvc, token, webURL)
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
		r.Get("/permissions", m.handler.GetPermissions)

		r.Route("/workspaces/{workspaceID}", func(r chi.Router) {
			r.Route("/roles", func(r chi.Router) {
				r.Post("/", m.handler.CreateRole)
				r.Get("/", m.handler.GetRoles)
				r.Get("/{roleID}", m.handler.GetRole)
				r.Put("/{roleID}", m.handler.UpdateRole)
				r.Delete("/{roleID}", m.handler.DeleteRole)
			})

			r.Route("/members", func(r chi.Router) {
				r.Post("/", m.handler.AddMember)
				r.Get("/", m.handler.GetMembers)
				r.Get("/{memberID}", m.handler.GetMember)
				r.Put("/{memberID}", m.handler.UpdateMember)
				r.Delete("/{memberID}", m.handler.DeleteMember)
			})

			r.Route("/invitations", func(r chi.Router) {
				r.Post("/", m.handler.AddMembers)
				r.Get("/", m.handler.GetInvitations)
				r.Post("/{invitationID}/resend", m.handler.ResendInvitation)
				r.Post("/{invitationID}/revoke", m.handler.RevokeInvitation)
			})

			r.Route("/groups", func(r chi.Router) {
				r.Post("/", m.handler.CreateGroup)
				r.Get("/", m.handler.GetGroups)
				r.Get("/{groupID}", m.handler.GetGroup)
				r.Put("/{groupID}", m.handler.UpdateGroup)
				r.Delete("/{groupID}", m.handler.DeleteGroup)
				r.Post("/{groupID}/assign", m.handler.AssignMember)
				r.Delete("/{groupID}/unassign/{memberID}", m.handler.UnassignMember)
			})
		})

		r.Route("/member", func(r chi.Router) {
			r.Post("/{workspaceID}", m.handler.AddMember)
			r.Post("/{workspaceID}/invite", m.handler.AddMembers)
			r.Get("/{workspaceID}/invite", m.handler.GetInvitations)
			r.Get("/{workspaceID}", m.handler.GetMembers)
			r.Get("/{workspaceID}/{memberID}", m.handler.GetMember)
			r.Put("/{workspaceID}/{memberID}", m.handler.UpdateMember)
			r.Delete("/{workspaceID}/{memberID}", m.handler.DeleteMember)
		})
	})
}
