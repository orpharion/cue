// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package myers_test

import (
	"testing"

	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/diff/difftest"
	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/diff/myers"
)

func TestDiff(t *testing.T) {
	difftest.DiffTest(t, myers.ComputeEdits)
}
