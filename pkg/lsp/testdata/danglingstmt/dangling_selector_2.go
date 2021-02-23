package danglingstmt

import "cuelang.org/go/pkg/lsp/foo"

func _() {
	foo. //@rank(" //", Foo)
	var _ = []string{foo.} //@rank("}", Foo)
}
