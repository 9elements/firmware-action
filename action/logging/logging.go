// SPDX-License-Identifier: MIT

// Package logging for dealing with logging, log settings and log format
package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

// InitLogger is used to initialize logger
func InitLogger(level slog.Level, opts ...PrettyHandlerOption) {
	// Deal with log level
	//   LevelDebug Level = -4
	//   LevelInfo  Level = 0
	//   LevelWarn  Level = 4
	//   LevelError Level = 8
	_ = slog.SetLogLoggerLevel(level)

	options := []PrettyHandlerOption{WithLevel(level)}
	options = append(options, opts[:]...)
	pHandler := NewPrettyHandler(
		os.Stdout,
		options[:]...,
	)

	slog.SetDefault(slog.New(pHandler))
}

// TrimNameFunction just trims the name
func TrimNameFunction(pc uintptr) string {
	// 'runtime.FuncForPC(pc).Name()' is nice and all, but it will return this monstrosity:
	//   github.com/9elements/firmware-action/action/<package>.<func>...
	// So this function is just to trim it down
	// Usage:
	//   pc, _, _, _ := runtime.Caller(0)
	//   name := logging.TrimFunctionName(pc)
	return filepath.Base(runtime.FuncForPC(pc).Name())
}
