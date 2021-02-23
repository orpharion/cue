package other

import "cuelang.org/go/pkg/lsp/rename/crosspkg"

func Other() {
	crosspkg.Bar
	crosspkg.Foo() //@rename("Foo", "Flamingo")
}
