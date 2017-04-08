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

// func TestReadProtected(t *testing.T) {
//
// 	str := `func tomate() {
// 		var expr string = "some"
// 		expr2 := "other"
// }`
// 	d := stringTokenizer(str)
//
// 	y := NewReadSlowly(d)
// 	x := NewReadProtected(y)
// 	interpret := NewInterpreter(x)
// 	interpret.Next()
// 	interpret.Next()
// 	interpret.Next()
// 	interpret.Next()
// 	y.Slow = 3
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("The code did not panic")
// 		}
// 	}()
// 	fmt.Println(interpret.Next())
// }
//
// type ReadSlowly struct {
// 	Reader TokenerReaderOK
// 	Slow   int
// }
//
// func (r *ReadSlowly) NextToken() Tokener {
// 	<-time.After(time.Second * time.Duration(r.Slow))
// 	return r.Reader.NextToken()
// }
//
// func NewReadSlowly(r TokenerReaderOK) *ReadSlowly {
// 	ret := &ReadSlowly{
// 		Reader: r,
// 	}
// 	return ret
// }
