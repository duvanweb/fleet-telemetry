package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/fx"

	"fleet/alert-service/internal/infrastructure/api/controllers"
)

// Controllers holds all HTTP controllers injected via FX.
type Controllers struct {
	fx.In

	Health *controllers.Health
}

// Router wraps the chi mux with its controllers.
type Router struct {
	controllers Controllers
	server      *chi.Mux
}

// NewRouter creates and returns a new Router.
func NewRouter(server *chi.Mux, c Controllers) *Router {
	return &Router{controllers: c, server: server}
}

func (r *Router) start(basePath string) http.Handler {
	r.server.Use(middleware.Logger)
	r.server.Use(middleware.Recoverer)
	r.server.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.server.Route(basePath, func(route chi.Router) {
		route.Get("/health", r.controllers.Health.GetHealth)
		route.Get("/ready", r.controllers.Health.GetReady)
	})

	return r.server
}
