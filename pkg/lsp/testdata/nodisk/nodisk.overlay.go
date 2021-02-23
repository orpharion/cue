package nodisk

import (
	"cuelang.org/go/pkg/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}
