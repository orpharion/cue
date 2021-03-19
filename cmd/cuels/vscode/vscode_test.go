package vscode

import "testing"

func TestSemanticTokens(t *testing.T) {
	sts, _ := SemanticTokensFromSource("test", "a: {b: 1}")
	b, _ := sts.Json()
	println(string(b))
}
