// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mod

import (
	"context"

	"cuelang.org/go/pkg/golang_x_tools_internal/event"
	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/protocol"
	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/source"
)

func Format(ctx context.Context, snapshot source.Snapshot, fh source.FileHandle) ([]protocol.TextEdit, error) {
	ctx, done := event.Start(ctx, "mod.Format")
	defer done()

	pm, err := snapshot.ParseMod(ctx, fh)
	if err != nil {
		return nil, err
	}
	formatted, err := pm.File.Format()
	if err != nil {
		return nil, err
	}
	// Calculate the edits to be made due to the change.
	diff, err := snapshot.View().Options().ComputeEdits(fh.URI(), string(pm.Mapper.Content), string(formatted))
	if err != nil {
		return nil, err
	}
	return source.ToProtocolEdits(pm.Mapper, diff)
}
