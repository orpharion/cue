package other

import (
	references "cuelang.org/go/pkg/golang_x_tools_internal/lsp/references"
)

func GetXes() []references.X {
	return []references.X{
		{
			Y: 1, //@mark(GetXesY, "Y"),refs("Y", typeXY, GetXesY, anotherXY)
		},
	}
}

func _() {
	references.Q = "hello" //@mark(assignExpQ, "Q")
	bob := func(_ string) {}
	bob(references.Q) //@mark(bobExpQ, "Q")
}
