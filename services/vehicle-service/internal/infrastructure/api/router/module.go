package router

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"fleet/shared/pkg/logger"
	"fleet/vehicle-service/internal/infrastructure/api/controllers"
)

// Module registers all router and HTTP server components with FX.
func Module() fx.Option {
	return fx.Module(
		"api",
		fx.Provide(
			chi.NewRouter,
			NewRouter,
			controllers.NewHealth,
			controllers.NewVehicle,
		),
		fx.Invoke(registerHooks),
	)
}

func registerHooks(lc fx.Lifecycle, shutdown fx.Shutdowner, router *Router, log logger.Logger) {
	server := &http.Server{Addr: ":8081", Handler: router.start("/api")}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Infow(ctx, "vehicle-service listening", "addr", server.Addr)
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Errorw(ctx, "server error", "error", err)
					_ = shutdown.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
