package httpfx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"golang.org/x/crypto/acme/autocert"

	"github.com/tesserical/geck/application"
	geckhttp "github.com/tesserical/geck/transport/http"

	"github.com/tesserical/enclave/internal/globallog"
)

type startServerDeps struct {
	fx.In

	Lifecycle fx.Lifecycle
	Echo      *echo.Echo
	Config    serverConfig
}

func startServer(deps startServerDeps) {
	deps.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				globallog.Logger().InfoContext(ctx, "starting http server",
					slog.String("addr", deps.Config.Address),
				)
				var err error
				if deps.Config.EnableTLS {
					// TODO: use TLS configuration from the config
					err = deps.Echo.StartTLS(deps.Config.Address, "", "")
				} else if deps.Config.EnableAutoTLS {
					deps.Echo.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
					err = deps.Echo.StartAutoTLS(deps.Config.Address)
				} else {
					err = deps.Echo.Start(deps.Config.Address)
				}
				if errors.Is(err, http.ErrServerClosed) {
					return
				} else if err != nil {
					globallog.Logger().ErrorContext(ctx, "failed during http server execution",
						slog.String("error", err.Error()),
					)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			globallog.Logger().InfoContext(ctx, "stopping http server")
			return deps.Echo.Shutdown(ctx)
		},
	})
}

type registerServerEndpointsDeps struct {
	fx.In
	Echo        *echo.Echo
	Logger      *slog.Logger
	App         application.Application
	Controllers []geckhttp.Controller `group:"http_controllers"`
}

func registerServerEndpoints(deps registerServerEndpointsDeps) {
	pathPrefix := geckhttp.RegisterServerEndpoints(deps.Echo, deps.App, deps.Controllers)
	globallog.Logger().Debug("registered http server endpoints",
		slog.String("path_prefix", pathPrefix),
		slog.Int("total_controllers", len(deps.Controllers)),
	)
}
