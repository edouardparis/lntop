package pubsub

import (
	"context"
	"sync"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
)

func Run(ctx context.Context, app *app.App, sub chan *events.Event) error {
	logger := app.Logger.With(logging.String("logger", "pubsub"))
	wg := &sync.WaitGroup{}

	logger.Debug("Starting...")
	wg.Wait()
	return nil
}
