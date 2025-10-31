// SPDX-License-Identifier: MIT

// Package logging for dealing with logging, log settings and log format
//
//	Inspiration:
//		https://github.com/golang/example/tree/master/slog-handler-guide
//		https://dusted.codes/creating-a-pretty-console-logger-using-gos-slog-package
package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"slices"
	"strings"
	"sync"
)

var (
	errInnerHandle          = errors.New("calling inner handler's Handle failed")
	errInnerHandleUnmarshal = errors.New("unmarshal output on inner handler's Handle failed")
	bufferSize              = 1024
)

// PrettyHandler is for prettier logs
type PrettyHandler struct {
	baseHandler slog.Handler
	opts        *slog.HandlerOptions
	mutex       *sync.Mutex
	buffer      *bytes.Buffer
	json        bool
	indent      bool
	skipframes  int
	output      io.Writer
}

// PrettyHandlerOption is for functional option pattern
type PrettyHandlerOption func(*PrettyHandler)

// WithJSON enables or disables JSON output (handy for programmatic processing out stderr)
func WithJSON(json bool) PrettyHandlerOption {
	return func(h *PrettyHandler) {
		h.json = json
	}
}

// WithIndent enables or disables indentation (handy for programmatic processing out stderr)
func WithIndent(indent bool) PrettyHandlerOption {
	return func(h *PrettyHandler) {
		h.indent = indent
	}
}

// WithLevel changes log level
func WithLevel(level slog.Level) PrettyHandlerOption {
	return func(h *PrettyHandler) {
		h.opts.Level = level
		h.opts.AddSource = level == slog.LevelDebug
	}
}

// WithSkipFrames changes how many stack frames to ascend in runtime.Caller()
//
//	0 = logging.(*PrettyHandler).Handle
//	1 = slog.(*Logger).log
//	2 = slog.Debug || slog.Info || slog.Warn || slog.Error
//	3 = place where we call slog.<LEVEL>
func WithSkipFrames(frames int) PrettyHandlerOption {
	return func(h *PrettyHandler) {
		h.skipframes = frames
	}
}

// NewPrettyHandler returns a PrettyHandler
func NewPrettyHandler(output io.Writer, opts ...PrettyHandlerOption) *PrettyHandler {
	buffer := &bytes.Buffer{}

	options := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}

	h := &PrettyHandler{
		baseHandler: slog.NewJSONHandler(
			buffer,
			&options,
		),
		mutex:      &sync.Mutex{},
		buffer:     buffer,
		indent:     true,
		skipframes: 3,
		output:     output,
		opts:       &options,
	}
	for _, opt := range opts {
		opt(h)
	}

	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}

	return h
}

// Enabled reports whether the handler handles records at the given level
func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// WithAttrs returns a new Handler whose attributes consist of both the receiver's attributes
// and the arguments
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// TODO
	return &PrettyHandler{
		baseHandler: h.baseHandler.WithAttrs(attrs),
		mutex:       h.mutex,
		buffer:      h.buffer,
		indent:      h.indent,
		skipframes:  h.skipframes,
		output:      h.output,
		opts:        h.opts,
	}
}

// WithGroup returns a new Handler with the given group appended to the receiver's existing groups
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// TODO
	return &PrettyHandler{
		baseHandler: h.baseHandler.WithGroup(name),
		mutex:       h.mutex,
		buffer:      h.buffer,
		indent:      h.indent,
		skipframes:  h.skipframes,
		output:      h.output,
		opts:        h.opts,
	}
}

// Handle handles the Record
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add location where slog.<LEVEL> has been called
	pc, _, _, _ := runtime.Caller(h.skipframes)
	funcName := TrimNameFunction(pc)
	r.Add("origin of this message", funcName)

	// Convert Record into map
	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	// Prepare temporary buffer
	tmpBuffer := make([]byte, 0, bufferSize)

	if h.json {
		// Pretty-print JSON
		var bytes []byte
		if h.indent {
			bytes, err = json.MarshalIndent(attrs, "", "  ")
		} else {
			bytes, err = json.Marshal(attrs)
		}

		if err != nil {
			return err
		}

		tmpBuffer = append(tmpBuffer, bytes[:]...)
	} else {
		// Don't do JSON, rather use some more human-readable format
		var details strings.Builder

		ignore := []string{"level", "msg"}
		//   ^^^ ignore is to omit fields from 'details' section because already in main body of log message
		for key, val := range attrs {
			if slices.Contains(ignore, key) {
				continue
			}

			switch val.(type) {
			case bool:
				details.WriteString(fmt.Sprintf("    - %s: %t\n", key, val))
			default:
				details.WriteString(fmt.Sprintf("    - %s: %s\n", key, val))
			}
		}

		tmpBuffer = fmt.Appendf(
			tmpBuffer,
			"[%-7s] %s\n%s",
			attrs["level"],
			attrs["msg"],
			details.String(),
		)
	}

	// Write into output
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, err = h.output.Write(tmpBuffer)

	return err
}

func (h *PrettyHandler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	// Deal with concurrency
	h.mutex.Lock()

	defer func() {
		h.buffer.Reset()
		h.mutex.Unlock()
	}()

	// Use inner handler (JSONHandler) to compute JSON object
	//   h.baseHandler is initialized with h.buffer as output
	if err := h.baseHandler.Handle(ctx, r); err != nil {
		return nil, errors.Join(errInnerHandle, err)
	}

	// Unmarshal the buffer into map and return that
	var attrs map[string]any
	if err := json.Unmarshal(h.buffer.Bytes(), &attrs); err != nil {
		return nil, errors.Join(errInnerHandleUnmarshal, err)
	}

	return attrs, nil
}
