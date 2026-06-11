package auth

import (
	"github.com/findardi/Wadi/server/internal/auth/handler"
	"github.com/findardi/Wadi/server/internal/auth/repository"
	"github.com/findardi/Wadi/server/internal/auth/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *handler.AuthHandler
}

func NewModule(pool *pgxpool.Pool, otp service.OTPService, jwt service.JWTService) *Module {
	r := repository.New(pool)
	s := service.NewAuthService(r, otp, jwt)
	h := handler.NewAuthHandler(s)

	return &Module{
		handler: h,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", m.handler.Register)
		r.Post("/login", m.handler.Login)
	})
}
