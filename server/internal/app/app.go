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

	"github.com/findardi/Wadi/server/internal/access"
	accessrepo "github.com/findardi/Wadi/server/internal/access/repository"
	accessservice "github.com/findardi/Wadi/server/internal/access/service"
	"github.com/findardi/Wadi/server/internal/auth"
	authrepo "github.com/findardi/Wadi/server/internal/auth/repository"
	authservice "github.com/findardi/Wadi/server/internal/auth/service"
	"github.com/findardi/Wadi/server/internal/invitation"
	"github.com/findardi/Wadi/server/internal/platform/config"
	"github.com/findardi/Wadi/server/internal/platform/oauth"
	"github.com/findardi/Wadi/server/internal/platform/otp"
	"github.com/findardi/Wadi/server/internal/platform/ratelimit"
	"github.com/findardi/Wadi/server/internal/platform/response"
	"github.com/findardi/Wadi/server/internal/platform/sender"
	"github.com/findardi/Wadi/server/internal/platform/token"
	"github.com/findardi/Wadi/server/internal/workspace"
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
	limiter := ratelimit.NewMemory()

	ghCfg := config.LoadOAuth("OAUTH_GITHUB")
	ggCfg := config.LoadOAuth("OAUTH_GOOGLE")
	providers := map[string]oauth.Provider{
		"github": oauth.NewGithub(ghCfg.ClientID, ghCfg.ClientSecret, ghCfg.RedirectURL),
		"google": oauth.NewGoogle(ggCfg.ClientID, ggCfg.ClientSecret, ggCfg.RedirectURL),
	}

	webURL := config.GetEnv("WEB_URL", "http://localhost:5173")

	authsvc := authservice.NewAuthService(authrepo.New(pool), otpGen, jwtGen, mailer, nil)
	accessSvc := accessservice.NewAccessService(accessrepo.New(pool), mailer, authsvc, otpGen, webURL)

	authModule := auth.NewModule(pool, otpGen, jwtGen, mailer, limiter, providers, accessSvc)
	workspaceModule := workspace.NewModule(pool, jwtGen, accessSvc)
	accessModule := access.NewModule(pool, jwtGen, mailer, authsvc, otpGen, webURL)
	invitationModule := invitation.NewModule(pool, jwtGen)

	r := chi.NewRouter()
	registerGlobalMiddleware(r)

	r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, "Server Listen", nil)
	})

	authModule.RegisterRoutes(r)
	workspaceModule.RegisterRoutes(r)
	accessModule.RegisterRoutes(r)
	invitationModule.RegisterRoutes(r)

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
