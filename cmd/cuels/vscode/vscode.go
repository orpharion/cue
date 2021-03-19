package vscode

import (
	"bytes"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/scanner"
	"cuelang.org/go/cue/token"
	"cuelang.org/go/internal/source"
	"cuelang.org/go/pkg/golang_x_tools_internal/span"
	"encoding/json"
	"flag"
	"fmt"
	"sort"
	"context"

	"cuelang.org/go/pkg/lsp/protocol"
	"cuelang.org/go/pkg/golang_x_tools_internal/tool"
)
// json.Marshal escapes certain html characters
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

type SemanticTokensLegendStruct struct {
	TokenTypes     []string `json:"tokenTypes"`
	TokenModifiers []string `json:"tokenModifiers"`
}

type SemanticToken struct {
	DiffLine       int
	DiffChar       int
	Length         int
	TokenType      int
	TokenModifiers int
}

func (st *SemanticToken) Array() []int {
	return []int{st.DiffLine, st.DiffChar, st.Length, st.TokenType, st.TokenModifiers}
}

func SemanticTokenFromScan(lastPos token.Pos, pos token.Pos, tok token.Token, lit string) SemanticToken {
	return SemanticToken{
		DiffLine:       pos.Line() - lastPos.Line(),
		DiffChar:       pos.Column() - lastPos.Column(),
		Length:         len(lit), // todo (ado) fixme find the correct length.
		TokenType:      int(tok),
		TokenModifiers: 0,
	}
}

type SemanticTokens []SemanticToken

func (sts SemanticTokens) Array() []int {
	var a []int
	for _, st := range sts {
		a = append(a, st.Array()...)
	}
	return a
}

func (sts SemanticTokens) Json() ([]byte, error) {
	return jsonMarshal(sts.Array())
}

var SemanticTokensLegend SemanticTokensLegendStruct
var SemanticTokensLegendJson string

func init() {
	for tok := range token.Tokens {
		SemanticTokensLegend.TokenTypes = append(SemanticTokensLegend.TokenTypes, token.Tokens[tok])
	}
	b, _ := jsonMarshal(SemanticTokensLegend)
	SemanticTokensLegendJson = string(b)
	print(SemanticTokensLegendJson)
}

func SemanticTokensFromSource(filename string, src interface{}) (SemanticTokens, error) {
	text, terr := source.Read(filename, src)
	if terr != nil {
		return SemanticTokens{}, terr
	}
	s := scanner.Scanner{}
	var errs errors.Error
	var mode scanner.Mode
	eh := func(pos token.Pos, msg string, args []interface{}) {
		errs = errors.Append(errs, errors.Newf(pos, msg, args...))
	}
	f := token.NewFile(filename, 0, len(text))
	lastPos := f.Pos(0, token.RelPos(0))
	var (
		pos token.Pos
		tok token.Token
		lit string
	)
	s.Init(f, text, eh, mode) // todo (ado) always zero offset?
	var sts SemanticTokens
	for true {
		pos, tok, lit = s.Scan()
		st := SemanticTokenFromScan(lastPos, pos, tok, lit)
		lastPos = pos
		sts = append(sts, st)
		if tok == token.EOF {
			break
		}
	}
	return sts, errs
}


// symbols implements the symbols verb for gopls
type symbols struct {
	app *Application
}

func (r *symbols) Name() string      { return "symbols" }
func (r *symbols) Usage() string     { return "<file>" }
func (r *symbols) ShortHelp() string { return "display selected file's symbols" }
func (r *symbols) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprint(f.Output(), `
Example:
  $ gopls symbols helper/helper.go
`)
	f.PrintDefaults()
}
func (r *symbols) Run(ctx context.Context, args ...string) error {
	if len(args) != 1 {
		return tool.CommandLineErrorf("symbols expects 1 argument (position)")
	}

	conn, err := r.app.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.terminate(ctx)

	from := span.Parse(args[0])

	p := protocol.SemanticTokens{
		ResultID: "",
		Data:     nil,
	}

	p := protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.URIFromSpanURI(from.URI()),
		},
	}
	symbols, err := conn.SemanticTokensFull(ctx, &p)
	if err != nil {
		return err
	}
	for _, s := range symbols {
		if m, ok := s.(map[string]interface{}); ok {
			s, err = mapToSymbol(m)
			if err != nil {
				return err
			}
		}
		switch t := s.(type) {
		case protocol.DocumentSymbol:
			printDocumentSymbol(t)
		case protocol.SymbolInformation:
			printSymbolInformation(t)
		}
	}
	return nil
}

func mapToSymbol(m map[string]interface{}) (interface{}, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	if _, ok := m["selectionRange"]; ok {
		var s protocol.DocumentSymbol
		if err := json.Unmarshal(b, &s); err != nil {
			return nil, err
		}
		return s, nil
	}

	var s protocol.SymbolInformation
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return s, nil
}

func printDocumentSymbol(s protocol.DocumentSymbol) {
	fmt.Printf("%s %s %s\n", s.Name, s.Kind, positionToString(s.SelectionRange))
	// Sort children for consistency
	sort.Slice(s.Children, func(i, j int) bool {
		return s.Children[i].Name < s.Children[j].Name
	})
	for _, c := range s.Children {
		fmt.Printf("\t%s %s %s\n", c.Name, c.Kind, positionToString(c.SelectionRange))
	}
}

func printSymbolInformation(s protocol.SymbolInformation) {
	fmt.Printf("%s %s %s\n", s.Name, s.Kind, positionToString(s.Location.Range))
}

func positionToString(r protocol.Range) string {
	return fmt.Sprintf("%v:%v-%v:%v",
		r.Start.Line+1,
		r.Start.Character+1,
		r.End.Line+1,
		r.End.Character+1,
	)
}

