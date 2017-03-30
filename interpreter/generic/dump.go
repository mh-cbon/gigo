package generic

import (
	"fmt"
	"strings"

	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
)

// type PrettyPrint struct {
// }
//
// func (p PrettyPrint) Format(s fmt.State, verb rune) {
// 	switch verb {
// 	case 'v':
// 		if s.Flag('+') {
// 			io.WriteString(s, p)
// 			return
// 		}
// 		fallthrough
// 	case 's':
// 		io.WriteString(s, p)
// 	case 'q':
// 		fmt.Fprintf(s, "%q", p)
// 	}
// }

func Dump(src Expressioner, lvl int) {
	x := strings.Repeat(" ", lvl)
	fmt.Printf("%v%-6v %-20T Tokens(%v)\n", x, "begin", src, len(src.GetExprs()))
	exprs := src.GetExprs()
	for _, e := range exprs {
		if len(e.GetExprs()) == 0 {
			tok := e.GetTokens()[0]
			fmt.Printf("%40v %v %20v %q\n", x, tok.GetPos().String(),
				glanglexer.TokenType(tok.GetType()),
				tok.GetValue())
		} else {
			Dump(e, lvl+1)
		}
	}
	fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "end", src, len(src.GetExprs()))
}

func DumpTokens(tokens []Tokener, lvl int) {
	x := strings.Repeat(" ", lvl)
	fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "begin", "<noname>", len(tokens))
	for _, tok := range tokens {
		fmt.Printf("%40v %v %20v %q\n", x, tok.GetPos().String(),
			glanglexer.TokenType(tok.GetType()),
			tok.GetValue())
	}
	fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "end", "<noname>", len(tokens))
}
