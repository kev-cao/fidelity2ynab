/*
This package adds support for handler groups to be added to slog. This allows
multiple handlers to be attached to each slog logger, each of which can have its own
log level. This is useful for sending logs to multiple destinations, such as stdout
and a file, with different log levels for each.
*/
package log

import (
	"context"
	"os"

	"kevincao.dev/fidelity2ynab/pkg/secrets"
	"kevincao.dev/fidelity2ynab/pkg/twilio"

	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

var logger *slog.Logger

func init() {
	group := HandlerGroup{}
	creds, err := secrets.GetSecrets()
	if err == nil && os.Getenv("ENV") != "dev" {
		// Automatically send error logs to the Twilio number
		group.AddHandler(slog.NewTextHandler(twilio.NewTwilioWriter(
			creds.TwilioAccountSid,
			creds.TwilioApiSecret,
			creds.TwilioNumber,
			creds.TwilioToNumber,
		), &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
	}
	stdoutOptions := slog.HandlerOptions{}
	if os.Getenv("ENV") == "dev" {
		stdoutOptions.Level = slog.LevelDebug
	}
	group.AddHandler(slog.NewTextHandler(os.Stdout, &stdoutOptions))
	logger = slog.New(group)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger.Log(ctx, level, msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, args...)
}

type HandlerGroup struct {
	handlers []slog.Handler
}

var _ slog.Handler = HandlerGroup{}

// Returns true if any of the handlers in the group are enabled at the given level
func (h HandlerGroup) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(context.Background(), slog.LevelInfo) {
			return true
		}
	}
	return false
}

// Passes the record to all handlers in the group. All handlers will be called,
// even if one of them returns an error. The first error encountered will be
// returned.
func (h HandlerGroup) Handle(ctx context.Context, record slog.Record) error {
	group := errgroup.Group{}
	group.SetLimit(len(h.handlers))

	for _, handler := range h.handlers {
		handler := handler
		if handler.Enabled(ctx, record.Level) {
			group.Go(func() error {
				return handler.Handle(ctx, record)
			})
		}
	}

	err := group.Wait()
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}

// Creates a new HandlerGroup with the given attributes added to each handler.
func (h HandlerGroup) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for _, handler := range h.handlers {
		newHandlers = append(newHandlers, handler.WithAttrs(attrs))
	}
	h.handlers = newHandlers
	return h
}

// Creates a new HandlerGroup with the given group added to each handler.
func (h HandlerGroup) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for _, handler := range h.handlers {
		newHandlers = append(newHandlers, handler.WithGroup(name))
	}
	h.handlers = newHandlers
	return h
}

// Adds a handler to the group
func (h *HandlerGroup) AddHandler(handler slog.Handler) {
	h.handlers = append(h.handlers, handler)
}
