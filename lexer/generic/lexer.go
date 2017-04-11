package generic

import (
	"sort"

	lexer "github.com/mh-cbon/state-lexer"
)

// generic tokens
const (
	WsToken lexer.TokenType = iota
	CommentLineToken
	CommentBlockToken
	WordToken
	TextToken
	EOFToken
)

// UnknownTokenLabel is the label for unknow ntokens.
const UnknownTokenLabel = "token unknown"

// TokenName Helper function
func TokenName(tok lexer.Token) string {
	ret := TokenType(tok.Type)
	if ret == UnknownTokenLabel {
		ret = tok.Value + " " + ret
	}
	return ret
}

// TokenType Helper function
func TokenType(Type lexer.TokenType) string {
	switch Type {
	case WsToken:
		return "WsToken"
	case CommentBlockToken:
		return "CommentBlockToken"
	case CommentLineToken:
		return "CommentLineToken"
	case WordToken:
		return "WordToken"
	case TextToken:
		return "TextToken"
	case EOFToken:
		return "EOFToken"
	}
	return UnknownTokenLabel
}

// NotWs Helper function
func NotWs(f func(lexer.Token)) func(lexer.Token) {
	return lexer.Not(WsToken, f)
}

// NotComments Helper function
func NotComments(f func(lexer.Token)) func(lexer.Token) {
	return lexer.Not(CommentLineToken, lexer.Not(CommentBlockToken, f))
}

// Lexer ...
type Lexer struct {
	Words   []Word
	Printer func(Type lexer.TokenType) string
}

// Word ...
type Word struct {
	Value string
	Type  lexer.TokenType
	// Sep           bool
	TextWord bool
	// BeginOnly     bool
	IsBlockIgnore bool
	BlockSepEnd   string
	ExcludeSepEnd bool
	CanEscape     bool
	EscapeStr     string
}

// Words is a sortable list of Word.
type Words []Word

func (slice Words) Len() int {
	return len(slice)
}

func (slice Words) Less(i, j int) bool {
	return len(slice[i].Value) < len(slice[j].Value)
}

func (slice Words) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// GetExactWord return a word whose content==w
func (g *Lexer) GetExactWord(w string) (Word, bool) {
	for _, word := range g.Words {
		if word.Value == w {
			return word, true
		}
	}
	return Word{}, false
}

// StartHere ...
func (g *Lexer) StartHere(l *lexer.L) lexer.StateFunc {

	sort.Sort(sort.Reverse(Words(g.Words)))

	return g.process
}

// StartHere ...
func (g *Lexer) process(l *lexer.L) lexer.StateFunc {
	r := l.Next()
	if r == lexer.EOFRune {
		return nil
	}

	c := l.Current()
	if w, ok := g.GetExactWord(c); ok {
		s := g.getSimilarWords(l, w.Value)
		if len(s) > 0 {
			unreadWord(l, c)
			for _, ws := range s {
				if peekWord(l, ws.Value) == ws.Value {
					if !ws.TextWord || (len(l.Current()) == len(ws.Value) && isNonWord(l.Peek())) {
						unreadWord(l, ws.Value)
						g.Emit(l, WordToken)
						peekWord(l, ws.Value)
						readBlockIgnore(l, ws)
						g.Emit(l, ws.Type)
						return g.process
					}
				}
			}
			peekWord(l, w.Value)
		}

		readBlockIgnore(l, w)

		if !w.TextWord || (len(l.Current()) == len(w.Value) && isNonWord(l.Peek())) {
			g.Emit(l, w.Type)
		}

	} else if s := g.getStartingWords(l, string(r)); len(s) > 0 {
		l.Rewind()
		for _, ws := range s {
			if peekWord(l, ws.Value) == ws.Value {
				if !ws.TextWord || (len(l.Current()) == len(ws.Value) && isNonWord(l.Peek())) {
					unreadWord(l, ws.Value)
					g.Emit(l, WordToken)
					peekWord(l, ws.Value)

					readBlockIgnore(l, ws)
					g.Emit(l, ws.Type)

					return g.process
				} else {
					unreadWord(l, ws.Value)
				}
			}
		}
		l.Next()
	}

	return g.process
}

func readBlockIgnore(l *lexer.L, w Word) {

	if w.IsBlockIgnore {
		if w.CanEscape {
			readBlock(l, w.Value, w.EscapeStr)
		} else {
			readUntil(l, w.BlockSepEnd)
		}
		if w.ExcludeSepEnd {
			unreadWord(l, w.BlockSepEnd)
		}
	}
}

// StartHere ...
func (g *Lexer) getStartingWords(l *lexer.L, w string) []Word {
	ret := []Word{}
	for _, word := range g.Words {
		if len(word.Value) >= len(w) && word.Value[:len(w)] == w {
			ret = append(ret, word)
		}
	}
	return ret
}

// StartHere ...
func (g *Lexer) getSimilarWords(l *lexer.L, w string) []Word {
	ret := []Word{}
	for _, word := range g.Words {
		if len(word.Value) > len(w) && word.Value[:len(w)] == w {
			ret = append(ret, word)
		}
	}
	return ret
}

func isNonWord(ch rune) bool {
	return (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && ch != '_'
}

func rewindAll(l *lexer.L) {
	for {
		l.Rewind()
		if l.Current() == "" {
			break
		}
	}
}

func peekWord(l *lexer.L, w string) string {
	f := ""
	for range w {
		f += string(l.Next())
	}

	if f != w {
		unreadWord(l, w)
	}
	return f
}

func unreadWord(l *lexer.L, w string) {
	for range w {
		l.Rewind()
	}
}

func readBlock(l *lexer.L, blockTerm string, escapeStr string) {
	for {
		readUntil(l, blockTerm)
		unreadWord(l, blockTerm)
		unreadWord(l, escapeStr)
		escaped := peekWord(l, escapeStr) == escapeStr
		if !escaped {
			for range escapeStr {
				l.Next()
			}
			peekWord(l, blockTerm)
			break
		} else {
			peekWord(l, blockTerm)
		}
	}
}

func readUntil(l *lexer.L, w string) {
	for {
		if f := peekWord(l, w); f == w {
			break
		}
		l.Next()
	}
}

// Emit will receive a token type and push a new token with the current analyzed
// value into the tokens channel.
func (g *Lexer) Emit(l *lexer.L, t lexer.TokenType) {
	if l.Current() != "" {
		l.Emit(t)
	}
}

// New gigo lexer
func New() *Lexer {
	return &Lexer{
		Printer: TokenType,
		Words:   []Word{},
	}
}
