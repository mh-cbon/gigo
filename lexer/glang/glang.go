package glang

import (
	generic "github.com/mh-cbon/gigo/lexer/generic"
	lexer "github.com/mh-cbon/state-lexer"
)

// tokens for a golang source code.
const (
	NumberToken lexer.TokenType = generic.TextToken + iota + 1
	NlToken

	PackageToken

	TypeToken
	StructToken
	ImplementsToken
	TemplateToken
	ImportToken
	InterfaceToken

	ConstToken
	ConstNameToken
	VarToken

	FuncToken

	BracketOpenToken
	BracketCloseToken
	ParenOpenToken
	ParenCloseToken

	AssignToken
	TypeAssignToken
	SmallerToken
	GreaterToken
	EqToken
	NegateToken
	NeqToken
	GteqToken
	SmeqToken

	ColonToken
	SemiColonToken
	CommaToken

	ReturnToken
	ForToken
	RangeToken
	DeferToken
	ElseToken
	IfToken
	ElipseToken
	TplOpenToken
	TplCloseToken

	PoireauToken
	PoireauPointerToken

	BraceOpenToken
	BraceCloseToken

	TrueToken
	FalseToken
	//-

	AddToken
	SubToken
	MulToken
	QuoToken
	RemToken
	AndToken
	OrToken
	XorToken
	ShlToken
	ShrToken
	AndNotToken

	AddAssignToken
	SubAssignToken
	MulAssignToken
	QuoAssignToken
	RemAssignToken

	AndAssignToken
	OrAssignToken
	XorAssignToken
	ShlAssignToken
	ShrAssignToken
	AndNotAssignToken

	LAndToken
	LOrToken
	ArrowToken
	IncToken
	DecToken
	DotToken

	ChanToken
	BreakToken
	ContinueToken
	GoToken
	GotoToken
	MapToken
	FallthroughToken
	DefaultToken

	StringToken
	IntToken
	Int8Token
	Int16Token
	Int32Token
	Int64Token
	UintToken
	Uint8Token
	Uint16Token
	Uint32Token
	Uint64Token
	FloatToken
	Float32Token
	Float64Token
)

// TokenName Helper function
func TokenName(tok lexer.Token) string {
	return TokenType(tok.Type)
}

// TokenType Helper function
func TokenType(Type lexer.TokenType) string {
	ret := generic.TokenType(Type)
	if ret != "token unknown" {
		return ret
	}
	switch Type {
	case TrueToken:
		return "TrueToken"
	case FalseToken:
		return "FalseToken"
	case NumberToken:
		return "numberToken"
	case VarToken:
		return "varToken"
	case ConstNameToken:
		return "ConstNameToken"
	case ConstToken:
		return "constToken"
	case ColonToken:
		return "ColonToken"
	case SemiColonToken:
		return "semiColonToken"
	case CommaToken:
		return "CommaToken"
	case ReturnToken:
		return "returnToken"
	case ForToken:
		return "forToken"
	case RangeToken:
		return "rangeToken"
	case DeferToken:
		return "deferToken"
	case FuncToken:
		return "funcToken"
	case ElseToken:
		return "elseToken"
	case IfToken:
		return "ifToken"
	case ElipseToken:
		return "ElipseToken"
	case PackageToken:
		return "packageToken"
	case TypeToken:
		return "typeToken"
	case StructToken:
		return "structToken"
	case ImplementsToken:
		return "implementsToken"
	case NlToken:
		return "nlToken"
	case AssignToken:
		return "assignToken"
	case TypeAssignToken:
		return "TypeAssignToken"
	case BraceOpenToken:
		return "BraceOpenToken"
	case BraceCloseToken:
		return "BraceCloseToken"
	case BracketOpenToken:
		return "bracketOpenToken"
	case BracketCloseToken:
		return "bracketCloseToken"
	case ParenOpenToken:
		return "parenOpenToken"
	case ParenCloseToken:
		return "parenCloseToken"
	case SmallerToken:
		return "smallerToken"
	case GreaterToken:
		return "greaterToken"
	case EqToken:
		return "eqToken"
	case NegateToken:
		return "NegateToken"
	case NeqToken:
		return "NeqToken"
	case GteqToken:
		return "GteqToken"
	case SmeqToken:
		return "SmeqToken"
	case TemplateToken:
		return "templateToken"
	case ImportToken:
		return "ImportToken"
	case InterfaceToken:
		return "interfaceToken"
	case TplOpenToken:
		return "TplOpenToken"
	case TplCloseToken:
		return "TplCloseToken"
	case PoireauToken:
		return "PoireauToken"
	case PoireauPointerToken:
		return "PoireauPointerToken"
	case AddToken:
		return "AddToken"
	case SubToken:
		return "SubToken"
	case MulToken:
		return "MulToken"
	case QuoToken:
		return "QuoToken"
	case RemToken:
		return "RemToken"

	case AndToken:
		return "AndToken"
	case OrToken:
		return "OrToken"
	case XorToken:
		return "XorToken"
	case ShlToken:
		return "ShlToken"
	case ShrToken:
		return "ShrToken"
	case AndNotToken:
		return "AndNotToken"

	case AddAssignToken:
		return "AddAssignToken"
	case SubAssignToken:
		return "SubAssignToken"
	case MulAssignToken:
		return "MulAssignToken"
	case QuoAssignToken:
		return "QuoAssignToken"
	case RemAssignToken:
		return "RemAssignToken"

	case AndAssignToken:
		return "AndAssignToken"
	case OrAssignToken:
		return "OrAssignToken"
	case XorAssignToken:
		return "XorAssignToken"
	case ShlAssignToken:
		return "ShlAssignToken"
	case ShrAssignToken:
		return "ShrAssignToken"
	case AndNotAssignToken:
		return "AndNotAssignToken"

	case LAndToken:
		return "LAndToken"
	case LOrToken:
		return "LOrToken"
	case ArrowToken:
		return "ArrowToken"
	case IncToken:
		return "IncToken"
	case DecToken:
		return "DecToken"
	case DotToken:
		return "DotToken"
	case ChanToken:
		return "ChanToken"
	case BreakToken:
		return "BreakToken"
	case ContinueToken:
		return "ContinueToken"
	case GoToken:
		return "GoToken"
	case GotoToken:
		return "GotoToken"
	case MapToken:
		return "MapToken"
	case FallthroughToken:
		return "FallthroughToken"
	case DefaultToken:
		return "DefaultToken"

	// type keywors
	case StringToken:
		return "StringToken"
	case IntToken:
		return "IntToken"
	case Int8Token:
		return "Int8Token"
	case Int16Token:
		return "Int16Token"
	case Int32Token:
		return "Int32Token"
	case Int64Token:
		return "Int64Token"
	case UintToken:
		return "UintToken"
	case Uint8Token:
		return "Uint8Token"
	case Uint16Token:
		return "Uint16Token"
	case Uint32Token:
		return "Uint32Token"
	case Uint64Token:
		return "Uint64Token"
	case FloatToken:
		return "FloatToken"
	case Float32Token:
		return "Float32Token"
	case Float64Token:
		return "Float64Token"
	}
	return "token unknown"
}

// New ...
func New() *generic.Lexer {
	return &generic.Lexer{
		Printer: TokenType,
		Words: []generic.Word{
			// comments
			generic.Word{Value: "//", Type: generic.CommentLineToken, IsBlockIgnore: true, BlockSepEnd: "\n", ExcludeSepEnd: true},
			generic.Word{Value: "/*", Type: generic.CommentBlockToken, IsBlockIgnore: true, BlockSepEnd: "*/"},
			// Texts
			generic.Word{Value: "\"", Type: generic.TextToken, IsBlockIgnore: true, BlockSepEnd: "\"", CanEscape: true, EscapeStr: "\\"},
			generic.Word{Value: "`", Type: generic.TextToken, IsBlockIgnore: true, BlockSepEnd: "`"},
			// ws
			generic.Word{Value: " ", Type: generic.WsToken},
			generic.Word{Value: "\t", Type: generic.WsToken},
			generic.Word{Value: "\n", Type: NlToken},

			// Operators and delimiters
			generic.Word{Value: "+", Type: AddToken},
			generic.Word{Value: "-", Type: SubToken},
			generic.Word{Value: "*", Type: MulToken},
			generic.Word{Value: "/", Type: QuoToken},
			generic.Word{Value: "%", Type: RemToken},

			generic.Word{Value: "&", Type: AndToken},
			generic.Word{Value: "|", Type: OrToken},
			generic.Word{Value: "^", Type: XorToken},
			generic.Word{Value: "<<", Type: ShlToken},
			generic.Word{Value: ">>", Type: ShrToken},
			generic.Word{Value: "&^", Type: AndNotToken},

			generic.Word{Value: "+=", Type: AddAssignToken},
			generic.Word{Value: "-=", Type: SubAssignToken},
			generic.Word{Value: "*=", Type: MulAssignToken},
			generic.Word{Value: "/=", Type: QuoAssignToken},
			generic.Word{Value: "%=", Type: RemAssignToken},

			generic.Word{Value: "&=", Type: AndAssignToken},
			generic.Word{Value: "|=", Type: OrAssignToken},
			generic.Word{Value: "^=", Type: XorAssignToken},
			generic.Word{Value: "<<=", Type: ShlAssignToken},
			generic.Word{Value: ">>=", Type: ShrAssignToken},
			generic.Word{Value: "&^=", Type: AndNotAssignToken},

			generic.Word{Value: "&&", Type: LAndToken},
			generic.Word{Value: "||", Type: LOrToken},
			generic.Word{Value: "<-", Type: ArrowToken},
			generic.Word{Value: "++", Type: IncToken},
			generic.Word{Value: "--", Type: DecToken},

			generic.Word{Value: "==", Type: EqToken},
			generic.Word{Value: "<", Type: SmallerToken},
			generic.Word{Value: ">", Type: GreaterToken},
			generic.Word{Value: "=", Type: AssignToken},
			generic.Word{Value: "!", Type: NegateToken},

			generic.Word{Value: "!=", Type: NeqToken},
			generic.Word{Value: "<=", Type: SmeqToken},
			generic.Word{Value: ">=", Type: GteqToken},
			generic.Word{Value: ":=", Type: TypeAssignToken},
			generic.Word{Value: "...", Type: ElipseToken},

			generic.Word{Value: "(", Type: ParenOpenToken},
			generic.Word{Value: "[", Type: BracketOpenToken},
			generic.Word{Value: "{", Type: BraceOpenToken},
			generic.Word{Value: ",", Type: CommaToken},
			generic.Word{Value: ".", Type: DotToken},
			generic.Word{Value: "}", Type: BraceCloseToken},
			generic.Word{Value: "]", Type: BracketCloseToken},
			generic.Word{Value: ")", Type: ParenCloseToken},
			generic.Word{Value: ";", Type: SemiColonToken},
			generic.Word{Value: ":", Type: ColonToken},

			generic.Word{Value: "<:", Type: TplOpenToken},

			// generic.Word{Value: "<", Type: TplOpenToken},
			// generic.Word{Value: ">", Type: TplCloseToken},
			// type keywords
			generic.Word{Value: "string", Type: StringToken, TextWord: true},
			generic.Word{Value: "int", Type: IntToken, TextWord: true},
			generic.Word{Value: "int8", Type: Int8Token, TextWord: true},
			generic.Word{Value: "int16", Type: Int16Token, TextWord: true},
			generic.Word{Value: "int32", Type: Int32Token, TextWord: true},
			generic.Word{Value: "int64", Type: Int64Token, TextWord: true},
			generic.Word{Value: "uint", Type: UintToken, TextWord: true},
			generic.Word{Value: "uint8", Type: Uint8Token, TextWord: true},
			generic.Word{Value: "uint16", Type: Uint16Token, TextWord: true},
			generic.Word{Value: "uint32", Type: Uint32Token, TextWord: true},
			generic.Word{Value: "uint64", Type: Uint64Token, TextWord: true},
			generic.Word{Value: "float", Type: FloatToken, TextWord: true},
			generic.Word{Value: "float32", Type: Float32Token, TextWord: true},
			generic.Word{Value: "float64", Type: Float64Token, TextWord: true},

			// ... keywords
			generic.Word{Value: "chan", Type: ChanToken, TextWord: true},
			generic.Word{Value: "break", Type: BreakToken, TextWord: true},
			generic.Word{Value: "continue", Type: ContinueToken, TextWord: true},
			generic.Word{Value: "go", Type: GoToken, TextWord: true},
			generic.Word{Value: "goto", Type: GotoToken, TextWord: true},
			generic.Word{Value: "map", Type: MapToken, TextWord: true},
			generic.Word{Value: "fallthrough", Type: FallthroughToken, TextWord: true},
			generic.Word{Value: "default", Type: DefaultToken, TextWord: true},
			generic.Word{Value: "if", Type: IfToken, TextWord: true},
			generic.Word{Value: "true", Type: TrueToken, TextWord: true}, // the real one, neo.
			generic.Word{Value: "false", Type: FalseToken, TextWord: true},
			generic.Word{Value: "for", Type: ForToken, TextWord: true},
			generic.Word{Value: "var", Type: VarToken, TextWord: true},
			generic.Word{Value: "const", Type: ConstToken, TextWord: true},
			generic.Word{Value: "func", Type: FuncToken, TextWord: true},
			generic.Word{Value: "type", Type: TypeToken, TextWord: true},
			generic.Word{Value: "else", Type: ElseToken, TextWord: true},
			generic.Word{Value: "range", Type: RangeToken, TextWord: true},
			generic.Word{Value: "defer", Type: DeferToken, TextWord: true},
			generic.Word{Value: "struct", Type: StructToken, TextWord: true},
			generic.Word{Value: "return", Type: ReturnToken, TextWord: true},
			generic.Word{Value: "package", Type: PackageToken, TextWord: true},
			generic.Word{Value: "import", Type: ImportToken, TextWord: true},
			generic.Word{Value: "template", Type: TemplateToken, TextWord: true},
			generic.Word{Value: "constname", Type: ConstNameToken, TextWord: true},
			generic.Word{Value: "interface", Type: InterfaceToken, TextWord: true},
			generic.Word{Value: "implements", Type: ImplementsToken, TextWord: true},
			generic.Word{Value: "poireau", Type: PoireauToken, TextWord: true},
			generic.Word{Value: "*poireau", Type: PoireauPointerToken, TextWord: true},
		},
	}
}
