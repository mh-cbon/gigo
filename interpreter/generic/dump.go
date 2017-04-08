package generic

import (
	"fmt"
	"strings"

	genericlexer "github.com/mh-cbon/gigo/lexer/generic"
)

// Dump recursively prints an expression
func Dump(src Expressioner) {
	DumpWithNamer(src, genericlexer.TokenType, 0)
}
func DumpWithNamer(src Expressioner, typer TokenTyper, lvl int) {
	x := strings.Repeat(" ", lvl)
	T := fmt.Sprintf("%T", src)
	T = strings.Replace(T, "*glang.", "", -1)
	exprs := src.GetExprs()
	for i, e := range exprs {
		if len(e.GetExprs()) == 0 {
			tok := e.GetTokens()[0]
			yy := ""
			if len(src.GetExprs()) == 1 {
				yy = fmt.Sprintf("=> %v %v token", T, len(src.GetExprs()))
			} else if i == 0 {
				yy = fmt.Sprintf("-> %v %v tokens", T, len(src.GetExprs()))
			} else if i == len(src.GetExprs())-1 {
				yy = fmt.Sprintf("<- %v", T)
			}
			fmt.Printf("%-40v", x+yy)

			fmt.Printf("%-8v%-20v %q\n",
				tok.GetPos().String(),
				typer(tok.GetType()),
				tok.GetValue())
		} else {
			if i == 0 {
				yy := fmt.Sprintf("-> %v %v tokens", T, len(src.GetExprs()))
				fmt.Printf("%-40v\n", x+yy)
			}
			DumpWithNamer(e, typer, lvl+1)
			if i == len(src.GetExprs())-1 {
				yy := fmt.Sprintf("<- %v %v tokens", T, len(src.GetExprs()))
				fmt.Printf("%-40v\n", x+yy)
			}
		}
	}

	if len(src.GetExprs()) < 1 || len(src.GetExprs()) > 2000 {
		yy := fmt.Sprintf("<- %v %v tokens", T, len(src.GetExprs()))
		fmt.Printf("%-40v\n", x+yy)
	}
	// fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "<-", src, len(src.GetExprs()))
}

// DumpTokens prints a list of tokens
func DumpTokens(tokens []Tokener) {
	DumpTokensWithNamer(tokens, genericlexer.TokenType)
}
func DumpTokensWithNamer(tokens []Tokener, typer TokenTyper) {
	x := strings.Repeat(" ", 2)
	fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "begin", "<noname>", len(tokens))
	for _, tok := range tokens {
		fmt.Printf("%40v %v %20v %q\n", x, tok.GetPos().String(),
			typer(tok.GetType()),
			tok.GetValue())
	}
	fmt.Printf("%v%-6v %-20T tokens(%v)\n", x, "end", "<noname>", len(tokens))
}
