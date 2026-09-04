package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/fx"

	"fleet/simulator/internal/infrastructure/api/controllers"
)

// Controllers holds all HTTP controllers injected via FX.
type Controllers struct {
	fx.In

	Health    *controllers.Health
	Simulator *controllers.Simulator
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

		route.Post("/simulator/start", r.controllers.Simulator.Start)
		route.Post("/simulator/stop", r.controllers.Simulator.Stop)
		route.Get("/simulator/status", r.controllers.Simulator.GetStatus)
		route.Get("/simulator/scenarios", r.controllers.Simulator.GetScenarios)
		route.Post("/simulator/scenarios/{scenario}/start", r.controllers.Simulator.StartScenario)
	})

	return r.server
}
