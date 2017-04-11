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

	expected := "Line:Pos | TokenName:tok.Type | \"tok.Value\"\n1:0      | funcToken:60       | \"func\"\n1:4      | WsToken:0          | \" \"\n1:5      | WordToken:3        | \"tomate\"\n"

	if buf.String() != expected {
		t.Errorf("Unexpected Prettyprint output\nwant=\n%q--\ngot=\n%q--", expected, buf.String())
	}
}
