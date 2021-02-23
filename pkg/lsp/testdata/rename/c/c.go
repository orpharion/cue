package c

import "cuelang.org/go/pkg/lsp/rename/b"

func _() {
	b.Hello() //@rename("Hello", "Goodbye")
}
