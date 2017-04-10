package generic

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	lexer "github.com/mh-cbon/state-lexer"
)

// TokenReader reads *lexer.Token
type TokenReader interface {
	NextToken() *lexer.Token
}

// TokenerReader reads Tokener
type TokenerReader interface {
	NextToken() Tokener
}

// ReadTokenWithPos reads lexer.Token, outputs Tokener
type ReadTokenWithPos struct {
	Reader TokenReader
	line   int
	pos    int
}

// NextToken advance to the next token.
// returns nil on eof.
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

// NewReadTokenWithPos initializes a TokenerReader from a TokenerReader
func NewReadTokenWithPos(r TokenReader) *ReadTokenWithPos {
	return &ReadTokenWithPos{
		Reader: r,
		line:   1,
		pos:    0,
	}
}

// ReadNPrettyPrint is a TokenerReader that pretty prints what it reads.
type ReadNPrettyPrint struct {
	Reader TokenerReader
	Namer  TokenerNamer
	w      *tabwriter.Writer
	d      bool
}

// Flush prints on the underlying write.
func (r *ReadNPrettyPrint) Flush() {
	r.w.Flush()
}

// NextToken returns next token and prints it.
func (r *ReadNPrettyPrint) NextToken() Tokener {
	if !r.d {
		r.d = true
		fmt.Fprintf(r.w, "%v\t %v\t %q\n", "Line:Pos", "TokenName:tok.Type", "tok.Value")
	}
	if next := r.Reader.NextToken(); next != nil {
		fmt.Fprintf(r.w, "%0000d:%000d\t %v:%000d\t %q\n",
			next.GetPos().Line,
			next.GetPos().Pos, r.Namer(next), next.GetType(), next.GetValue())
		return next
	} else if next == nil {
		fmt.Fprintf(r.w, "%v\t %v\t %q\n", "Line:Pos", "TokenName:tok.Type", "tok.Value")
		r.Flush()
	}
	return nil
}

// NewReadNPrettyPrint initializes a prety printer reader to out, Namer resolves token types.
func NewReadNPrettyPrint(r TokenerReader, Namer TokenerNamer, out io.Writer) *ReadNPrettyPrint {
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
