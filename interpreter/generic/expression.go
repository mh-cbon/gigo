package generic

import (
	"fmt"

	genericlexer "github.com/mh-cbon/gigo/lexer/generic"
	lexer "github.com/mh-cbon/state-lexer"
)

type ExprReceiver interface {
	AddExpr(expr Tokener)
	AddExprs(expr []Tokener)
}
type ScopeReceiver interface {
	ExprReceiver
	GetName() string
	FinalizeErr(*SyntaxError) error
}

type TokenPos struct {
	Line int
	Pos  int
}

func (t TokenPos) String() string {
	return fmt.Sprintf("%3d:%3d", t.Line, t.Pos)
}

type TokenWithPos struct {
	lexer.Token
	Pos TokenPos
}

func (f *TokenWithPos) SetValue(s string) {
	f.Value = s
}
func (f *TokenWithPos) SetType(s lexer.TokenType) {
	f.Type = s
}
func (f *TokenWithPos) GetExprs() []Expressioner {
	return []Expressioner{}
}
func (f *TokenWithPos) GetTokens() []Tokener {
	return []Tokener{f}
}
func (f *TokenWithPos) GetToken(T lexer.TokenType) Tokener {
	if f.GetType() == T {
		return f
	}
	return nil
}
func (f *TokenWithPos) GetPos() TokenPos {
	return f.Pos
}
func (f *TokenWithPos) HasToken(T lexer.TokenType) bool {
	return f.Type == T
}
func (f *TokenWithPos) Remove(e Expressioner) bool {
	return false
}
func (f *TokenWithPos) PrependExpr(expr Tokener) {
}
func (f *TokenWithPos) PrependExprs(exprs []Tokener) {
}
func (f *TokenWithPos) SetTokenValue(T lexer.TokenType, v string) {
	if f.GetType() == T {
		f.SetValue(v)
	}
}

func NewTokenWithPos(t lexer.Token, line, pos int) *TokenWithPos {
	return &TokenWithPos{
		Token: t,
		Pos: TokenPos{
			Pos:  pos,
			Line: line,
		},
	}
}

func NewTokenEOF() *TokenWithPos {
	return NewTokenWithPos(lexer.Token{Type: genericlexer.EOFToken}, -1, -1)
}

type Tokener interface {
	GetType() lexer.TokenType
	GetValue() string
	SetValue(string)
	SetType(lexer.TokenType)
	String() string
	GetPos() TokenPos
}

type Expressioner interface {
	GetExprs() []Expressioner
	GetTokens() []Tokener
	HasToken(lexer.TokenType) bool
	Remove(Expressioner) bool
	PrependExpr(Tokener)
	PrependExprs([]Tokener)
	GetToken(T lexer.TokenType) Tokener
	SetTokenValue(T lexer.TokenType, v string)
}

type Expression struct {
	Tokens []Tokener
}

// Filter root tokens of type T.
func (f *Expression) Filter(T lexer.TokenType) []Tokener {
	var ret []Tokener
	for _, t := range f.Tokens {
		if t.GetType() == T {
			ret = append(ret, t)
		}
	}
	return ret
}

// FilterToken filter root tokens of type T and returns the first one.
func (f *Expression) FilterToken(T lexer.TokenType) Tokener {
	t := f.Filter(T)
	if len(t) > 0 {
		return t[0]
	}
	return nil
}

// GetToken implements Tokener.
func (f *Expression) GetToken(T lexer.TokenType) Tokener {
	return f.FilterToken(T)
}

// GetExprIndex returns index of a root token matching given expression.
func (f *Expression) GetExprIndex(e Expressioner) int {
	for i, t := range f.Tokens {
		if t == e.(Tokener) {
			return i
		}
	}
	return -1
}

// GetTokenIndex returns index of a root token matching given Tokener.
func (f *Expression) GetTokenIndex(e lexer.TokenType) int {
	for i, t := range f.Tokens {
		if t.GetType() == e {
			return i
		}
	}
	return -1
}

// Replace a root token with given expression.
func (f *Expression) Replace(old Expressioner, nnew Tokener) bool {
	if index := f.GetExprIndex(old); index > -1 {
		f.Tokens[index] = nnew
		return true
	}
	return false
}

// InsertAfter a nnew root token after given ref expression.
func (f *Expression) InsertAfter(ref Expressioner, nnew Tokener) bool {
	if index := f.GetExprIndex(ref); index > -1 {
		f.InsertAt(index, nnew)
		return true
	}
	return false
}

// InsertAt a nnew root token at index.
func (f *Expression) InsertAt(index int, nnew Tokener) {
	f.Tokens = append(f.Tokens[:index], append([]Tokener{nnew}, f.Tokens[index:]...)...)
}

// RemoveAt removes root token at index.
func (f *Expression) RemoveAt(index int) {
	f.Tokens = append(f.Tokens[:index], f.Tokens[index+1:]...)
}

// RemoveT a root token of type T.
func (f *Expression) RemoveT(t lexer.TokenType) bool {
	if index := f.GetTokenIndex(t); index > -1 {
		f.Tokens = append(f.Tokens[:index], f.Tokens[index+1:]...)
		return true
	}
	return false
}

// Remove a root token of this expression.
func (f *Expression) Remove(e Expressioner) bool {
	if index := f.GetExprIndex(e); index > -1 {
		f.RemoveAt(index)
		return true
	}
	return false
}

// RemoveAll root tokens matching those expressions.
func (f *Expression) RemoveAll(e []Expressioner) int {
	ret := 0
	for _, t := range e {
		if f.Remove(t) {
			ret++
		}
	}
	return ret
}

// HasToken recursively for a token of the given type.
func (f *Expression) HasToken(T lexer.TokenType) bool {
	for _, t := range f.GetExprs() {
		if t.HasToken(T) {
			return true
		}
	}
	return false
}

// GetTokens returns the list of root-tokens.
func (f *Expression) GetTokens() []Tokener {
	return f.Tokens
}

// SetTokenValue recusrively change the Value of tokens of type T.
func (f *Expression) SetTokenValue(T lexer.TokenType, v string) {
	for _, t := range f.GetExprs() {
		t.SetTokenValue(T, v)
	}
}

// GetExprs returns a list of root expression.
func (f *Expression) GetExprs() []Expressioner {
	var ret []Expressioner
	for _, t := range f.Tokens {
		ret = append(ret, t.(Expressioner))
	}
	return ret
}
func (f *Expression) AddExpr(expr Tokener) {
	if expr == nil || expr == Tokener(nil) {
		panic("rrr")
	}
	f.Tokens = append(f.Tokens, expr)
}
func (f *Expression) AddExprs(exprs []Tokener) {
	for _, expr := range exprs {
		if expr == nil || expr == Tokener(nil) {
			panic("rrr")
		}
	}
	f.Tokens = append(f.Tokens, exprs...)
}
func (f *Expression) PrependExpr(expr Tokener) {
	if expr == nil || expr == Tokener(nil) {
		panic("rrr")
	}
	f.Tokens = append([]Tokener{expr}, f.Tokens...)
}
func (f *Expression) PrependExprs(exprs []Tokener) {
	for _, expr := range exprs {
		if expr == nil || expr == Tokener(nil) {
			panic("rrr")
		}
	}
	f.Tokens = append(exprs, f.Tokens...)
}
func (f *Expression) String() string {
	s := ""
	for _, e := range f.Tokens {
		s += e.String()
	}
	return s
}
