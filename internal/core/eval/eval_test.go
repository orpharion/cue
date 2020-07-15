// Copyright 2020 CUE Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eval_test

import (
	"flag"
	"fmt"
	"testing"

	"cuelang.org/go/internal/core/debug"
	"cuelang.org/go/internal/core/eval"
	"cuelang.org/go/internal/core/validate"
	"cuelang.org/go/internal/cuetxtar"
	"cuelang.org/go/internal/legacy/cue"
	"cuelang.org/go/pkg/strings"
	"github.com/rogpeppe/go-internal/txtar"
)

var (
	update = flag.Bool("update", false, "update the test files")
	todo   = flag.Bool("todo", false, "run tests marked with #todo-compile")
)

func TestEval(t *testing.T) {
	test := cuetxtar.TxTarTest{
		Root:   "../../../cue/testdata",
		Name:   "eval",
		Update: *update,
		Skip:   alwaysSkip,
		ToDo:   needFix,
	}

	if *todo {
		test.ToDo = nil
	}

	r := cue.NewRuntime()

	test.Run(t, func(t *cuetxtar.Test) {
		a := t.ValidInstances()

		v, err := r.Build(a[0])
		if err != nil {
			t.Fatal(err)
		}

		e := eval.New(r)
		ctx := e.NewContext(v)
		v.Finalize(ctx)

		t.Log(e.Stats())

		if b := validate.Validate(r, v, &validate.Config{
			AllErrors: true,
		}); b != nil {
			fmt.Fprintln(t, "Errors:")
			t.WriteErrors(b.Err)
			fmt.Fprintln(t, "")
			fmt.Fprintln(t, "Result:")
		}

		if v == nil {
			return
		}

		debug.WriteNode(t, r, v, &debug.Config{Cwd: t.Dir})
		fmt.Fprintln(t)
	})
}

var alwaysSkip = map[string]string{
	"compile/erralias": "compile error",
}

var needFix = map[string]string{
	"export/027": "cycle",
	"export/028": "cycle",
	"export/030": "cycle",

	"cycle/025_cannot_resolve_references_that_would_be_ambiguous": "cycle",
	"compile/scope": "cycle",
}

// TestX is for debugging. Do not delete.
func TestX(t *testing.T) {
	in := `
-- cue.mod/module.cue --
module: "example.com"

-- in.cue --
	`

	if strings.HasSuffix(strings.TrimSpace(in), ".cue --") {
		t.Skip()
	}

	a := txtar.Parse([]byte(in))
	instance := cuetxtar.Load(a, "/tmp/test")[0]
	if instance.Err != nil {
		t.Fatal(instance.Err)
	}

	r := cue.NewRuntime()

	v, err := r.Build(instance)
	if err != nil {
		t.Fatal(err)
	}
	t.Error(debug.NodeString(r, v, nil))

	e := eval.New(r)
	ctx := e.NewContext(v)
	v.Finalize(ctx)

	t.Error(debug.NodeString(r, v, nil))

	t.Log(e.Stats())
}