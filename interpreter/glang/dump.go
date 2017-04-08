package glang

import (
	genericinterpreter "github.com/mh-cbon/gigo/interpreter/generic"
	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
)

// Dump recursively prints an expression
func Dump(src genericinterpreter.Expressioner) {
	genericinterpreter.DumpWithNamer(src, glanglexer.TokenType, 0)
}

// DumpTokens prints a list of tokens
func DumpTokens(tokens []genericinterpreter.Tokener) {
	genericinterpreter.DumpTokensWithNamer(tokens, glanglexer.TokenType)
}
