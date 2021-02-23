package errors

import (
	"cuelang.org/go/pkg/lsp/types"
)

func _() {
	bob.Bob() //@complete(".")
	types.b //@complete(" //", Bob_interface)
}
