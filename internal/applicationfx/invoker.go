package applicationfx

import (
	"log/slog"
	"runtime"

	"github.com/bosonicalio/geck/application"

	"github.com/bosonicalio/enclave/internal/globallog"
)

func logAppStart(config application.Application) {
	globallog.Logger().Info("starting application",
		slog.String("name", config.Name),
		slog.String("environment", config.Environment.String()),
		slog.String("version", config.Version.String()),
		slog.String("instance_id", config.InstanceID),
		slog.Group("runtime",
			slog.Int("cpus", runtime.NumCPU()),
			slog.Group("go",
				slog.String("version", runtime.Version()),
				slog.String("os", runtime.GOOS),
				slog.String("arch", runtime.GOARCH),
			),
		),
	)
}
