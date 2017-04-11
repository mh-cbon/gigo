package generic

import (
	"fmt"

	genericlexer "github.com/mh-cbon/gigo/lexer/generic"
	lexer "github.com/mh-cbon/state-lexer"
)

// ExprReceiver receives tokens
type ExprReceiver interface {
	AddExpr(expr Tokener)
	AddExprs(expr []Tokener)
}

// ScopeReceiver is a context of token reception
type ScopeReceiver interface {
	ExprReceiver
	GetName() string
	FinalizeErr(*ParseError) error
}

// TokenPos is a token position
type TokenPos struct {
	Line int
	Pos  int
}

func (t TokenPos) String() string {
	return fmt.Sprintf("%3d:%-3d", t.Line, t.Pos)
}

// TokenWithPos is a token with a pos.
type TokenWithPos struct {
	lexer.Token
	Pos TokenPos
}

//SetValue implements Tokener.
func (f *TokenWithPos) SetValue(s string) {
	f.Value = s
}

//SetType implements Tokener.
func (f *TokenWithPos) SetType(s lexer.TokenType) {
	f.Type = s
}

//GetExprs implements Expressioner. its a noop.
func (f *TokenWithPos) GetExprs() []Expressioner {
	return []Expressioner{}
}

//GetTokens implements Expressioner. returns an array of one token.
func (f *TokenWithPos) GetTokens() []Tokener {
	return []Tokener{f}
}

//GetToken implements Expressioner. returns the token if its type matches T.
func (f *TokenWithPos) GetToken(T lexer.TokenType) Tokener {
	if f.GetType() == T {
		return f
	}
	return nil
}

//GetPos implements Expressioner.
func (f *TokenWithPos) GetPos() TokenPos {
	return f.Pos
}

//HasToken implements Expressioner. returns true if token type is T.
func (f *TokenWithPos) HasToken(T lexer.TokenType) bool {
	return f.Type == T
}

//Remove implements Expressioner. It is a noop.
func (f *TokenWithPos) Remove(e Expressioner) bool { return false }

//PrependExpr implements Expressioner. It is a noop.
func (f *TokenWithPos) PrependExpr(expr Tokener) {}

//PrependExprs implements Expressioner. It is a noop.
func (f *TokenWithPos) PrependExprs(exprs []Tokener) {}

//SetTokenValue implements Expressioner. It changes value of the token if its type matches T.
func (f *TokenWithPos) SetTokenValue(T lexer.TokenType, v string) {
	if f.GetType() == T {
		f.SetValue(v)
	}
}

//GetTokensAtLine implements Expressioner. It returns the tokens if its pos matches line.
func (f *TokenWithPos) GetTokensAtLine(line int) []Tokener {
	if f.Pos.Line == line {
		return []Tokener{f}
	}
	return []Tokener{}
}

//FindAll implements Expressioner. It returns the token if its type matches T.
func (f *TokenWithPos) FindAll(T lexer.TokenType) []Expressioner {
	if f.GetType() == T {
		return []Expressioner{f}
	}
	return []Expressioner{}
}

// NewTokenWithPos creates a positionned Tokener of token.
func NewTokenWithPos(t lexer.Token, line, pos int) *TokenWithPos {
	return &TokenWithPos{
		Token: t,
		Pos: TokenPos{
			Pos:  pos,
			Line: line,
		},
	}
}

// NewTokenEOF creates an EOF Tokener.
func NewTokenEOF() *TokenWithPos {
	return NewTokenWithPos(lexer.Token{Type: genericlexer.EOFToken}, -1, -1)
}

// Tokener is an interface that defiens methods to change value of an underlyong token.
type Tokener interface {
	GetType() lexer.TokenType
	GetValue() string
	SetValue(string)
	SetType(lexer.TokenType)
	String() string
	GetPos() TokenPos
}

// Expressioner is an interface that defines methods to manipulate a token of tokens.
type Expressioner interface {
	GetExprs() []Expressioner
	GetTokens() []Tokener
	HasToken(lexer.TokenType) bool
	Remove(Expressioner) bool
	String() string
	PrependExpr(Tokener)
	PrependExprs([]Tokener)
	GetToken(lexer.TokenType) Tokener
	SetTokenValue(lexer.TokenType, string)
	GetTokensAtLine(int) []Tokener
	FindAll(lexer.TokenType) []Expressioner
}

// Expression isa  token of tokens.
type Expression struct {
	Tokens []Tokener
}

// First get 1st.
func (f *Expression) First() Tokener {
	if len(f.Tokens) > 0 {
		return f.Tokens[0]
	}
	return nil
}

// Until returns all tokens until it met u.
func (f *Expression) Until(u lexer.TokenType) *Expression {
	ret := &Expression{}
	for _, t := range f.Tokens {
		if t.GetType() == u {
			break
		} else {
			ret.Tokens = append(ret.Tokens, t)
		}
	}
	return ret
}

// GetPos get 1st pos.
func (f *Expression) GetPos() TokenPos {
	return f.First().GetPos() //let see if need to be pointer
}

// GetType get 1st TokenType.
func (f *Expression) GetType() lexer.TokenType {
	return f.First().GetType() //let see if need to be pointer
}

// GetValue get 1st value.
func (f *Expression) GetValue() string {
	return f.First().GetValue() //let see if need to be pointer
}

// SetType set 1st TokenType.
func (f *Expression) SetType(T lexer.TokenType) {
	f.First().SetType(T) //let see if need to be pointer
}

// SetValue set 1st value.
func (f *Expression) SetValue(v string) {
	f.First().SetValue(v) //let see if need to be pointer
}

// GetTokensAtLine finds all tokens at line.
func (f *Expression) GetTokensAtLine(line int) []Tokener {
	ret := []Tokener{}
	for _, t := range f.GetExprs() {
		ret = append(ret, t.GetTokensAtLine(line)...)
	}
	return ret
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
		if t.(Expressioner) == e {
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

// MustInsertAfter ...
func (f *Expression) MustInsertAfter(ref Expressioner, nnew Tokener) {
	if !f.InsertAfter(ref, nnew) {
		panic("MustInsertAfter failed")
	}
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

// MustRemove ...
func (f *Expression) MustRemove(e Expressioner) {
	if !f.Remove(e) {
		panic("MustRemove failed")
	}
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

// FindAll returns all expressions of type T.
func (f *Expression) FindAll(T lexer.TokenType) []Expressioner {
	var ret []Expressioner
	for _, t := range f.GetExprs() {
		ret = append(ret, t.FindAll(T)...)
	}
	return ret
}

// AddExpr appends a Tokener
func (f *Expression) AddExpr(expr Tokener) {
	f.Tokens = append(f.Tokens, expr)
}

// AddExprs appends many Tokener
func (f *Expression) AddExprs(exprs []Tokener) {
	f.Tokens = append(f.Tokens, exprs...)
}

// PrependExpr prepends a Tokener
func (f *Expression) PrependExpr(expr Tokener) {
	f.Tokens = append([]Tokener{expr}, f.Tokens...)
}

// PrependExprs prepends many Tokener
func (f *Expression) PrependExprs(exprs []Tokener) {
	f.Tokens = append(exprs, f.Tokens...)
}
func (f *Expression) String() string {
	s := ""
	for _, e := range f.Tokens {
		s += e.String()
	}
	return s
}
