package nodisk

import (
	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}
