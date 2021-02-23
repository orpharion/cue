package other

import "cuelang.org/go/pkg/golang_x_tools_internal/lsp/rename/crosspkg"

func Other() {
	crosspkg.Bar
	crosspkg.Foo() //@rename("Foo", "Flamingo")
}
