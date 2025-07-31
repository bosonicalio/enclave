package httpfx

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	geckhttp "github.com/bosonicalio/geck/transport/http"

	"github.com/bosonicalio/enclave/internal/osenv"
)

// ServerModule is the `uber/fx` module of the [geckhttp] package, aimed for HTTP servers.
//
// This module uses `labstack/echo` as HTTP framework for internal operations.
var ServerModule = fx.Module("enclave/transport/http/server",
	fx.Provide(
		osenv.ParseAs[serverConfig],
		newServer,
	),
	fx.Invoke(
		registerServerEndpoints,
		startServer,
	),
)

// -- Factory --

func newServer(cfg serverConfig) *echo.Echo {
	return geckhttp.NewEchoServer(
		geckhttp.WithServerErrorResponseCodec(cfg.ErrResponseCodec),
	)
}
