package generic

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	lexer "github.com/mh-cbon/state-lexer"
)

// TokenReader return Token or nil.
type TokenReader func() *lexer.Token

// TokenerReader return Tokener or nil.
type TokenerReader func() Tokener

// PositionnedTokenReader turns a TokenReader into a feed of Tokener.
func PositionnedTokenReader(reader TokenReader) func() Tokener {
	var line = 1
	var pos = 0
	return func() Tokener {
		if next := reader(); next != nil {
			tok := NewTokenWithPos(*next, line, pos)
			if k := strings.Count(tok.Value, "\n"); k > 0 {
				line += k
				pos = -1
			}
			pos += len(tok.Value)
			return tok
		}
		return nil
	}
}

// PrettyPrint a feed of token being read.
func PrettyPrint(reader TokenerReader, namer TokenerNamer) func() Tokener {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "%v\t %v\t %q\n", "Line:Pos", "TokenName:tok.Type", "tok.Value")
	return func() Tokener {
		if next := reader(); next != nil {
			fmt.Fprintf(w, "%0000d:%000d\t %v:%000d\t %q\n",
				next.GetPos().Line,
				next.GetPos().Pos, namer(next), next.GetType(), next.GetValue())
			return next
		} else if next == nil {
			fmt.Fprintf(w, "%v\t %v\t %q\n", "Line:Pos", "TokenName:tok.Type", "tok.Value")
			w.Flush()
		}
		return nil
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
