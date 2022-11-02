package logger

import (
	"github.com/getsentry/sentry-go"
	// sentrygin "github.com/getsentry/sentry-go/gin"
)

type LoggerSentry struct {
	hub *sentry.Hub
}

func InitClient(dsn string, traceRate float64) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,

		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		return err
	}

	return nil
}

// func New() LoggerSentry {

// }

// func SetGinOptions() {
// 	handler := sentrygin.New(sentrygin.Options {
// 		Repanic: true,
// 	})

// }
