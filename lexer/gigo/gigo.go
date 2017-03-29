package gigo

import (
	generic "github.com/mh-cbon/gigo/lexer/generic"
	glang "github.com/mh-cbon/gigo/lexer/glang"
	lexer "github.com/mh-cbon/state-lexer"
)

const (
	xxxToken lexer.TokenType = iota
)

// TokenName Helper function
func TokenName(tok lexer.Token) string {
	return TokenType(tok.Type)
}

// TokenType Helper function
func TokenType(Type lexer.TokenType) string {
	return glang.TokenType(Type)
}

// New ...
func New() *generic.Lexer {
	ret := glang.New()
	ret.Printer = TokenType
	return ret
}
