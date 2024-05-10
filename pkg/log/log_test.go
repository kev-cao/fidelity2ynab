package log

import (
	"testing"

	"golang.org/x/exp/slog"
)

type testSliceWriter struct {
	slice *[]string
}

func (t *testSliceWriter) Write(p []byte) (n int, err error) {
	*t.slice = append(*t.slice, string(p))
	return len(p), nil
}

// Tests that the log group correctly filters logs based on the handler's level
func TestLogGroup(t *testing.T) {
	errsOnly := make([]string, 0, 4)
	warnsAndAbove := make([]string, 0, 4)
	infoAndAbove := make([]string, 0, 4)
	allLogs := make([]string, 0, 4)

	group := HandlerGroup{}
	group.AddHandler(slog.NewTextHandler(
		&testSliceWriter{&errsOnly},
		&slog.HandlerOptions{Level: slog.LevelError},
	))
	group.AddHandler(slog.NewTextHandler(
		&testSliceWriter{&warnsAndAbove},
		&slog.HandlerOptions{Level: slog.LevelWarn},
	))
	group.AddHandler(slog.NewTextHandler(
		&testSliceWriter{&infoAndAbove}, &slog.HandlerOptions{},
	))
	group.AddHandler(slog.NewTextHandler(
		&testSliceWriter{&allLogs},
		&slog.HandlerOptions{Level: slog.LevelDebug},
	))

	logger := slog.New(group)
	logger.Debug("debug sentinel")
	logger.Info("info sentinel")
	logger.Warn("warn sentinel")
	logger.Error("error sentinel")

	if len(errsOnly) != 1 {
		t.Errorf("Expected 1 error log, got %s", errsOnly)
	}
	if len(warnsAndAbove) != 2 {
		t.Errorf("Expected 2 warn logs, got %s", warnsAndAbove)
	}
	if len(infoAndAbove) != 3 {
		t.Errorf("Expected 3 warn logs, got %s", warnsAndAbove)
	}
	if len(allLogs) != 4 {
		t.Errorf("Expected 4 total logs, got %s", allLogs)
	}
}
