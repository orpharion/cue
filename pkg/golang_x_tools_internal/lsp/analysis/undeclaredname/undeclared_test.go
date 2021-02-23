// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package undeclaredname_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/analysis/undeclaredname"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, undeclaredname.Analyzer, "a")
}
