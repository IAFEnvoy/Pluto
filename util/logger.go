package util

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
)

var LOGGER = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
	Level:      slog.LevelDebug,
	TimeFormat: time.Kitchen,
	NoColor:    false,
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.MessageKey {
			return slog.Attr{
				Key:   a.Key,
				Value: slog.StringValue("\x1B[1m" + a.Value.String() + "\033[0m"),
			}
		}
		return a
	},
}))
