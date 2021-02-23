// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"

	"cuelang.org/go/pkg/golang_x_tools_internal/event"
	"cuelang.org/go/pkg/lsp/debug/tag"
	"cuelang.org/go/pkg/lsp/protocol"
	"cuelang.org/go/pkg/lsp/source"
)

func (s *Server) documentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) ([]interface{}, error) {
	ctx, done := event.Start(ctx, "lsp.Server.documentSymbol")
	defer done()

	snapshot, fh, ok, release, err := s.beginFileRequest(ctx, params.TextDocument.URI, source.Go)
	defer release()
	if !ok {
		return []interface{}{}, err
	}
	docSymbols, err := source.DocumentSymbols(ctx, snapshot, fh)
	if err != nil {
		event.Error(ctx, "DocumentSymbols failed", err, tag.URI.Of(fh.URI()))
		return []interface{}{}, nil
	}
	// Convert the symbols to an interface array.
	// TODO: Remove this once the lsp deprecates SymbolInformation.
	symbols := make([]interface{}, len(docSymbols))
	for i, s := range docSymbols {
		if snapshot.View().Options().HierarchicalDocumentSymbolSupport {
			symbols[i] = s
			continue
		}
		// If the client does not support hierarchical document symbols, then
		// we need to be backwards compatible for now and return SymbolInformation.
		symbols[i] = protocol.SymbolInformation{
			Name:       s.Name,
			Kind:       s.Kind,
			Deprecated: s.Deprecated,
			Location: protocol.Location{
				URI:   params.TextDocument.URI,
				Range: s.Range,
			},
		}
	}
	return symbols, nil
}
