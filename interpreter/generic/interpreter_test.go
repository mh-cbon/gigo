package generic

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	// genericlexer "github.com/mh-cbon/gigo/lexer/generic"

	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
	lexer "github.com/mh-cbon/state-lexer"
)

func TestNextUntilEOF(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	expected := []Tokener{
		&TokenWithPos{Token: lexer.Token{Type: 17, Value: "func"}, Pos: TokenPos{Line: 1, Pos: 0}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 1, Pos: 4}},
		&TokenWithPos{Token: lexer.Token{Type: 3, Value: "tomate"}, Pos: TokenPos{Line: 1, Pos: 5}},
		&TokenWithPos{Token: lexer.Token{Type: 20, Value: "("}, Pos: TokenPos{Line: 1, Pos: 11}},
		&TokenWithPos{Token: lexer.Token{Type: 21, Value: ")"}, Pos: TokenPos{Line: 1, Pos: 12}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 1, Pos: 13}},
		&TokenWithPos{Token: lexer.Token{Type: 45, Value: "{"}, Pos: TokenPos{Line: 1, Pos: 14}},
		&TokenWithPos{Token: lexer.Token{Type: 6, Value: "\n"}, Pos: TokenPos{Line: 1, Pos: 15}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: "\t"}, Pos: TokenPos{Line: 2, Pos: 0}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: "\t"}, Pos: TokenPos{Line: 2, Pos: 1}},
		&TokenWithPos{Token: lexer.Token{Type: 16, Value: "var"}, Pos: TokenPos{Line: 2, Pos: 2}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 2, Pos: 5}},
		&TokenWithPos{Token: lexer.Token{Type: 3, Value: "expr"}, Pos: TokenPos{Line: 2, Pos: 6}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 2, Pos: 10}},
		&TokenWithPos{Token: lexer.Token{Type: 85, Value: "string"}, Pos: TokenPos{Line: 2, Pos: 11}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 2, Pos: 17}},
		&TokenWithPos{Token: lexer.Token{Type: 22, Value: "="}, Pos: TokenPos{Line: 2, Pos: 18}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 2, Pos: 19}},
		&TokenWithPos{Token: lexer.Token{Type: 4, Value: "\"some\""}, Pos: TokenPos{Line: 2, Pos: 20}},
		&TokenWithPos{Token: lexer.Token{Type: 6, Value: "\n"}, Pos: TokenPos{Line: 2, Pos: 26}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: "\t"}, Pos: TokenPos{Line: 3, Pos: 0}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: "\t"}, Pos: TokenPos{Line: 3, Pos: 1}},
		&TokenWithPos{Token: lexer.Token{Type: 3, Value: "expr2"}, Pos: TokenPos{Line: 3, Pos: 2}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 3, Pos: 7}},
		&TokenWithPos{Token: lexer.Token{Type: 23, Value: ":="}, Pos: TokenPos{Line: 3, Pos: 8}},
		&TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 3, Pos: 10}},
		&TokenWithPos{Token: lexer.Token{Type: 4, Value: "\"other\""}, Pos: TokenPos{Line: 3, Pos: 11}},
		&TokenWithPos{Token: lexer.Token{Type: 6, Value: "\n"}, Pos: TokenPos{Line: 3, Pos: 18}},
		&TokenWithPos{Token: lexer.Token{Type: 46, Value: "}"}, Pos: TokenPos{Line: 4, Pos: 0}},
	}

	c := 0
	for {
		tok := interpret.Next()
		if tok == nil {
			break
		}
		if c < len(expected) {
			if !compareToken(t, c, expected[c], tok) {
				return
			}
		}
		c++
	}

	want := len(expected)
	got := c
	if want != got {
		t.Errorf("Wrong number of tokens want=%v, got=%v", want, got)
	}
}

func TestNextAfterEOF(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	for {
		tok := interpret.Next()
		if tok == nil {
			break
		}
	}

	if got := interpret.Next(); got != nil {
		t.Errorf("Next must return nil got=%v", got)
	}

	if interpret.Ended() == false {
		t.Errorf("Ended must return true got=%v", false)
	}

}

func TestNext(t *testing.T) {
	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	got := interpret.Next()
	want := &TokenWithPos{Token: lexer.Token{Type: 17, Value: "func"}, Pos: TokenPos{Line: 1, Pos: 0}}

	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
}

func TestRewind(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	interpret.Rewind()

	want := interpret.Next()

	interpret.Rewind()

	got := interpret.Next()

	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}

	interpret.Rewind()
	interpret.Rewind()
	interpret.Rewind()
	got = interpret.Next()

	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
}

func TestRewindAll(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	want := interpret.Next()
	interpret.Next()
	interpret.Next()
	interpret.Next()

	interpret.RewindAll()

	got := interpret.Next()

	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
}

func TestLast(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	if got := interpret.Last(); got != nil {
		t.Errorf("Last must return nil got=%v", got)
	}

	want := interpret.Next()

	got := interpret.Last()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}

	got = interpret.Last()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
}

func TestCurrent(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	if got := interpret.Current(); len(got) > 0 {
		t.Errorf("Current must 0 len got=%v", got)
	}

	want := interpret.Next()

	if ct := interpret.Current(); len(ct) == 0 {
		t.Errorf("Current must > 0 len got=%v", ct)
	} else {
		got := ct[0]
		if !compareToken(t, 0, want, got) {
			fmt.Printf("want %#v\n", want)
			fmt.Printf("got %#v\n", got)
			return
		}

		want1 := interpret.Next()
		if ct = interpret.Current(); len(ct) != 2 {
			t.Errorf("Current must 2 len got=%v", ct)
		} else {
			got = ct[0]
			if !compareToken(t, 0, want, got) {
				fmt.Printf("want %#v\n", want)
				fmt.Printf("got %#v\n", got)
			}
			got = ct[1]
			if !compareToken(t, 0, want1, got) {
				fmt.Printf("want %#v\n", want1)
				fmt.Printf("got %#v\n", got)
			}
		}
	}
}

func TestFlush(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	interpret.Flush()
	want := &TokenWithPos{Token: lexer.Token{Type: 17, Value: "func"}, Pos: TokenPos{Line: 1, Pos: 0}}
	got := interpret.Next()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
	interpret.Flush()
	want = &TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 1, Pos: 4}}
	got = interpret.Next()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
	want = &TokenWithPos{Token: lexer.Token{Type: 3, Value: "tomate"}, Pos: TokenPos{Line: 1, Pos: 5}}
	got = interpret.Next()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
	interpret.Flush()
	want = &TokenWithPos{Token: lexer.Token{Type: 20, Value: "("}, Pos: TokenPos{Line: 1, Pos: 11}}
	got = interpret.Next()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
	want = &TokenWithPos{Token: lexer.Token{Type: 21, Value: ")"}, Pos: TokenPos{Line: 1, Pos: 12}}
	got = interpret.Next()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
	interpret.Rewind()
	interpret.Rewind()
	want = &TokenWithPos{Token: lexer.Token{Type: 20, Value: "("}, Pos: TokenPos{Line: 1, Pos: 11}}
	got = interpret.Next()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
	interpret.Flush()
	want = &TokenWithPos{Token: lexer.Token{Type: 21, Value: ")"}, Pos: TokenPos{Line: 1, Pos: 12}}
	got = interpret.Next()
	if !compareToken(t, 0, want, got) {
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
	}
}

func TestEmit(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	if got := interpret.Emit(); len(got) > 0 {
		t.Errorf("Emit must 0 len got=%v", got)
	}

	want := interpret.Next()

	if ct := interpret.Emit(); len(ct) != 1 {
		t.Errorf("Emit must 1 len got=%v", ct)
	} else {
		got := ct[0]
		if !compareToken(t, 0, want, got) {
			fmt.Printf("want %#v\n", want)
			fmt.Printf("got %#v\n", got)
		}

		want = interpret.Next()
		if ct = interpret.Emit(); len(ct) != 1 {
			t.Errorf("Emit must 1 len got=%v", ct)
		} else {
			got = ct[0]
			if !compareToken(t, 0, want, got) {
				fmt.Printf("want %#v\n", want)
				fmt.Printf("got %#v\n", got)
			}
		}
	}
}

func TestPeekOne(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	want := &TokenWithPos{Token: lexer.Token{Type: 17, Value: "func"}, Pos: TokenPos{Line: 1, Pos: 0}}
	got := interpret.PeekOne()
	compareToken(t, 0, want, got)

	igot := len(interpret.Current())
	iwant := 0
	if igot != iwant {
		t.Errorf("Wrong number of current tokens, got=%v want=%v", igot, iwant)
		return
	}

	interpret.ReadN(2)

	want = &TokenWithPos{Token: lexer.Token{Type: 3, Value: "tomate"}, Pos: TokenPos{Line: 1, Pos: 5}}
	got = interpret.PeekOne()
	compareToken(t, 0, want, got)

	igot = len(interpret.Current())
	iwant = 2
	if igot != iwant {
		t.Errorf("Wrong number of current tokens, got=%v want=%v", igot, iwant)
		return
	}
}

func TestPeekN(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	want := &TokenWithPos{Token: lexer.Token{Type: 17, Value: "func"}, Pos: TokenPos{Line: 1, Pos: 0}}
	got := interpret.PeekN(1)
	compareToken(t, 0, want, got[0])

	igot := len(interpret.Current())
	iwant := 0
	if igot != iwant {
		t.Errorf("Wrong number of current tokens, got=%v want=%v", igot, iwant)
		return
	}

	interpret.ReadN(2)

	want = &TokenWithPos{Token: lexer.Token{Type: 3, Value: "tomate"}, Pos: TokenPos{Line: 1, Pos: 5}}
	got = interpret.PeekN(1)
	compareToken(t, 0, want, got[0])

	igot = len(interpret.Current())
	iwant = 2
	if igot != iwant {
		t.Errorf("Wrong number of current tokens, got=%v want=%v", igot, iwant)
		return
	}
}

func TestPeek(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	got := interpret.Peek()
	compareToken(t, 0, nil, got)

	got = interpret.Peek(17)
	want := &TokenWithPos{Token: lexer.Token{Type: 17, Value: "func"}, Pos: TokenPos{Line: 1, Pos: 0}}
	compareToken(t, 0, want, got)

	interpret.Next()

	got = interpret.Peek(17)
	compareToken(t, 0, nil, got)

	interpret.Emit()

	want = &TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 1, Pos: 4}}
	got = interpret.Peek(0)
	compareToken(t, 0, want, got)

	igot := len(interpret.Current())
	iwant := 0
	if igot != iwant {
		t.Errorf("Wrong number of current tokens, got=%v want=%v", igot, iwant)
		return
	}

	interpret.ReadN(2)

	want = &TokenWithPos{Token: lexer.Token{Type: 20, Value: "("}, Pos: TokenPos{Line: 1, Pos: 11}}
	got = interpret.Peek(20)
	compareToken(t, 0, want, got)

	igot = len(interpret.Current())
	iwant = 2
	if igot != iwant {
		t.Errorf("Wrong number of current tokens, got=%v want=%v", igot, iwant)
		return
	}
}

func TestPeekUntil(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	got := interpret.PeekUntil()
	compareToken(t, 0, nil, got)

	got = interpret.PeekOne()
	want := &TokenWithPos{Token: lexer.Token{Type: 17, Value: "func"}, Pos: TokenPos{Line: 1, Pos: 0}}
	compareToken(t, 0, want, got)

	got = interpret.PeekUntil(17)
	want = &TokenWithPos{Token: lexer.Token{Type: 17, Value: "func"}, Pos: TokenPos{Line: 1, Pos: 0}}
	compareToken(t, 0, want, got)

	interpret.Next()

	got = interpret.PeekUntil(17)
	compareToken(t, 0, nil, got)

	interpret.Emit()

	want = &TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 1, Pos: 4}}
	got = interpret.PeekUntil(0)
	compareToken(t, 0, want, got)

	igot := len(interpret.Current())
	iwant := 0
	if igot != iwant {
		t.Errorf("Wrong number of current tokens, got=%v want=%v", igot, iwant)
		return
	}

	interpret.ReadN(2)

	want = &TokenWithPos{Token: lexer.Token{Type: 20, Value: "("}, Pos: TokenPos{Line: 1, Pos: 11}}
	got = interpret.PeekUntil(20)
	compareToken(t, 0, want, got)

	igot = len(interpret.Current())
	iwant = 2
	if igot != iwant {
		t.Errorf("Wrong number of current tokens, got=%v want=%v", igot, iwant)
		return
	}
}

func TestRead(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	got := interpret.Read()
	compareToken(t, 0, nil, got)

	want := interpret.PeekOne()

	got = interpret.Read(glanglexer.FuncToken)
	compareToken(t, 0, want, got)

	got = interpret.PeekOne()
	want = &TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 1, Pos: 4}}
	compareToken(t, 0, want, got)

	got = interpret.Read(3, 0)
	compareToken(t, 0, want, got)

	got = interpret.PeekOne()
	want = &TokenWithPos{Token: lexer.Token{Type: 3, Value: "tomate"}, Pos: TokenPos{Line: 1, Pos: 5}}
	compareToken(t, 0, want, got)

	interpret.Rewind()
	got = interpret.PeekOne()
	want = &TokenWithPos{Token: lexer.Token{Type: 0, Value: " "}, Pos: TokenPos{Line: 1, Pos: 4}}
	compareToken(t, 0, want, got)
}

func TestReadMany(t *testing.T) {

	str := `func tomate () {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	compareTokensLen(t, "ReadMany", 0, 0, interpret.ReadMany())
	compareTokensLen(t, "ReadMany", 0, 0, interpret.ReadMany(0))
	compareTokensLen(t, "ReadMany", 0, 1, interpret.ReadMany(glanglexer.FuncToken))
	compareTokensLen(t, "ReadMany", 0, 0, interpret.ReadMany(glanglexer.BraceOpenToken))
	compareTokensLen(t, "ReadMany", 0, 4, interpret.ReadMany(glanglexer.ParenOpenToken, 3, 0))
	compareTokensLen(t, "Current", 0, 5, interpret.Current())
}

func TestGet(t *testing.T) {

	str := `func tomate () {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	got := interpret.Get()
	compareToken(t, 0, nil, got)

	want := interpret.PeekOne()

	got = interpret.Get(glanglexer.FuncToken)
	compareToken(t, 0, want, got)

	want = interpret.PeekOne()

	interpret.Rewind()
	got = interpret.PeekOne()
	compareToken(t, 0, want, got)
}

func TestGetMany(t *testing.T) {

	str := `func tomate () {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	compareTokensLen(t, "GetMany", 0, 0, interpret.GetMany())
	compareTokensLen(t, "GetMany", 0, 0, interpret.GetMany(0))
	compareTokensLen(t, "GetMany", 0, 1, interpret.GetMany(glanglexer.FuncToken))
	compareTokensLen(t, "GetMany", 0, 0, interpret.GetMany(glanglexer.BraceOpenToken))
	compareTokensLen(t, "GetMany", 0, 4, interpret.GetMany(glanglexer.ParenOpenToken, 3, 0))
	compareTokensLen(t, "Current", 0, 0, interpret.Current())
}

func TestReadBlock(t *testing.T) {

	str := `func tomate (xx (yy)) {
		var expr string = "some"
		expr2 := "other"
}`
	d := stringTokenizer(str)
	interpret := NewInterpreter(d)

	tokens := interpret.ReadBlock(glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	compareTokensLen(t, "ReadBlock", 0, 0, tokens)

	interpret.GetMany(glanglexer.FuncToken, 0, 3)

	tokens = interpret.ReadBlock(glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	compareTokensLen(t, "ReadBlock", 0, 7, tokens)
}

func compareTokensLen(t *testing.T, r string, c int, want int, got []Tokener) bool {
	return compareLen(t, r+" tokens", c, want, len(got))
}

func compareLen(t *testing.T, r string, c int, want, got int) bool {
	if got != want {
		t.Errorf("Wrong number of %v at %v, got=%v want=%v", r, c, got, want)
		return false
	}
	return true
}

func compareToken(t *testing.T, c int, want, got Tokener) bool {

	if got == nil || got == Tokener(nil) {
		if want != nil {
			t.Errorf("Wrong result at %v want=%v, got=%v", c, want, got)
		}
		return want == nil
	}

	if want == nil || want == Tokener(nil) {
		if got != nil {
			t.Errorf("Wrong result at %v want=%v, got=%v", c, want, got)
		}
		return got == nil
	}

	twant := want.GetType()
	tgot := got.GetType()
	if twant != tgot {
		t.Errorf("Wrong token type at %v want=%v, got=%v", c, twant, tgot)
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
		return false
	}
	swant := want.GetValue()
	sgot := got.GetValue()
	if swant != sgot {
		t.Errorf("Wrong token value at %v want=%v, got=%v", c, swant, sgot)
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
		return false
	}
	swant = want.GetPos().String()
	sgot = got.GetPos().String()
	if swant != sgot {
		t.Errorf("Wrong token pos at %v want=%v, got=%v", c, swant, sgot)
		fmt.Printf("want %#v\n", want)
		fmt.Printf("got %#v\n", got)
		return false
	}
	return true
}

func stringTokenizer(content string) TokenerReader {

	var buf bytes.Buffer
	buf.WriteString(content)

	return makeLexerReader(&buf)
}

func makeLexerReader(r io.Reader) TokenerReader {
	l := lexer.New(r, (glanglexer.New()).StartHere)
	l.ErrorHandler = func(e string) {}

	reader := NewReadTokenWithPos(l)

	return reader
}
