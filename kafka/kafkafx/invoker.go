package kafkafx

import (
	"context"
	"errors"
	"log/slog"

	"go.uber.org/fx"

	"github.com/hadroncorp/geck/transport/stream/kafka"
)

func registerManagerReaders(rm kafka.ReaderManager, logger *slog.Logger, controllers []kafka.Controller) {
	for _, controller := range controllers {
		controller.RegisterReaders(rm)
	}
	logger.Info("registered kafka readers to manager", slog.Int("count", len(controllers)))
}

func startManagerReaders(lc fx.Lifecycle, logger *slog.Logger, rm kafka.ReaderManager) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Info("starting kafka reader manager")
				if err := rm.Start(); err != nil && !errors.Is(err, kafka.ErrReaderManagerClosed) {
					logger.Error("failed to start kafka reader manager", slog.String("error", err.Error()))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping kafka reader manager")
			return rm.Close(ctx)
		},
	})
}
