package auth

import (
	"context"
	"net/http"

	"github.com/findardi/Wadi/server/internal/auth/handler"
	"github.com/findardi/Wadi/server/internal/auth/repository"
	"github.com/findardi/Wadi/server/internal/auth/service"
	"github.com/findardi/Wadi/server/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *handler.AuthHandler
	mw      *middleware.Middleware
}

type userStatusReader struct {
	repo *repository.Repository
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

func NewModule(pool *pgxpool.Pool, otp service.OTPService, jwt service.JWTService, mail service.MailService) *Module {
	r := repository.New(pool)
	s := service.NewAuthService(r, otp, jwt, mail)
	h := handler.NewAuthHandler(s)
	mw := middleware.New(jwt, userStatusReader{repo: r})

	return &Module{
		handler: h,
		mw:      mw,
	}
}

func (m *Module) RequireActive(next http.Handler) http.Handler {
	return m.mw.RequireActive(next)
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		// public
		r.Post("/register", m.handler.Register)
		r.Post("/login", m.handler.Login)
		r.Post("/refresh", m.handler.RefreshToken)
		r.Post("/forgot-password", m.handler.ForgotPassword)
		r.Post("/reset-password", m.handler.ResetPassword)
		r.Post("/validation-otp", m.handler.CheckOTP)

		// protected
		r.Group(func(r chi.Router) {
			r.Use(m.mw.RequireAuth)
			r.Post("/resend-otp", m.handler.ResendOTP)
			r.Post("/verify-email", m.handler.VerifyAccount)
			r.Post("/logout", m.handler.Logout)
			r.Get("/me", m.handler.GetMe)
		})
	})
}
