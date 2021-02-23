package c

import "cuelang.org/go/pkg/golang_x_tools_internal/lsp/rename/b"

func _() {
	b.Hello() //@rename("Hello", "Goodbye")
}
