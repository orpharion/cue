package danglingstmt

import "cuelang.org/go/pkg/golang_x_tools_internal/lsp/foo"

func _() {
	foo. //@rank(" //", Foo)
	var _ = []string{foo.} //@rank("}", Foo)
}
