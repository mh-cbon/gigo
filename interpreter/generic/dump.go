package generic

import (
	"fmt"
	"strings"

	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
)

// Dump recursively prints an expression
func Dump(src Expressioner) {
	dump(src, 0)
}
func dump(src Expressioner, lvl int) {
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
			dump(e, lvl+1)
		}
	}
	fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "end", src, len(src.GetExprs()))
}

// DumpTokens prints a list of tokens
func DumpTokens(tokens []Tokener) {
	dumpTokens(tokens, 0)
}
func dumpTokens(tokens []Tokener, lvl int) {
	x := strings.Repeat(" ", lvl)
	fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "begin", "<noname>", len(tokens))
	for _, tok := range tokens {
		fmt.Printf("%40v %v %20v %q\n", x, tok.GetPos().String(),
			glanglexer.TokenType(tok.GetType()),
			tok.GetValue())
	}
	fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "end", "<noname>", len(tokens))
}
