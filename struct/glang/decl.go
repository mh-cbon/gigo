package glang

import (
	"regexp"
	"strings"

	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
)

// ScopeDecl defines the source code origin, a file o string.
type ScopeDecl struct {
	genericinterperter.Expression
}

// GrepLine finds all token on given line.
func (f *ScopeDecl) GrepLine(line int) []genericinterperter.Tokener {
	return f.GetTokensAtLine(line)
}

// FindPackagesDecl returns all package declarations found.
func (f *ScopeDecl) FindPackagesDecl() []*PackageDecl {
	var ret []*PackageDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*PackageDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindImplementsTypes returns all implements declarations found.
func (f *ScopeDecl) FindImplementsTypes() []*ImplementDecl {
	var ret []*ImplementDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*ImplementDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindStructsTypes returns all struct declarations found.
func (f *ScopeDecl) FindStructsTypes() []*StructDecl {
	var ret []*StructDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*StructDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindTemplatesTypes returns all template declarations found.
func (f *ScopeDecl) FindTemplatesTypes() []*TemplateDecl {
	var ret []*TemplateDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*TemplateDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindInterfaces returns all interface type declarations.
func (f *ScopeDecl) FindInterfaces() []*InterfaceDecl {
	var ret []*InterfaceDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*InterfaceDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindFuncs returns all func declarations.
func (f *ScopeDecl) FindFuncs() []*FuncDecl {
	var ret []*FuncDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*FuncDecl); ok && x.IsTemplated() == false {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindTemplateFuncs returns all funcs with templating declarations.
func (f *ScopeDecl) FindTemplateFuncs() []FuncDeclarer {
	var ret []FuncDeclarer
	for _, t := range f.Tokens {
		if x, ok := t.(*TemplateFuncDecl); ok && x.IsDefine() == false {
			ret = append(ret, x)
		}
	}
	for _, t := range f.Tokens {
		if x, ok := t.(*FuncDecl); ok && x.IsTemplated() {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindDefineFuncs returns all <define> declarations.
func (f *ScopeDecl) FindDefineFuncs() []*TemplateFuncDecl {
	var ret []*TemplateFuncDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*TemplateFuncDecl); ok && x.IsDefine() {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindVarDecl returns all var declarations.
func (f *ScopeDecl) FindVarDecl() []*VarDecl {
	var ret []*VarDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*VarDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}

// FindConstDecl returns all const declarations.
func (f *ScopeDecl) FindConstDecl() []*ConstDecl {
	var ret []*ConstDecl
	for _, t := range f.Tokens {
		if x, ok := t.(*ConstDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}

type slugamer interface {
	GetSlugName() string
}
type namer interface {
	GetName() string
}

// FindSymbols returns declarations that matches given symbol name.
func (f *ScopeDecl) FindSymbols(symbol string) []genericinterperter.Expressioner {
	ret := []genericinterperter.Expressioner{}
	for _, t := range f.GetExprs() {
		if x, ok := t.(slugamer); ok {
			if strings.TrimSpace(x.GetSlugName()) == symbol { // should not need to trim here.
				ret = append(ret, t)
			}
		} else if x, ok := t.(namer); ok {
			if strings.TrimSpace(x.GetName()) == symbol { // should not need to trim here.
				ret = append(ret, t)
			}
		}
	}
	return ret
}

// StrDecl is a string source code.
type StrDecl struct {
	ScopeDecl
	Src string
}

// FinalizeErr contextualize an error for pretty printing.
func (f *StrDecl) FinalizeErr(err *genericinterperter.ParseError) error {
	return &genericinterperter.StringSyntaxError{Src: f.Src, Filepath: "<noname>", ParseError: *err}
}

// GetName implements ScopeReceiver
func (f *StrDecl) GetName() string {
	return "noname"
}

// FileDecl is a source code from a file.
type FileDecl struct {
	ScopeDecl
	Name string
}

// FinalizeErr contextualize an error for pretty printing.
func (f *FileDecl) FinalizeErr(err *genericinterperter.ParseError) error {
	return &genericinterperter.FileSyntaxError{Src: f.GetName(), ParseError: *err}
}

// GetName implements ScopeReceiver
func (f *FileDecl) GetName() string {
	return f.Name
}

// PackageDecl for package <name> declaratons
type PackageDecl struct {
	genericinterperter.Expression
	Name genericinterperter.Tokener
}

func (p *PackageDecl) String() string {
	return p.Expression.String()
}

// GetName returns the name of the package.
func (p *PackageDecl) GetName() string {
	return p.Name.GetValue()
}

// NewPackageDecl creates a new PackageDecl
func NewPackageDecl() *PackageDecl {
	return &PackageDecl{}
}

type BodyBlockDecl struct {
	genericinterperter.Expression
	Open  genericinterperter.Tokener
	Close genericinterperter.Tokener
}

func (p *BodyBlockDecl) String() string {
	return p.Expression.String()
}
func (p *BodyBlockDecl) GetOpen() string {
	return p.Open.GetValue()
}
func (p *BodyBlockDecl) GetClose() string {
	return p.Close.GetValue()
}

// NewBodyBlockDecl creates a new BodyBlockDecl
func NewBodyBlockDecl() *BodyBlockDecl {
	return &BodyBlockDecl{}
}

type StructDecl struct {
	genericinterperter.Expression
	Name    *IdentifierDecl
	Methods []FuncDeclarer
	Block   *PropsBlockDecl
}

func (p *StructDecl) String() string {
	return p.Expression.String()
}
func (p *StructDecl) GetName() string {
	return p.Name.String()
}
func (p *StructDecl) AddMethod(f FuncDeclarer) {
	p.Methods = append(p.Methods, f)
}

// NewStructDecl creates a new StructDecl
func NewStructDecl() *StructDecl {
	return &StructDecl{}
}

type TemplateDecl struct {
	genericinterperter.Expression
	Name    *IdentifierDecl
	Methods []FuncDeclarer
	Block   *PropsBlockDecl
}

func (t *TemplateDecl) SetDelims(l, r string) {
	t.SetTokenValue(glanglexer.TplOpenToken, l)
	t.SetTokenValue(glanglexer.TplCloseToken, r)
}

func (t *TemplateDecl) GetSlugName() string {
	return t.Name.GetSlugName()
}
func (t *TemplateDecl) Template() string {
	ret := t.String()
	for _, m := range t.Methods {
		ret += m.String()
	}
	return ret
}
func (t *TemplateDecl) String() string {
	return t.Expression.String()
}
func (t *TemplateDecl) GetName() string {
	return t.Name.String()
}
func (t *TemplateDecl) AddMethod(f FuncDeclarer) {
	t.Methods = append(t.Methods, f)
}

// NewTemplateDecl creates a new TemplateDecl
func NewTemplateDecl() *TemplateDecl {
	return &TemplateDecl{}
}

type InterfaceDecl struct {
	genericinterperter.Expression
	Name  *IdentifierDecl
	Block *SignsBlockDecl
}

func (p *InterfaceDecl) String() string {
	return p.Expression.String()
}
func (p *InterfaceDecl) GetName() string {
	return p.Name.GetValue()
}

// NewInterfaceDecl creates a new InterfaceDecl
func NewInterfaceDecl() *InterfaceDecl {
	return &InterfaceDecl{}
}

type ImplementDecl struct {
	genericinterperter.Expression
	Name              *IdentifierDecl
	ImplementTemplate genericinterperter.Tokener
	Methods           []FuncDeclarer
}

func (p *ImplementDecl) AddMethod(f FuncDeclarer) {
	p.Methods = append(p.Methods, f)
}
func (p *ImplementDecl) SetImplKeyword(s string) {
	p.SetTokenValue(glanglexer.ImplementsToken, s)
}
func (p *ImplementDecl) String() string {
	return p.Expression.String()
}
func (p *ImplementDecl) GetName() string {
	return p.Name.GetValue()
}
func (p *ImplementDecl) GetImplementTemplate() string {
	return p.ImplementTemplate.String()
}
func (p *ImplementDecl) GetBlock() *PropsBlockDecl {
	return p.FilterToken(glanglexer.BracketOpenToken).(*PropsBlockDecl)
}

// NewImplementDecl creates a new ImplementDecl
func NewImplementDecl() *ImplementDecl {
	return &ImplementDecl{}
}

type PoireauDecl struct {
	genericinterperter.Expression
	ImplementTemplate *IdentifierDecl
}

func (p *PoireauDecl) String() string {
	return p.Expression.String()
}
func (p *PoireauDecl) IsPointer() bool {
	return p.GetType() == glanglexer.PoireauPointerToken
}
func (p *PoireauDecl) GetImplementTemplate() string {
	return p.ImplementTemplate.String()
}

// NewPoireauDecl creates a new PoireauDecl
func NewPoireauDecl() *PoireauDecl {
	return &PoireauDecl{}
}

type TemplateFuncDecl struct {
	genericinterperter.Expression
	Func     *FuncDecl
	Modifier *BodyBlockDecl
}

// NewTemplateFuncDecl creates a new TemplateFuncDecl
func NewTemplateFuncDecl() *TemplateFuncDecl {
	return &TemplateFuncDecl{}
}

func (t *TemplateFuncDecl) SetDelims(l, r string) {
	t.SetTokenValue(glanglexer.TplOpenToken, l)
	t.SetTokenValue(glanglexer.TplCloseToken, r)
}
func (t *TemplateFuncDecl) GetReceiverType() *IdentifierDecl {
	return t.Func.GetReceiverType()
}
func (t *TemplateFuncDecl) GetName() string {
	return t.Func.GetName()
}
func (t *TemplateFuncDecl) IsMethod() bool {
	return t.Func.IsMethod()
}
func (t *TemplateFuncDecl) IsDefine() bool {
	return t.Modifier.String() == "<define>"
}
func (p *TemplateFuncDecl) String() string {
	return p.Expression.String()
}
func (t *TemplateFuncDecl) IsTemplated() bool {
	return true
}
func (t *TemplateFuncDecl) GetModifier() *BodyBlockDecl {
	return t.Modifier
}
func (t *TemplateFuncDecl) GetReceiver() *PropsBlockDecl {
	return t.Func.GetReceiver()
}
func (t *TemplateFuncDecl) GetBody() *BodyBlockDecl {
	return t.Func.GetBody()
}
func (t *TemplateFuncDecl) GetArgs() *PropsBlockDecl {
	return t.Func.GetArgs()
}
func (t *TemplateFuncDecl) GetArgsBlock() []*PropDecl {
	return t.Func.GetArgsBlock()
}
func (t *TemplateFuncDecl) GetArgsNames() []*IdentifierDecl {
	return t.Func.GetArgsNames()
}

// FuncDeclarer is a func or a template func
type FuncDeclarer interface {
	genericinterperter.Expressioner
	IsMethod() bool
	IsTemplated() bool
	GetName() string
	// GetSlugName() string
	GetReceiverType() *IdentifierDecl
	GetReceiver() *PropsBlockDecl
	GetModifier() *BodyBlockDecl
	GetBody() *BodyBlockDecl
	GetArgs() *PropsBlockDecl
	GetPos() genericinterperter.TokenPos
	GetArgsBlock() []*PropDecl
	GetArgsNames() []*IdentifierDecl
}

type FuncDecl struct {
	genericinterperter.Expression
	Receiver *PropsBlockDecl
	Name     *IdentifierDecl
	Params   *PropsBlockDecl
	Out      *PropsBlockDecl
	Body     *BodyBlockDecl
}

func (p *FuncDecl) GetArgs() *PropsBlockDecl {
	return p.Params
}
func (p *FuncDecl) GetArgsBlock() []*PropDecl {
	return p.Params.Props
}
func (p *FuncDecl) GetArgsNames() []*IdentifierDecl {
	ret := []*IdentifierDecl{}
	for _, p := range p.Params.Props {
		ret = append(ret, p.Name)
	}
	return ret
}
func (p *FuncDecl) GetBody() *BodyBlockDecl {
	return p.Body
}
func (p *FuncDecl) GetReceiver() *PropsBlockDecl {
	return p.Receiver
}
func (p *FuncDecl) GetReceiverType() *IdentifierDecl {
	return p.Receiver.Props[0].Type
}
func (p *FuncDecl) GetName() string {
	return p.Name.String()
}
func (p *FuncDecl) IsTemplated() bool {
	return p.HasToken(glanglexer.TplOpenToken) && p.HasToken(glanglexer.TplCloseToken)
}
func (p *FuncDecl) IsMethod() bool {
	return p.Receiver != nil && len(p.Receiver.Props) > 0
}
func (p *FuncDecl) String() string {
	return p.Expression.String()
}
func (p *FuncDecl) GetModifier() *BodyBlockDecl {
	return nil
}

// NewFuncDecl creates a new FuncDecl
func NewFuncDecl() *FuncDecl {
	return &FuncDecl{}
}

type SignsBlockDecl struct {
	genericinterperter.Expression
	Underlying []*IdentifierDecl
	Signs      []*FuncDecl
}

func (p *SignsBlockDecl) String() string {
	return p.Expression.String()
}
func (p *SignsBlockDecl) AddUnderlying(Type *IdentifierDecl) {
	p.Underlying = append(p.Underlying, Type)
	p.Expression.AddExpr(Type)
}
func (p *SignsBlockDecl) Add(sign *FuncDecl) {
	p.Signs = append(p.Signs, sign)
	p.Expression.AddExpr(sign)
}

// NewSignsBlockDecl creates a new SignsBlockDecl
func NewSignsBlockDecl() *SignsBlockDecl {
	return &SignsBlockDecl{}
}

type PropsBlockDecl struct {
	genericinterperter.Expression
	Poireaux   []*PoireauDecl
	Underlying []*IdentifierDecl
	Props      []*PropDecl
}

func (p *PropsBlockDecl) String() string {
	return p.Expression.String()
}
func (p *PropsBlockDecl) AddUnderlying(Type *IdentifierDecl) {
	p.Underlying = append(p.Underlying, Type)
	p.Expression.AddExpr(Type)
}
func (p *PropsBlockDecl) AddPoireau(Mutation *PoireauDecl) {
	p.Poireaux = append(p.Poireaux, Mutation)
	p.Expression.AddExpr(Mutation)
}
func (p *PropsBlockDecl) Add(Name *IdentifierDecl, Type *IdentifierDecl) *PropDecl {
	prop := NewPropDecl()
	prop.Name = Name
	prop.Type = Type
	prop.AddExpr(Name)
	prop.AddExpr(Type)
	p.Props = append(p.Props, prop)
	p.Expression.AddExpr(prop)
	return prop
}
func (p *PropsBlockDecl) AddT(Type *IdentifierDecl) *PropDecl {
	prop := NewPropDecl()
	prop.Type = Type
	prop.AddExpr(Type)
	p.Props = append(p.Props, prop)
	p.Expression.AddExpr(prop)
	return prop
}

// NewPropsBlockDecl creates a new PropsBlockDecl
func NewPropsBlockDecl() *PropsBlockDecl {
	return &PropsBlockDecl{}
}

type AssignsBlockDecl struct {
	genericinterperter.Expression
	Assigns []*AssignDecl
}

func (p *AssignsBlockDecl) String() string {
	return p.Expression.String()
}
func (p *AssignsBlockDecl) Add(left, leftType, right genericinterperter.Tokener) *AssignDecl {
	as := NewAssignDecl()
	as.Left = left
	as.LeftType = leftType
	as.Right = right
	p.Assigns = append(p.Assigns, as)
	p.Expression.AddExpr(as)
	return as
}
func (p *AssignsBlockDecl) GetAssignments() []*AssignDecl {
	return p.Assigns
}

// NewAssignsBlockDecl creates a new AssignsBlockDecl
func NewAssignsBlockDecl() *AssignsBlockDecl {
	return &AssignsBlockDecl{}
}

type AssignDecl struct {
	genericinterperter.Expression
	Left     genericinterperter.Tokener
	LeftType genericinterperter.Tokener
	Assign   genericinterperter.Tokener
	Right    genericinterperter.Tokener
}

func (p *AssignDecl) String() string {
	return p.Expression.String()
}
func (p *AssignDecl) GetLeft() string {
	return p.Left.GetValue()
}
func (p *AssignDecl) GetLeftType() string {
	return p.LeftType.GetValue()
}
func (p *AssignDecl) GetAssign() string {
	return p.Assign.GetValue()
}
func (p *AssignDecl) GetRight() string {
	return p.Right.GetValue()
}
func (p *AssignDecl) GetAssignments() []*AssignDecl {
	return []*AssignDecl{p}
}

// NewAssignDecl creates a new AssignDecl
func NewAssignDecl() *AssignDecl {
	return &AssignDecl{}
}

type PropDecl struct {
	genericinterperter.Expression
	Name *IdentifierDecl
	Type *IdentifierDecl
}

func (p *PropDecl) String() string {
	return p.Expression.String()
}
func (p *PropDecl) GetName() string {
	return p.Name.GetValue()
}
func (p *PropDecl) GetPropType() string {
	return p.Type.GetValue()
}

// NewPropDecl creates a new PropDecl
func NewPropDecl() *PropDecl {
	return &PropDecl{}
}

type IdentifierDecl struct {
	genericinterperter.Expression
}

var re = regexp.MustCompile("(<[^>]+>)")

func (p *IdentifierDecl) IsImplement() bool {
	return p.GetType() == glanglexer.PoireauToken || p.GetType() == glanglexer.PoireauPointerToken
}

func (p *IdentifierDecl) GetSlugName() string {
	s := p.Expression.String()
	s = re.ReplaceAllString(s, "")
	s = strings.Replace(s, "*", "", -1)
	return strings.TrimSpace(s)
}
func (p *IdentifierDecl) String() string {
	return p.Expression.String()
}

// NewIdentifierDecl creates a new IdentifierDecl
func NewIdentifierDecl() *IdentifierDecl {
	return &IdentifierDecl{}
}

type AssignDeclarer interface {
	GetAssignments() []*AssignDecl
}

type VarDecl struct {
	genericinterperter.Expression
	Assignments []AssignDeclarer
}

func (p *VarDecl) GetAssignments() []*AssignDecl {
	var ret []*AssignDecl
	for _, z := range p.Assignments {
		ret = append(ret, z.GetAssignments()...)
	}
	return ret
}

func (p *VarDecl) AddAssignment(a AssignDeclarer) {
	p.Assignments = append(p.Assignments, a)
}
func (p *VarDecl) String() string {
	return p.Expression.String()
}

// NewVarDecl creates a new VarDecl
func NewVarDecl() *VarDecl {
	return &VarDecl{}
}

type ConstDecl struct {
	genericinterperter.Expression
	Assignments []AssignDeclarer
}

func (p *ConstDecl) GetAssignments() []*AssignDecl {
	var ret []*AssignDecl
	for _, z := range p.Assignments {
		ret = append(ret, z.GetAssignments()...)
	}
	return ret
}
func (p *ConstDecl) AddAssignment(a AssignDeclarer) {
	p.Assignments = append(p.Assignments, a)
}
func (p *ConstDecl) String() string {
	return p.Expression.String()
}

// NewConstDecl creates a new ConstDecl
func NewConstDecl() *ConstDecl {
	return &ConstDecl{}
}

type ExpressionDecl struct {
	genericinterperter.Expression
}

func (p *ExpressionDecl) String() string {
	return p.Expression.String()
}

// NewExpressionDecl creates a new ExpressionDecl
func NewExpressionDecl() *ExpressionDecl {
	return &ExpressionDecl{}
}

// type CommentGroupDecl struct {
// 	// genericinterperter.Tokener
// 	genericinterperter.Expression
// }
//
// func (p *CommentGroupDecl) String() string {
// 	return p.Expression.String()
// }
// func (p *CommentGroupDecl) HasAny() bool {
// 	return len(p.Expression.Tokens) > 0
// }
//
// func NewCommentGroupDecl() *CommentGroupDecl {
// 	c := &CommentGroupDecl{}
// 	return c
// }
