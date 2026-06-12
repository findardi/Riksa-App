package auth

import (
	"github.com/findardi/Wadi/server/internal/auth/handler"
	"github.com/findardi/Wadi/server/internal/auth/repository"
	"github.com/findardi/Wadi/server/internal/auth/service"
	"github.com/findardi/Wadi/server/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *handler.AuthHandler
	mw      *middleware.Middleware
}

func NewModule(pool *pgxpool.Pool, otp service.OTPService, jwt service.JWTService) *Module {
	r := repository.New(pool)
	s := service.NewAuthService(r, otp, jwt)
	h := handler.NewAuthHandler(s)
	mw := middleware.New(jwt)

	return &Module{
		handler: h,
		mw:      mw,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		// publik
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
		})
	})
}
