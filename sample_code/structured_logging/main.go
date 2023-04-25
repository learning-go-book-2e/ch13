package main

import (
	"golang.org/x/exp/slog"
	"os"
	"time"
)

func main() {
	slog.Debug("debug log message")
	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.Error("error log message")

	userID := "fred"
	loginCount := 20
	slog.Info("user login", "id", userID, "login_count", loginCount)

	options := slog.HandlerOptions{Level: slog.LevelDebug}
	handler := options.NewJSONHandler(os.Stderr)
	mySlog := slog.New(handler)
	lastLogin := time.Date(2023, 01, 01, 11, 50, 00, 00, time.UTC)
	mySlog.Debug("debug message", "id", userID, "last_login", lastLogin)

	mySlog.LogAttrs(nil, slog.LevelInfo, "faster logging", slog.String("id", userID), slog.Time("last_login", lastLogin))

	myLog := slog.NewLogLogger(mySlog.Handler(), slog.LevelDebug)
	myLog.Println("using the mySlog Handler")
}
