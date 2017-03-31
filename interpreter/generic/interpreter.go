package generic

import (
	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
	lexer "github.com/mh-cbon/state-lexer"
)

type Interpreter struct {
	position int
	Tokens   []Tokener
	Scope    ScopeReceiver
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		position: -1,
	}
}

func (I *Interpreter) Next() Tokener {
	if I.position < len(I.Tokens) {
		I.position++
		if I.position < len(I.Tokens) {
			return I.Tokens[I.position]
		}
	}
	return nil
}
func (I *Interpreter) Rewind() {
	I.position--
	if I.position < -1 {
		I.position = -1
	}
}
func (I *Interpreter) Last() Tokener {
	if I.position > -1 && I.position < len(I.Tokens) {
		return I.Tokens[I.position]
	}
	return nil
}

func (I *Interpreter) Ended() bool {
	return len(I.Tokens) == 0
}

func (I *Interpreter) Current() []Tokener {
	if I.position < 0 {
		return I.Tokens[:0]
	}
	if I.position >= len(I.Tokens) {
		return I.Tokens[0:]
	}
	return I.Tokens[:I.position+1]
}

func (I *Interpreter) Emit() []Tokener {
	toks := []Tokener{}
	c := I.Current()
	if len(c) > 0 {
		toks = append(toks, c...)
		I.Flush()
	}
	return toks
}

func (I *Interpreter) Flush() {
	if I.position+1 < len(I.Tokens) {
		I.Tokens = I.Tokens[I.position+1:]
	} else {
		I.Tokens = I.Tokens[:0]
	}
	I.position = -1
}

func (I *Interpreter) PeekOne() Tokener {
	t := I.Next()
	I.Rewind()
	return t
}
func (I *Interpreter) PeekN(n int) []Tokener {
	var ret []Tokener
	for i := 0; i < n; i++ {
		ret = append(ret, I.Next())
	}
	for i := 0; i < n; i++ {
		I.Rewind()
	}
	return ret
}

// func (I *Interpreter) MustRead(T lexer.TokenType) Tokener {
// 	t := I.Read(T)
// 	if t == nil {
// 		panic(I.Debug("Token is unexpected: ", T))
// 	}
// 	return t
// }

func (I *Interpreter) Peek(T lexer.TokenType) Tokener {
	t := I.Next()
	if t == nil {
	} else if t.GetType() != T {
		t = nil
	}
	I.Rewind()
	return t
}

func (I *Interpreter) Read(T lexer.TokenType) Tokener {
	t := I.Next()
	if t == nil || t.GetType() != T {
		I.Rewind()
		t = nil
	}
	return t
}

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

// bad idea ?
func (I *Interpreter) Get(T lexer.TokenType) Tokener {
	if t := I.Read(T); t != nil {
		I.Flush()
		return t
	}
	return nil
}

// bad idea ?
func (I *Interpreter) GetMany(T lexer.TokenType) []Tokener {
	ret := I.ReadMany(T)
	I.Flush()
	return ret
}

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

func (I *Interpreter) Debug(reason string, wantedTypes ...lexer.TokenType) error {
	wanted := []string{}
	for _, w := range wantedTypes {
		wanted = append(wanted, glanglexer.TokenType(w))
	}
	n := I.Last()
	if n == nil {
		n = I.Next()
		I.Rewind()
	} else if n == nil {
		// tbd adjust the position
		n = NewTokenEOF()
	}
	got := glanglexer.TokenType(n.GetType())
	return I.Scope.FinalizeErr(NewParseError(n, reason, got, wanted))
}
