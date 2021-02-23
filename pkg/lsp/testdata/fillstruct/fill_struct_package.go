package fillstruct

import (
	h2 "net/http"

	"cuelang.org/go/pkg/lsp/fillstruct/data"
)

func unexported() {
	a := data.B{}   //@suggestedfix("}", "refactor.rewrite")
	_ = h2.Client{} //@suggestedfix("}", "refactor.rewrite")
}
