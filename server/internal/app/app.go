package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/findardi/Wadi/server/internal/auth"
	"github.com/findardi/Wadi/server/internal/platform/config"
	"github.com/findardi/Wadi/server/internal/platform/otp"
	"github.com/findardi/Wadi/server/internal/platform/response"
	"github.com/findardi/Wadi/server/internal/platform/sender"
	"github.com/findardi/Wadi/server/internal/platform/token"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	router chi.Router
	addr   string
}

func New(pool *pgxpool.Pool, otpSecret, addr, jwtSecret string) *App {
	otpGen := otp.New(otpSecret)
	jwtGen := token.New(jwtSecret)

	mailCfg, _ := config.LoadMailConfig()
	mailer := sender.New(mailCfg)

	authModule := auth.NewModule(pool, otpGen, jwtGen, mailer)

	r := chi.NewRouter()
	registerGlobalMiddleware(r)

	r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, "Server Listen", nil)
	})

	authModule.RegisterRoutes(r)

	return &App{
		router: r,
		addr:   addr,
	}
}

func (a *App) Run() error {
	srv := &http.Server{
		Addr:    a.addr,
		Handler: a.router,
	}

	go func() {
		log.Printf("server running on %s", a.addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("listen: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
