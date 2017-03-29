package generic

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	lexer "github.com/mh-cbon/state-lexer"
)

type TokenReader func() *lexer.Token
type TokenerReader func() Tokener

type TokenNamer func(lexer.Token) string
type TokenerNamer func(Tokener) string

func TokenerName(h TokenNamer) TokenerNamer {
	return func(t Tokener) string {
		return h(lexer.Token{Type: t.GetType(), Value: t.GetValue()})
	}
}

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
