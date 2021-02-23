// Package another has another type.
package another

import (
	other "cuelang.org/go/pkg/golang_x_tools_internal/lsp/references/other"
)

func _() {
	xes := other.GetXes()
	for _, x := range xes {
		_ = x.Y //@mark(anotherXY, "Y"),refs("Y", typeXY, anotherXY, GetXesY)
	}
}
