// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdtest

import (
	"fmt"
	"testing"

	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/tests"
	"cuelang.org/go/pkg/golang_x_tools_internal/span"
)

func (r *runner) SuggestedFix(t *testing.T, spn span.Span, actionKinds []string, expectedActions int) {
	uri := spn.URI()
	filename := uri.Filename()
	args := []string{"fix", "-a", fmt.Sprintf("%s", spn)}
	for _, kind := range actionKinds {
		if kind == "refactor.rewrite" {
			t.Skip("refactor.rewrite is not yet supported on the command line")
		}
	}
	args = append(args, actionKinds...)
	got, stderr := r.NormalizeGoplsCmd(t, args...)
	if stderr == "ExecuteCommand is not yet supported on the command line" {
		t.Skipf(stderr)
	}
	want := string(r.data.Golden("suggestedfix_"+tests.SpanName(spn), filename, func() ([]byte, error) {
		return []byte(got), nil
	}))
	if want != got {
		t.Errorf("suggested fixes failed for %s:\n%s", filename, tests.Diff(t, want, got))
	}
}