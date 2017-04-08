package generic

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	lexer "github.com/mh-cbon/state-lexer"
)

type TokenReaderOK interface {
	NextToken() *lexer.Token
}

type TokenerReaderOK interface {
	NextToken() Tokener
}

type ReadTokenWithPos struct {
	Reader TokenReaderOK
	line   int
	pos    int
}

func (r *ReadTokenWithPos) NextToken() Tokener {
	if next := r.Reader.NextToken(); next != nil {
		tok := NewTokenWithPos(*next, r.line, r.pos)
		if k := strings.Count(tok.Value, "\n"); k > 0 {
			r.line += k
			r.pos = -1
		}
		r.pos += len(tok.Value)
		return tok
	}
	return nil
}

func NewReadTokenWithPos(r TokenReaderOK) *ReadTokenWithPos {
	return &ReadTokenWithPos{
		Reader: r,
		line:   1,
		pos:    0,
	}
}

type ReadNPrettyPrint struct {
	Reader TokenerReaderOK
	Namer  TokenerNamer
	w      *tabwriter.Writer
	d      bool
}

func (r *ReadNPrettyPrint) Flush() {
	r.w.Flush()
}
func (r *ReadNPrettyPrint) NextToken() Tokener {
	if !r.d {
		r.d = true
		fmt.Fprintf(r.w, "%v\t %v\t %q\n", "Line:Pos", "TokenName:tok.Type", "tok.Value")
		r.Flush()
	}
	if next := r.Reader.NextToken(); next != nil {
		fmt.Fprintf(r.w, "%0000d:%000d\t %v:%000d\t %q\n",
			next.GetPos().Line,
			next.GetPos().Pos, r.Namer(next), next.GetType(), next.GetValue())
		r.Flush()
		return next
	} else if next == nil {
		fmt.Fprintf(r.w, "%v\t %v\t %q\n", "Line:Pos", "TokenName:tok.Type", "tok.Value")
		r.Flush()
	}
	return nil
}

func NewReadNPrettyPrint(r TokenerReaderOK, Namer TokenerNamer, out io.Writer) *ReadNPrettyPrint {
	return &ReadNPrettyPrint{
		Reader: r,
		Namer:  Namer,
		w:      tabwriter.NewWriter(out, 0, 0, 1, ' ', tabwriter.Debug),
	}
}

// TokenNamer gives the name of a Token.
type TokenNamer func(lexer.Token) string

// TokenerNamer gives the name of a Tokener.
type TokenerNamer func(Tokener) string

// TokenerName turns a TokenNamer into a TokenerNamer.
func TokenerName(h TokenNamer) TokenerNamer {
	return func(t Tokener) string {
		return h(lexer.Token{Type: t.GetType(), Value: t.GetValue()})
	}
}

// TokenTyper gives the name of a Token type.
type TokenTyper func(lexer.TokenType) string

// ------------------
// old stuff

// // TokenReader return Token or nil.
// type TokenReader func() *lexer.Token
//
// // TokenerReader return Tokener or nil.
// type TokenerReader func() Tokener
//
// // PositionnedTokenReader turns a TokenReader into a feed of Tokener.
// func PositionnedTokenReader(reader TokenReader) func() Tokener {
// 	var line = 1
// 	var pos = 0
// 	return func() Tokener {
// 		if next := reader(); next != nil {
// 			tok := NewTokenWithPos(*next, line, pos)
// 			if k := strings.Count(tok.Value, "\n"); k > 0 {
// 				line += k
// 				pos = -1
// 			}
// 			pos += len(tok.Value)
// 			return tok
// 		}
// 		return nil
// 	}
// }
//
// // PrettyPrint a feed of token being read.
// func PrettyPrint(reader TokenerReader, namer TokenerNamer) func() Tokener {
// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
// 	fmt.Fprintf(w, "%v\t %v\t %q\n", "Line:Pos", "TokenName:tok.Type", "tok.Value")
// 	return func() Tokener {
// 		if next := reader(); next != nil {
// 			fmt.Fprintf(w, "%0000d:%000d\t %v:%000d\t %q\n",
// 				next.GetPos().Line,
// 				next.GetPos().Pos, namer(next), next.GetType(), next.GetValue())
// 			return next
// 		} else if next == nil {
// 			fmt.Fprintf(w, "%v\t %v\t %q\n", "Line:Pos", "TokenName:tok.Type", "tok.Value")
// 			w.Flush()
// 		}
// 		return nil
// 	}
// }
