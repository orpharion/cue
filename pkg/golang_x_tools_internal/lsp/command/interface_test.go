// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package command_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/command/gen"
	"cuelang.org/go/pkg/golang_x_tools_internal/testenv"
)

func TestGenerated(t *testing.T) {
	testenv.NeedsGoBuild(t) // This is a lie. We actually need the source code.

	onDisk, err := ioutil.ReadFile("command_gen.go")
	if err != nil {
		t.Fatal(err)
	}

	generated, err := gen.Generate()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(onDisk, generated) {
		t.Error("command_gen.go is stale -- regenerate")
	}
}
