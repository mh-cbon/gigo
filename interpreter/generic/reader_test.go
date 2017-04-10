package generic

import (
	"bytes"
	"testing"

	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
)

func TestReadNPrettyPrint(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)

	var buf bytes.Buffer
	x := NewReadNPrettyPrint(d, TokenerName(glanglexer.TokenName), &buf)
	interpret := NewInterpreter(x)
	interpret.Next()
	interpret.Next()
	interpret.Next()
	x.Flush()
	interpret.Next()

	expected := `Line:Pos | TokenName:tok.Type | "tok.Value"
1:0      | funcToken:17       | "func"
1:4      | wsToken:0          | " "
1:5      | wordToken:3        | "tomate"
`

	if buf.String() != expected {
		t.Errorf("Unexpected Prettyprint output\nwant=%q\ngot =%q", expected, buf.String())
	}
}
