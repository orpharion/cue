// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"

	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/mod"
	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/protocol"
	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/source"
)

func (s *Server) formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	snapshot, fh, ok, release, err := s.beginFileRequest(ctx, params.TextDocument.URI, source.UnknownKind)
	defer release()
	if !ok {
		return nil, err
	}
	switch fh.Kind() {
	case source.Mod:
		return mod.Format(ctx, snapshot, fh)
	case source.Go:
		return source.Format(ctx, snapshot, fh)
	}
	return nil, nil
}
