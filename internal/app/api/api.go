package apiapp

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/handlers/user"
	changeuseremail "github.com/k6mil6/birthday-notificator/internal/api/http/handlers/user/change/email"
	changenotificationdate "github.com/k6mil6/birthday-notificator/internal/api/http/handlers/user/change/notification_date"
	userlogin "github.com/k6mil6/birthday-notificator/internal/api/http/handlers/user/login"
	userregister "github.com/k6mil6/birthday-notificator/internal/api/http/handlers/user/register"
	usersubscribe "github.com/k6mil6/birthday-notificator/internal/api/http/handlers/user/subscribe"
	usersubscriptions "github.com/k6mil6/birthday-notificator/internal/api/http/handlers/user/subscriptions"
	userunsubscribe "github.com/k6mil6/birthday-notificator/internal/api/http/handlers/user/unsubscribe"
	"github.com/k6mil6/birthday-notificator/internal/api/http/middleware/identity"
	mwlogger "github.com/k6mil6/birthday-notificator/internal/api/http/middleware/logger"
	"log/slog"
	"net/http"
)

type App struct {
	log    *slog.Logger
	router *chi.Mux
	server *http.Server
}

func New(
	ctx context.Context,
	log *slog.Logger,
	port int,
	auth httpserver.Auth,
	userService httpserver.UserService,
	secret string,
) *App {
	router := chi.NewRouter()

	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/api/v1/user/register", userregister.New(ctx, log, auth))
	router.Post("/api/v1/user/login", userlogin.New(ctx, log, auth))

	// routes with authentication
	routerWithAuth := chi.NewRouter()
	routerWithAuth.Use(identity.New(log, secret))

	routerWithAuth.Get("/api/v1/user", user.New(ctx, log, userService))
	routerWithAuth.Get("/api/v1/user/subscriptions", usersubscriptions.New(ctx, log, userService))
	routerWithAuth.Post("/api/v1/user/subscribe", usersubscribe.New(ctx, log, userService))
	routerWithAuth.Post("/api/v1/user/unsubscribe", userunsubscribe.New(ctx, log, userService))
	routerWithAuth.Post("/api/v1/user/change/email", changeuseremail.New(ctx, log, userService))
	routerWithAuth.Post("/api/v1/user/change/notification_date", changenotificationdate.New(ctx, log, userService))

	router.Mount("/", routerWithAuth)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &App{
		log:    log,
		router: router,
		server: server,
	}
}

func (a *App) Run() error {
	a.log.Info("starting server", slog.String("address", a.server.Addr))
	return a.server.ListenAndServe()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
