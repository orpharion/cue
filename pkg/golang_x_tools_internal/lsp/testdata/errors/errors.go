package errors

import (
	"cuelang.org/go/pkg/golang_x_tools_internal/lsp/types"
)

func _() {
	bob.Bob() //@complete(".")
	types.b //@complete(" //", Bob_interface)
}
