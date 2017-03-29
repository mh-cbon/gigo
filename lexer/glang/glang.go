package glang

import (
	generic "github.com/mh-cbon/gigo/lexer/generic"
	lexer "github.com/mh-cbon/state-lexer"
)

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

	SemiToken
	SemiColonToken

	ReturnToken
	ForToken
	RangeToken
	DeferToken
	ElseToken
	IfToken
	ElipseToken
	TplOpenToken
	TplCloseToken
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
	case NumberToken:
		return "numberToken"
	case VarToken:
		return "varToken"
	case ConstNameToken:
		return "ConstNameToken"
	case ConstToken:
		return "constToken"
	case SemiToken:
		return "semiToken"
	case SemiColonToken:
		return "semiColonToken"
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
	}
	return "token unknown"
}

// New ...
func New() *generic.Lexer {
	return &generic.Lexer{
		Printer: TokenType,
		Words: []generic.Word{
			generic.Word{Value: "//", Type: generic.CommentLineToken, IsBlockIgnore: true, BlockSepEnd: "\n", ExcludeSepEnd: true},
			generic.Word{Value: "/*", Type: generic.CommentBlockToken, IsBlockIgnore: true, BlockSepEnd: "*/"},
			generic.Word{Value: "\"", Type: generic.TextToken, IsBlockIgnore: true, BlockSepEnd: "\"", CanEscape: true, EscapeStr: "\\"},
			generic.Word{Value: "```", Type: generic.TextToken, IsBlockIgnore: true, BlockSepEnd: "```"},
			generic.Word{Value: " ", Type: generic.WsToken},
			generic.Word{Value: "\t", Type: generic.WsToken},
			generic.Word{Value: "\n", Type: NlToken},
			generic.Word{Value: ";", Type: SemiToken},
			generic.Word{Value: ",", Type: SemiColonToken},
			generic.Word{Value: "{", Type: BracketOpenToken},
			generic.Word{Value: "}", Type: BracketCloseToken},
			generic.Word{Value: "(", Type: ParenOpenToken},
			generic.Word{Value: ")", Type: ParenCloseToken},
			// generic.Word{Value: "<", Type: TplOpenToken},
			// generic.Word{Value: ">", Type: TplCloseToken},
			generic.Word{Value: "<", Type: SmallerToken},
			generic.Word{Value: ">", Type: GreaterToken},
			generic.Word{Value: "=", Type: AssignToken},
			generic.Word{Value: ":=", Type: TypeAssignToken},
			generic.Word{Value: "==", Type: EqToken},
			generic.Word{Value: "!", Type: NegateToken},
			generic.Word{Value: "!=", Type: NeqToken},
			generic.Word{Value: ">=", Type: GteqToken},
			generic.Word{Value: "<=", Type: SmeqToken},
			generic.Word{Value: "if", Type: IfToken},
			generic.Word{Value: "...", Type: ElipseToken},
			generic.Word{Value: "for", Type: ForToken},
			generic.Word{Value: "var", Type: VarToken},
			generic.Word{Value: "const", Type: ConstToken},
			generic.Word{Value: "func", Type: FuncToken},
			generic.Word{Value: "type", Type: TypeToken},
			generic.Word{Value: "else", Type: ElseToken},
			generic.Word{Value: "range", Type: RangeToken},
			generic.Word{Value: "defer", Type: DeferToken},
			generic.Word{Value: "struct", Type: StructToken},
			generic.Word{Value: "return", Type: ReturnToken},
			generic.Word{Value: "package", Type: PackageToken},
			generic.Word{Value: "import", Type: ImportToken},
			generic.Word{Value: "template", Type: TemplateToken},
			generic.Word{Value: "constname", Type: ConstNameToken},
			generic.Word{Value: "interface", Type: InterfaceToken},
			generic.Word{Value: "implements", Type: ImplementsToken},
		},
	}
}
