// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"bytes"
	"context"

	"cuelang.org/go/pkg/golang_x_tools_internal/event"
	"cuelang.org/go/pkg/golang_x_tools_internal/event/core"
	"cuelang.org/go/pkg/golang_x_tools_internal/event/export"
	"cuelang.org/go/pkg/golang_x_tools_internal/event/label"
	"cuelang.org/go/pkg/golang_x_tools_internal/xcontext"
)

type contextKey int

const (
	clientKey = contextKey(iota)
)

func WithClient(ctx context.Context, client Client) context.Context {
	return context.WithValue(ctx, clientKey, client)
}

func LogEvent(ctx context.Context, ev core.Event, lm label.Map, mt MessageType) context.Context {
	client, ok := ctx.Value(clientKey).(Client)
	if !ok {
		return ctx
	}
	buf := &bytes.Buffer{}
	p := export.Printer{}
	p.WriteEvent(buf, ev, lm)
	msg := &LogMessageParams{Type: mt, Message: buf.String()}
	// Handle messages generated via event.Error, which won't have a level Label.
	if event.IsError(ev) {
		msg.Type = Error
	}
	go client.LogMessage(xcontext.Detach(ctx), msg)
	return ctx
}
