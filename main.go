// go generate super charged on steroids
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"

	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
	glanginterpreter "github.com/mh-cbon/gigo/interpreter/glang"
	genericlexer "github.com/mh-cbon/gigo/lexer/generic"
	gigolexer "github.com/mh-cbon/gigo/lexer/gigo"
	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
	glang "github.com/mh-cbon/gigo/struct/glang"
	"github.com/mh-cbon/state-lexer"
)

func main() {

	var symbol string
	flag.StringVar(&symbol, "symbol", "", "Find specified symbol name")

	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Wrong usage, should be")
		fmt.Println("go run main.go <cmd> <file>")
		fmt.Println("")
		fmt.Println("Available commands:")
		fmt.Println("dump: pretty print the interpretation result of a file")
		fmt.Println("gen: mutate a source file")
		panic("not enough arguments")
	}

	cmd := flag.Arg(0)
	path := flag.Arg(1)
	// f := must open os.Open("demo.gigo")

	fileDef := MustInterpretFile(path)

	if cmd == "str" || cmd == "s" {
		if symbol != "" {
			symbols := fileDef.FindSymbols(symbol)
			if len(symbols) > 0 {
				fmt.Println(symbols[0])
			} else {
				fmt.Println("No symbol found for ", symbol)
			}
		} else {
			fmt.Println(fileDef.String())
		}
	} else if cmd == "dump" || cmd == "d" {
		if symbol != "" {
			symbols := fileDef.FindSymbols(symbol)
			if len(symbols) > 0 {
				glanginterpreter.Dump(symbols[0])
			} else {
				fmt.Println("No symbol found for ", symbol)
			}
		} else {
			glanginterpreter.Dump(fileDef)
		}
	} else if cmd == "gen" || cmd == "g" {
		newDecl, err := mutate(fileDef)
		if err != nil {
			fmt.Printf("%#v\n", err)
			panic(err)
		}
		if symbol != "" {
			symbols := newDecl.FindSymbols(symbol)
			if len(symbols) > 0 {
				fmt.Println(symbols[0])
			} else {
				fmt.Println("No symbol found for ", symbol)
			}
		} else {
			fmt.Println(newDecl.String())
		}
	}
}

func mutate(fileDef *glang.FileDecl) (glang.ScopeReceiver, error) {

	allTplsFuncs := map[string]interface{}{
		"joinexpr": func(glue string, tokens interface{}) string {
			t := []genericinterperter.Tokener{}
			switch yy := tokens.(type) {
			case []genericinterperter.Tokener:
				t = append(t, yy...)
			case []*glang.PropDecl:
				for _, xx := range yy {
					t = append(t, xx)
				}
			case []*glang.IdentifierDecl:
				for _, xx := range yy {
					t = append(t, xx)
				}
			}
			ret := []string{}
			for _, xx := range t {
				ret = append(ret, xx.String())
			}
			return strings.Join(ret, glue)
		},
	}

	tplTypesFuncs := map[string]interface{}{}
	outData := &Tomate{
		implTplData: map[string]interface{}{},
	}

	/* At that moment the file is processed,
	all the template/type/struct/interface/func/ect declarations
	are well known.
	*/
	// prepare the source for its rendering

	var defineFunc []glang.FuncDeclarer
	structTypes := fileDef.FindStructsTypes()
	implTypes := fileDef.FindImplementsTypes()
	tplTypes := fileDef.FindTemplatesTypes()
	funcs := fileDef.FindFuncs()
	tplFuncs := fileDef.FindTemplateFuncs()
	defFuncs := fileDef.FindDefineFuncs()

	var attachMethod = func(m glang.FuncDeclarer) {
		for _, t := range tplTypes {
			if x, ok := m.GetReceiverType().First().(*glang.IdentifierDecl); ok {
				if t.GetSlugName() == x.GetSlugName() {
					t.AddMethod(m)
					return
				}
			}
		}
		panic("not found")
	}
	var attachImplMethod = func(m glang.FuncDeclarer) bool {
		if m.IsMethod() {
			for _, t := range implTypes {
				if x, ok := m.GetReceiverType().First().(*glang.IdentifierDecl); ok {
					if t.Name.GetSlugName() == x.GetSlugName() {
						t.AddMethod(m)
						return true
					}
				}
			}
		}
		// panic("not found")
		return false
	}

	// type XXX implements{}, needs to be replaced by a placeholder,
	// its template tokens values are changed to avoid further problems
	for _, i := range implTypes {
		name := fmt.Sprintf("placeholder%v", len(outData.placeholders))
		m := NewPlaceholderTypeMutation(name, i)
		outData.placeholders = append(outData.placeholders, m)
		fileDef.MustInsertAfter(i, m.PlaceholderDecl)
		fileDef.MustRemove(i)
		i.SetTokenValue(glanglexer.TplOpenToken, "<:")
		i.SetTokenValue(glanglexer.TplCloseToken, ":>")
	}
	// template XXX<Modifier> struct {}
	// are to be removed, they really just template expressions.
	for _, i := range tplTypes {
		fileDef.MustRemove(i)
		i.SetTokenValue(glanglexer.TplOpenToken, "<:")
		i.SetTokenValue(glanglexer.TplCloseToken, ":>")
	}
	// <Modifier> func ()
	// and
	// func(receiver<...>)...
	// are to be removed, they really just template expressions.
	// it also attaches the method to their type.
	for _, i := range tplFuncs {
		fileDef.MustRemove(i)
		i.SetTokenValue(glanglexer.TplOpenToken, "<:")
		i.SetTokenValue(glanglexer.TplCloseToken, ":>")
		// i.GetBody().SetTokenValue(glanglexer.GreaterToken, ":>") // trick unti fix.
		attachMethod(i)
	}
	// <define> func XXX ()
	// are to be removed because those funcs are injected into the template instances
	for _, i := range defFuncs {
		fileDef.MustRemove(i)
		i.SetTokenValue(glanglexer.TplOpenToken, "<:")
		i.SetTokenValue(glanglexer.TplCloseToken, ":>")
		defineFunc = append(defineFunc, i)

		//- define a template func
		// that func (tbd later) will be available in type declarations expressions like
		// - implement<>
		// - template<>
		name := i.GetName()
		tplTypesFuncs[name] = stubFunc(i.String())
		// the key difficulty in this feature is that the func string can not be
		// evaluated at runtime, so this whole template transforms step,
		// needs to be delayed to a new sub go program where the func body string can be written.
		// just refactroring of the current mess!
	}
	// regular go fund method are attached to ehir type.
	for _, i := range funcs {
		attachImplMethod(i)
	}

	for _, i := range structTypes {
		// declare regular structs as data protperties
		outData.implTplData[i.GetName()] = i
	}

	// for every declarations
	// - template XXXX struct{}
	// - func(of the template)...
	// do
	//- create a template.Template of its string
	//- create a template.Func of its mutation
	funcsForTypesMutators := map[string]interface{}{}
	for k, v := range tplTypesFuncs {
		funcsForTypesMutators[k] = v
	}
	for k, v := range allTplsFuncs {
		funcsForTypesMutators[k] = v
	}
	for _, i := range tplTypes {
		outData.tplTypesMutators = append(outData.tplTypesMutators, &TypeMutator{
			Decl:  i,
			funcs: funcsForTypesMutators,
		})
	}

	// need to remove comments, they are not understood by template.Template,
	// and if they contain the template syntax, it breaks becasue template evaluate them.
	// on the other hand, GigoInterpreter does not interpret comments, so it can t see and manage those
	// problematic strings. :/
	// finally the idea is to lacehold the comments, its kind of noop, works well.
	x := placeholdComments(genericlexer.CommentBlockToken, fileDef, "blockcomments")
	outData.placeholders = append(outData.placeholders, x...)
	y := placeholdComments(genericlexer.CommentLineToken, fileDef, "linecomments")
	outData.placeholders = append(outData.placeholders, y...)

	tplContent := fileDef.String()

	// execute the modified file tree with a taylor made template context.
	tpl := makeTplOfSource("gigo", tplContent, allTplsFuncs)

	var out bytes.Buffer
	if err := tpl.Execute(&out, outData); err != nil {
		return nil, genericinterperter.NewStringTplSyntaxError(err, "gigo", tplContent)
	}
	return InterpretString(fileDef.GetName(), out.String())
}

type Tomate struct {
	placeholders     []mutationExecuter
	tplTypesMutators []*TypeMutator
	implTplData      map[string]interface{}
}

func (t *Tomate) getPlaceholder(name string) mutationExecuter {
	for _, p := range t.placeholders {
		if p.getName() == name {
			return p
		}
	}
	return nil
}
func (t *Tomate) GetResult(name string) string {
	pl := t.getPlaceholder(name)

	if pl != nil {
		res, err := pl.execute(t.tplTypesMutators, t.implTplData)
		if err != nil {
			panic(err)
		}
		return res
	}
	return "not found"
}

type TemplateTplDot struct {
	*glang.StructDecl
	Args []interface{}
}

func (t *TemplateTplDot) ArgType(s interface{}) string {
	return reflect.TypeOf(s).Name()
}

func makeLexerReader(r io.Reader) genericinterperter.TokenerReaderOK {

	l := lexer.New(r, (gigolexer.New()).StartHere)
	l.ErrorHandler = func(e string) {}

	return genericinterperter.NewReadTokenWithPos(l)
}

func prettyPrinterLexer(reader genericinterperter.TokenerReaderOK) genericinterperter.TokenerReaderOK {

	namer := genericinterperter.TokenerName(gigolexer.TokenName)
	reader = genericinterperter.NewReadNPrettyPrint(reader, namer, os.Stdout)

	return reader
}

func InterpretFile(fileName string) (*glang.FileDecl, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := makeLexerReader(f)
	// reader = prettyPrinterLexer(reader)

	interpret := glanginterpreter.NewGigoInterpreter(reader)
	return interpret.ProcessFile(fileName)
}

func InterpretString(pkgName, content string) (*glang.StrDecl, error) {

	var buf bytes.Buffer
	buf.WriteString(content)
	reader := makeLexerReader(&buf)
	//reader = prettyPrinterLexer(reader)

	interpret := glanginterpreter.NewGigoInterpreter(reader)
	return interpret.ProcessStr(content)
}

func MustInterpretFile(fileName string) *glang.FileDecl {
	ret, err := InterpretFile(fileName)
	if err != nil {
		fmt.Printf("%#v\n", err)
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	return ret
}

func MustInterpretString(name, content string) *glang.StrDecl {
	ret, err := InterpretString(name, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
		panic(err)
	}
	return ret
}

var plToken lexer.TokenType = -200

func placeholderToken(name string, pos genericinterperter.TokenPos) *genericinterperter.TokenWithPos {
	tok := lexer.Token{
		Type:  plToken,
		Value: fmt.Sprintf("<:.GetResult \"%v\":>", name),
	}
	return genericinterperter.NewTokenWithPos(tok, pos.Line, pos.Pos)
}

func placeholdComments(T lexer.TokenType, src *glang.FileDecl, prefix string) []mutationExecuter {
	ret := []mutationExecuter{}
	for _, c := range src.FindAll(T) {
		name := fmt.Sprintf("placeholder%v%v", prefix, len(ret))
		m := NewPlaceholderMutation(name, c.GetTokens()[0])
		ret = append(ret, m)
		src.InsertAfter(c, m.PlaceholderDecl)
		src.Remove(c)
	}
	return ret
}

type mutationExecuter interface {
	execute(mutators []*TypeMutator, data interface{}) (string, error)
	getName() string
}

type placeholderMutation struct {
	OriginDecl      genericinterperter.Tokener
	PlaceholderDecl *genericinterperter.TokenWithPos
	Name            string
}

func (p *placeholderMutation) getName() string {
	return p.Name
}
func (p *placeholderMutation) execute(mutators []*TypeMutator, data interface{}) (string, error) {
	return p.OriginDecl.String(), nil
}

func NewPlaceholderMutation(name string, of genericinterperter.Tokener) *placeholderMutation {
	return &placeholderMutation{
		OriginDecl:      of,
		PlaceholderDecl: placeholderToken(name, of.GetPos()),
		Name:            name,
	}
}

type placeholderTypeMutation struct {
	mutation        *ImplTypeMutation
	PlaceholderDecl *genericinterperter.TokenWithPos
	Name            string
}

func (p *placeholderTypeMutation) getName() string {
	return p.Name
}
func (p *placeholderTypeMutation) execute(mutators []*TypeMutator, data interface{}) (string, error) {
	expr, err := p.mutation.mutate(mutators, data)
	res := ""
	if expr != nil {
		res = expr.String()
	}
	return res, err
}

func NewPlaceholderTypeMutation(name string, of *glang.ImplementDecl) *placeholderTypeMutation {
	return &placeholderTypeMutation{
		mutation:        &ImplTypeMutation{Decl: of},
		PlaceholderDecl: placeholderToken(name, of.GetPos()),
		Name:            name,
	}
}

func makeTplOfSource(name, src string, funcs map[string]interface{}) *template.Template {
	t, err := template.New(name).Funcs(funcs).Delims("<:", ":>").Parse(src)
	if err != nil {
		// fmt.Println(src)
		fErr := genericinterperter.NewStringTplSyntaxError(err, name, src)
		fmt.Printf("%#v\n", fErr)
		panic(fErr)
	}
	return t
}

func stubFunc(content string) func() error {
	return func() error {
		fmt.Println("template func content is")
		fmt.Println(content)
		return nil
	}
}

type TypeMutator struct {
	Decl  *glang.TemplateDecl
	funcs map[string]interface{}
}

func (t *TypeMutator) getTemplateStr() string {
	tplContent := ""
	// the template declares a type like this
	// template XXXX struct{}
	// it is needed to replace the template keyword by a type.
	// => type XXXX struct{}
	if y := t.Decl.GetToken(glanglexer.TemplateToken); y != nil {
		// y.SetType(glanglexer.TypeToken) // not needed to update
		y.SetValue("type")
	}
	tplContent += t.Decl.String()
	for _, m := range t.Decl.Methods {
		tplContent += m.String()
		if m.GetModifier() != nil { // test if there is a front modifier like <range $m :=...>
			tplContent += "<:end:>" // close the template expression, quick and dirty, but just works :)
		}
	}
	return tplContent
}
func (t *TypeMutator) execute(data interface{}) (string, error) {
	name := t.Decl.GetName()
	src := t.getTemplateStr()
	tpl := makeTplOfSource(name, src, t.funcs)
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	return buf.String(), err
}
func (t *TypeMutator) mutate(origin *glang.StructDecl, args ...interface{}) (*glang.StructDecl, error) {
	// the provided argument becomes the template root dot {{.}}
	arg := &TemplateTplDot{StructDecl: origin, Args: args}
	content, err := t.execute(arg)
	if err == nil {
		// we shall parse it
		newFileDef := MustInterpretString("random", content)
		newStruct := newFileDef.FindStructsTypes()[0] // a type for a type

		// should it be added to the current template data ?
		// structTypes = append(structTypes, newStruct)
		// note, it is expected the type gets added to the package repository.

		// dont forget to attach its method.
		for _, f := range newFileDef.FindFuncs() {
			newStruct.AddMethod(f)
		}
		return newStruct, nil
	}
	return origin, err
}

type ImplTypeMutation struct {
	scope genericinterperter.Expression
	Decl  *glang.ImplementDecl
	Res   []*glang.StructDecl
}

func (t *ImplTypeMutation) getMutationFunc(m *TypeMutator) func(origin *glang.StructDecl, args ...interface{}) (*glang.StructDecl, error) {
	return func(origin *glang.StructDecl, args ...interface{}) (*glang.StructDecl, error) {
		res, err := m.mutate(origin, args...)
		if err == nil {
			t.Res = append(t.Res, res)
		}
		return res, err
	}
}

func (t *ImplTypeMutation) mutate(mutators []*TypeMutator, data interface{}) (*glang.StrDecl, error) {
	// in a decl like implement<X Y Z>
	// X Y Z are func template of a template string "X Y Z"
	funcs := map[string]interface{}{}
	for _, m := range mutators {
		name := m.Decl.GetSlugName()
		funcs[name] = t.getMutationFunc(m)
	}

	tplContent := t.Decl.String()
	tpl := makeTplOfSource("gigo", tplContent, funcs)
	if err := tpl.Execute(ioutil.Discard, data); err != nil {
		fmt.Println(tplContent)
		return nil, err
	}
	// once the template "X Y Z" invoked => new struct type is added to t.Res

	// finalize the implements instruction into a regular struct
	// it becomes regular go code.
	// from => type xxx impements<y u i>{}
	// to => type xxx struct{}
	i := t.Decl
	i.SetTokenValue(glanglexer.ImplementsToken, "struct")
	i.RemoveT(glanglexer.TplOpenToken) // get ride of the template mutations

	strDecl := &glang.StrDecl{}

	if len(t.Res) > 0 {
		// if any mutations is found, get the last one,
		// and apply it to the original type
		last := t.Res[len(t.Res)-1]

		tok := genericinterperter.NewTokenWithPos(lexer.Token{Type: genericlexer.WordToken, Value: last.GetName()}, 0, 0)
		nl := genericinterperter.NewTokenWithPos(lexer.Token{Type: glanglexer.NlToken, Value: "\n"}, 0, 0)
		ws := genericinterperter.NewTokenWithPos(lexer.Token{Type: genericlexer.WsToken, Value: "\t"}, 0, 0)

		// define the last genrated type as an underlying type of i
		ID := glang.NewExpressionDecl()
		name := glang.NewIdentifierDecl()
		name.AddExpr(tok)
		ID.AddExpr(name)
		i.GetBlock().Underlying = append(i.GetBlock().Underlying, ID)
		i.GetBlock().InsertAt(1, nl)
		i.GetBlock().InsertAt(2, ws)
		i.GetBlock().InsertAt(3, ID)

		// add every generated types and all of their methods to the string decl
		for _, r := range t.Res {
			strDecl.AddExprs(r.Tokens)
			nl := genericinterperter.NewTokenWithPos(lexer.Token{Type: glanglexer.NlToken, Value: "\n"}, 0, 0)
			strDecl.AddExpr(nl)
			for _, m := range r.Methods {
				strDecl.AddExprs(m.GetTokens())
				nl := genericinterperter.NewTokenWithPos(lexer.Token{Type: glanglexer.NlToken, Value: "\n"}, 0, 0)
				strDecl.AddExpr(nl)
			}
		}
	}
	// finally add the original modified I struct to the string decl
	strDecl.AddExpr(i)

	return strDecl, nil
}
