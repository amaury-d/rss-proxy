package rss

import (
	"log/slog"
	"os"
)

var Logger = slog.New(
	slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	),
)
