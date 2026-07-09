package content

import (
	"context"
	"errors"

	accessrepo "github.com/findardi/Wadi/server/internal/access/repository"
	accessdb "github.com/findardi/Wadi/server/internal/access/repository/sqlc"
	auth "github.com/findardi/Wadi/server/internal/auth/repository"
	"github.com/findardi/Wadi/server/internal/content/handler"
	"github.com/findardi/Wadi/server/internal/content/repository"
	"github.com/findardi/Wadi/server/internal/content/service"
	"github.com/findardi/Wadi/server/internal/platform/middleware"
	"github.com/findardi/Wadi/server/internal/platform/permission"
	"github.com/findardi/Wadi/server/internal/platform/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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
	handler    *handler.ContentHandler
	mw         *middleware.Middleware
	accessRepo *accessrepo.Repository
	storage    storage.Storage
}

func NewModule(pool *pgxpool.Pool, verifier middleware.TokenVerifier, store storage.Storage) *Module {
	r := repository.New(pool)
	s := service.NewContentService(r, store)
	h := handler.NewContentHandler(s)

	mw := middleware.New(verifier, userStatusReader{repo: auth.New(pool)}, nil)

	return &Module{
		handler:    h,
		mw:         mw,
		accessRepo: accessrepo.New(pool),
		storage:    store,
	}
}

func (m *Module) workspaceMember(ctx context.Context, workspaceID, userID string) (*middleware.Membership, error) {
	var wID, uID pgtype.UUID

	if err := wID.Scan(workspaceID); err != nil {
		return nil, middleware.ErrResourceNotFound
	}
	if err := uID.Scan(userID); err != nil {
		return nil, middleware.ErrResourceNotFound
	}

	row, err := m.accessRepo.GetMembershipWithPermissions(ctx, accessdb.GetMembershipWithPermissionsParams{
		WorkspaceID: wID,
		UserID:      uID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, middleware.ErrResourceNotFound
	}
	if err != nil {
		return nil, err
	}

	return &middleware.Membership{
		Role:        row.RoleName,
		Permissions: row.Permissions,
		Status:      row.Status,
	}, nil
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/content", func(r chi.Router) {
		r.Use(m.mw.RequireAuth)
		r.Use(m.mw.RequireActive)

		r.Route("/workspaces/{workspaceID}/folders", func(r chi.Router) {
			r.Use(m.mw.RequireMember("workspaceID", m.workspaceMember))

			r.With(m.mw.RequirePermission(permission.PermFolderView)).Get("/", m.handler.GetFoldersTree)
			r.With(m.mw.RequirePermission(permission.PermFolderCreate)).Post("/", m.handler.CreateFolder)
			r.With(m.mw.RequirePermission(permission.PermFolderEdit)).Put("/{folderID}", m.handler.RenameFolder)
			r.With(m.mw.RequirePermission(permission.PermFolderEdit)).Patch("/{folderID}/move", m.handler.MoveFolder)
			r.With(m.mw.RequirePermission(permission.PermFolderDelete)).Delete("/{folderID}", m.handler.DeleteFolder)
		})
	})
}
