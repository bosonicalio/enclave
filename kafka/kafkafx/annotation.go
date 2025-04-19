package kafkafx

import (
	"github.com/hadroncorp/geck/transport/stream/kafka"
	"go.uber.org/fx"
)

// AsController is a helper function to annotate a function that returns a [stream.Controller].
func AsController(t any) any {
	return fx.Annotate(
		t,
		fx.As(new(kafka.Controller)),
		fx.ResultTags(`group:"kafka_controllers"`),
	)
}
