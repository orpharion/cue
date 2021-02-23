// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gopls (pronounced “go please”) is an LSP server for Go.
// The Language Server Protocol allows any text editor
// to be extended with IDE-like features;
// see https://langserver.org/ for details.
//
// See https://github.com/golang/tools/blob/master/gopls/README.md
// for the most up-to-date documentation.
package main // import "golang.org/x/tools/gopls"

import (
	"context"
	"os"

	"cuelang.org/go/cmd/cuels/internal/hooks"
	"cuelang.org/go/pkg/lsp/cmd"
	tool "cuelang.org/go/pkg/lsp_tool"
)

func main() {
	ctx := context.Background()
	tool.Main(ctx, cmd.New("cuels", "", nil, hooks.Options), os.Args[1:])
}
