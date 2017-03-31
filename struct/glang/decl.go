package glang

import (
	"regexp"
	"strings"

	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
)

type ScopeDecl struct {
	genericinterperter.Expression
}

func (f *ScopeDecl) GrepLine(line int) []genericinterperter.Tokener {
	return f.GetTokensAtLine(line)
}
func (f *ScopeDecl) FindPackagesDecl() []*PackageDecl {
	var ret []*PackageDecl
	for _, t := range f.Filter(glanglexer.PackageToken) {
		if x, ok := t.(*PackageDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}
func (f *ScopeDecl) FindImplementsTypes() []*ImplementDecl {
	var ret []*ImplementDecl
	for _, t := range f.Filter(glanglexer.ImplementsToken) {
		if x, ok := t.(*ImplementDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}
func (f *ScopeDecl) FindStructsTypes() []*StructDecl {
	var ret []*StructDecl
	for _, t := range f.Filter(glanglexer.StructToken) {
		if x, ok := t.(*StructDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}
func (f *ScopeDecl) FindTemplatesTypes() []*TemplateDecl {
	var ret []*TemplateDecl
	for _, t := range f.Filter(glanglexer.TemplateToken) {
		if x, ok := t.(*TemplateDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}
func (f *ScopeDecl) FindInterfaces() []*InterfaceDecl {
	var ret []*InterfaceDecl
	for _, t := range f.Filter(glanglexer.InterfaceToken) {
		if x, ok := t.(*InterfaceDecl); ok {
			ret = append(ret, x)
		}
	}
	return ret
}
func (f *ScopeDecl) FindFuncs() []*FuncDecl {
	var ret []*FuncDecl
	for _, t := range f.Filter(glanglexer.FuncToken) {
		if x, ok := t.(*FuncDecl); ok && x.IsTemplated() == false {
			ret = append(ret, x)
		}
	}
	return ret
}
func (f *ScopeDecl) FindTemplateFuncs() []FuncDeclarer {
	var ret []FuncDeclarer
	for _, t := range f.Filter(glanglexer.TplOpenToken) {
		if x, ok := t.(*TemplateFuncDecl); ok && x.IsDefine() == false {
			ret = append(ret, x)
		}
	}
	for _, t := range f.Filter(glanglexer.FuncToken) {
		if x, ok := t.(*FuncDecl); ok && x.IsTemplated() {
			ret = append(ret, x)
		}
	}
	return ret
}
func (f *ScopeDecl) FindDefineFuncs() []*TemplateFuncDecl {
	var ret []*TemplateFuncDecl
	for _, t := range f.Filter(glanglexer.TplOpenToken) {
		if x, ok := t.(*TemplateFuncDecl); ok && x.IsDefine() {
			ret = append(ret, x)
		}
	}
	return ret
}

type StrDecl struct {
	ScopeDecl
	Src string
}

func (f *StrDecl) FinalizeErr(err *genericinterperter.ParseError) error {
	return &genericinterperter.StringSyntaxError{Src: f.Src, Filepath: "<noname>", ParseError: *err}
}

func (p *StrDecl) GetName() string {
	return "noname"
}

type FileDecl struct {
	ScopeDecl
	Name string
}

func (f *FileDecl) FinalizeErr(err *genericinterperter.ParseError) error {
	return &genericinterperter.FileSyntaxError{Src: f.GetName(), ParseError: *err}
}

func (p *FileDecl) GetName() string {
	return p.Name
}

type PackageDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
	Name genericinterperter.Tokener
}

func (p *PackageDecl) String() string {
	return p.Expression.String()
}

func (p *PackageDecl) GetName() string {
	return p.Name.GetValue()
}
func NewPackageDecl(t genericinterperter.Tokener) *PackageDecl {
	return &PackageDecl{
		Tokener: t,
	}
}

type BodyBlockDecl struct {
	genericinterperter.Tokener
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
func NewBodyBlockDecl(t genericinterperter.Tokener) *BodyBlockDecl {
	return &BodyBlockDecl{
		Tokener: t,
	}
}

type StructDecl struct {
	genericinterperter.Tokener
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
func NewStructDecl(t genericinterperter.Tokener) *StructDecl {
	return &StructDecl{
		Tokener: t,
	}
}

type TemplateDecl struct {
	genericinterperter.Tokener
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

func NewTemplateDecl(t genericinterperter.Tokener) *TemplateDecl {
	return &TemplateDecl{
		Tokener: t,
	}
}

type InterfaceDecl struct {
	genericinterperter.Tokener
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
func NewInterfaceDecl(t genericinterperter.Tokener) *InterfaceDecl {
	return &InterfaceDecl{
		Tokener: t,
	}
}

type ImplementDecl struct {
	genericinterperter.Tokener
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
func NewImplementDecl(t genericinterperter.Tokener) *ImplementDecl {
	return &ImplementDecl{
		Tokener: t,
	}
}

type TemplateFuncDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
	Func     *FuncDecl
	Modifier *BodyBlockDecl
}

func NewTemplateFuncDecl(t genericinterperter.Tokener) *TemplateFuncDecl {
	return &TemplateFuncDecl{
		Tokener: t,
	}
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

type FuncDeclarer interface {
	genericinterperter.Expressioner
	IsMethod() bool
	IsTemplated() bool
	GetName() string
	// GetSlugName() string
	GetReceiverType() *IdentifierDecl
	GetReceiver() *PropsBlockDecl
	String() string
	GetModifier() *BodyBlockDecl
	GetBody() *BodyBlockDecl
	GetArgs() *PropsBlockDecl
	GetPos() genericinterperter.TokenPos
}

type FuncDecl struct {
	genericinterperter.Tokener
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
func NewFuncDecl(t genericinterperter.Tokener) *FuncDecl {
	return &FuncDecl{
		Tokener: t,
	}
}

type SignsBlockDecl struct {
	genericinterperter.Tokener
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
func NewSignsBlockDecl(t genericinterperter.Tokener) *SignsBlockDecl {
	return &SignsBlockDecl{
		Tokener: t,
	}
}

type PropsBlockDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
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
func (p *PropsBlockDecl) Add(Name *IdentifierDecl, Type *IdentifierDecl) *PropDecl {
	prop := NewPropDecl(Name)
	prop.Name = Name
	prop.Type = Type
	prop.AddExpr(Name)
	prop.AddExpr(Type)
	p.Props = append(p.Props, prop)
	p.Expression.AddExpr(prop)
	return prop
}
func (p *PropsBlockDecl) AddT(Type *IdentifierDecl) *PropDecl {
	prop := NewPropDecl(Type)
	prop.Type = Type
	prop.AddExpr(Type)
	p.Props = append(p.Props, prop)
	p.Expression.AddExpr(prop)
	return prop
}
func NewPropsBlockDecl(t genericinterperter.Tokener) *PropsBlockDecl {
	return &PropsBlockDecl{
		Tokener: t,
	}
}

type AssignsBlockDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
	Assigns []*AssignDecl
}

func (p *AssignsBlockDecl) String() string {
	return p.Expression.String()
}
func (p *AssignsBlockDecl) Add(left, leftType, right genericinterperter.Tokener) *AssignDecl {
	as := NewAssignDecl(left)
	as.Left = left
	as.LeftType = leftType
	as.Right = right
	p.Assigns = append(p.Assigns, as)
	p.Expression.AddExpr(as)
	return as
}
func NewAssignsBlockDecl(t genericinterperter.Tokener) *AssignsBlockDecl {
	return &AssignsBlockDecl{
		Tokener: t,
	}
}

type AssignDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
	Left     genericinterperter.Tokener
	LeftType genericinterperter.Tokener
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
func (p *AssignDecl) GetRight() string {
	return p.Right.GetValue()
}
func NewAssignDecl(t genericinterperter.Tokener) *AssignDecl {
	return &AssignDecl{
		Tokener: t,
	}
}

type PropDecl struct {
	genericinterperter.Tokener
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
func NewPropDecl(t genericinterperter.Tokener) *PropDecl {
	return &PropDecl{
		Tokener: t,
	}
}

type IdentifierDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
}

var re = regexp.MustCompile("(<[^>]+>)")

func (p *IdentifierDecl) GetSlugName() string {
	s := p.Expression.String()
	s = re.ReplaceAllString(s, "")
	s = strings.Replace(s, "*", "", -1)
	return strings.TrimSpace(s)
}
func (p *IdentifierDecl) String() string {
	return p.Expression.String()
}
func NewIdentifierDecl(t genericinterperter.Tokener) *IdentifierDecl {
	return &IdentifierDecl{
		Tokener: t,
	}
}

type VarDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
}

func (p *VarDecl) String() string {
	return p.Expression.String()
}
func NewVarDecl(t genericinterperter.Tokener) *VarDecl {
	return &VarDecl{
		Tokener: t,
	}
}

type ConstDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
}

func (p *ConstDecl) String() string {
	return p.Expression.String()
}
func NewConstDecl(t genericinterperter.Tokener) *ConstDecl {
	return &ConstDecl{
		Tokener: t,
	}
}

type ExpressionDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
}

func (p *ExpressionDecl) String() string {
	return p.Expression.String()
}
func NewExpressionDecl(t genericinterperter.Tokener) *ExpressionDecl {
	return &ExpressionDecl{
		Tokener: t,
	}
}

type CommentGroupDecl struct {
	genericinterperter.Tokener
	genericinterperter.Expression
}

func (p *CommentGroupDecl) String() string {
	return p.Expression.String()
}
func (p *CommentGroupDecl) HasAny() bool {
	return len(p.Expression.Tokens) > 0
}

func NewCommentGroupDecl(t genericinterperter.Tokener) *CommentGroupDecl {
	c := &CommentGroupDecl{
		Tokener: t,
	}
	c.AddExpr(t)
	return c
}
