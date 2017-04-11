package generic

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
	lexer "github.com/mh-cbon/state-lexer"
)

// Interpreter navigates a tokens list to produce a tokens tree.
type Interpreter struct {
	isEnded  bool
	Namer    TokenTyper
	position int
	try      int
	Reader   TokenerReader
	Tokens   []Tokener
	Scope    ScopeReceiver
}

// NewInterpreter makes an Interpreter starting at -1
func NewInterpreter(r TokenerReader) *Interpreter {
	return &Interpreter{
		Namer:    TokenTyper(glanglexer.TokenType),
		Reader:   r,
		position: -1,
	}
}

// Next gives the next token.
func (I *Interpreter) Next() Tokener {
	if I.position < len(I.Tokens) {
		I.try++
		// this helps to avoid infinite loops, on the other hand it limits size of tokens.
		if I.try > 10000 {
			n := I.Tokens[len(I.Tokens)-1]
			err := I.DebugAtToken(n, "Inifinite loop detected")
			fmt.Printf("%+v\n\n", err)
			fmt.Printf("%#v", err)
			os.Exit(1)
		}
		I.position++
		if I.position < len(I.Tokens) {
			return I.Tokens[I.position]
		}
		for {
			n := I.Reader.NextToken()
			if n == nil {
				I.isEnded = true
				break
			}
			I.Tokens = append(I.Tokens, n)
			return n
		}
	}
	return nil
}

// Rewind returns to the previous token, if any.
func (I *Interpreter) Rewind() {
	I.position--
	if I.position < -1 {
		I.position = -1
	}
}

// RewindAll returns to first position
func (I *Interpreter) RewindAll() {
	I.position = -1
}

// Last is the token at current position.
func (I *Interpreter) Last() Tokener {
	if I.position > -1 && I.position < len(I.Tokens) {
		return I.Tokens[I.position]
	}
	return nil
}

// Ended when the tokens list is empty.
func (I *Interpreter) Ended() bool {
	return I.isEnded
}

// Current unemitted tokens.
func (I *Interpreter) Current() []Tokener {
	if I.position < 0 {
		return I.Tokens[:0]
	}
	if I.position >= len(I.Tokens) {
		return I.Tokens[0:]
	}
	return I.Tokens[:I.position+1]
}

// Emit current tokens in buffer.
func (I *Interpreter) Emit() []Tokener {
	I.try = 0
	toks := []Tokener{}
	c := I.Current()
	if len(c) > 0 {
		toks = append(toks, c...)
		I.Flush()
	}
	return toks
}

// Flush current tokens in buffer.
func (I *Interpreter) Flush() {
	if I.position+1 < len(I.Tokens) {
		I.Tokens = I.Tokens[I.position+1:]
	} else {
		I.Tokens = I.Tokens[:0]
	}
	I.position = -1
}

// PeekOne returns next token without changing position.
func (I *Interpreter) PeekOne() Tokener {
	t := I.Next()
	I.Rewind()
	return t
}

// PeekN returns N tokens without changing position.
func (I *Interpreter) PeekN(n int) []Tokener {
	ret := I.ReadN(n)
	for i := 0; i < n; i++ {
		I.Rewind()
	}
	return ret
}

// Peek returns nil if the next token is not of type T.
func (I *Interpreter) Peek(T ...lexer.TokenType) Tokener {
	t := I.Read(T...)
	if t != nil {
		I.Rewind()
	}
	return t
}

// PeekUntil peeks anything until it met with any of T.
func (I *Interpreter) PeekUntil(T ...lexer.TokenType) Tokener {
	p := I.position
	var ret Tokener
	for {
		n := I.Next()
		if n == nil {
			I.position = p
			return nil
		}
		ret = n
		for _, t := range T {
			if n.GetType() == t {
				I.Rewind()
				return ret
			}
		}
	}
}

// Read advances the position if next token is of type T.
func (I *Interpreter) Read(Ts ...lexer.TokenType) Tokener {
	t := I.Next()
	found := false
	for _, T := range Ts {
		if !found && t != nil && t.GetType() == T {
			found = true
			break
		}
	}
	if !found {
		I.Rewind()
		t = nil
	}
	return t
}

// ReadN advances the position of n.
func (I *Interpreter) ReadN(n int) []Tokener {
	var ret []Tokener
	for i := 0; i < n; i++ {
		ret = append(ret, I.Next())
	}
	return ret
}

// ReadMany advances until next token is not of types T.
func (I *Interpreter) ReadMany(Ts ...lexer.TokenType) []Tokener {
	ret := []Tokener{}
	for {
		found := false
		for _, T := range Ts {
			if !found {
				if t := I.Read(T); t != nil {
					ret = append(ret, t)
					found = true
					break
				}
			}
		}
		if !found {
			break
		}
	}
	return ret
}

// Get reads a token T and flushes the buffer.
// bad idea ?
func (I *Interpreter) Get(T ...lexer.TokenType) Tokener {
	if t := I.Read(T...); t != nil {
		I.Flush()
		return t
	}
	return nil
}

// GetMany reads tokens until it is not of types T and flushes the buffer.
// bad idea ?
func (I *Interpreter) GetMany(T ...lexer.TokenType) []Tokener {
	ret := I.ReadMany(T...)
	I.Flush()
	return ret
}

// ReadBlock reads the tokens as a block delimited by open/close Type ({..}).
func (I *Interpreter) ReadBlock(open lexer.TokenType, close lexer.TokenType) []Tokener {

	var ret []Tokener

	if openTok := I.Read(open); openTok == nil {
		return ret
	}

	count := 1
	for {
		if openTok := I.Read(open); openTok != nil {
			count++

		} else if closeTok := I.Read(close); closeTok != nil {
			count--
			if count == 0 {
				break
			}

		} else {
			I.Next()
		}
	}
	I.Read(close)
	return I.Current()
}

// Debug produces a SyntaxError when unexpected tokens are found.
func (I *Interpreter) Debug(reason string, wantedTypes ...lexer.TokenType) error {
	n := I.Last()
	if n == nil {
		n = I.Next()
		I.Rewind()
	}

	return I.DebugAtToken(n, reason, wantedTypes...)
}

// DebugAtToken produces a SyntaxError at token T.
func (I *Interpreter) DebugAtToken(atToken Tokener, reason string, wantedTypes ...lexer.TokenType) error {
	wanted := []string{}
	for _, w := range wantedTypes {
		wanted = append(wanted, I.Namer(w))
	}

	if atToken == nil {
		// tbd adjust the position
		atToken = NewTokenEOF()
	}
	got := I.Namer(atToken.GetType())
	return I.Scope.FinalizeErr(NewParseError(errors.New(reason), atToken, got, wanted))
}
