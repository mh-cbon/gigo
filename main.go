package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	// f := must open os.Open("demo.gigo")

	fileName := "demo.gigo.go"
	fileDef, err := InterpretFile(fileName)
	if err != nil {
		panic(err)
	}

	/* At that moment the file is processed,
	all the template/type/struct/interface/func/ect declarations
	are well known.
	*/

	// prepare the source for its rendering

	var defineFunc []glang.FuncDeclarer
	pkgDecl := fileDef.FindPackageDecl()
	structTypes := fileDef.FindStructsTypes()
	implTypes := fileDef.FindImplementsTypes()
	tplTypes := fileDef.FindTemplatesTypes()
	funcs := fileDef.FindFuncs()
	tplFuncs := fileDef.FindTemplateFuncs()
	defFuncs := fileDef.FindDefineFuncs()

	var attachMethod = func(m glang.FuncDeclarer) {
		for _, t := range tplTypes {
			if t.GetSlugName() == m.GetReceiverType().GetSlugName() {
				t.AddMethod(m)
				return
			}
		}
		panic("not found")
	}
	var attachImplMethod = func(m glang.FuncDeclarer) bool {
		if m.IsMethod() {
			for _, t := range implTypes {
				if t.Name.GetSlugName() == m.GetReceiverType().GetSlugName() {
					t.AddMethod(m)
					return true
				}
			}
		}
		// panic("not found")
		return false
	}
	placeholders := map[string]*glang.ImplementDecl{}

	// type XXX implements{}, needs to be replaced by a placeholder,
	// its template tokens values are changed to avoid further problems
	for _, i := range implTypes {
		plname, pl := placeholderToken(i)
		placeholders[plname] = i
		fileDef.InsertAfter(i, pl)
		i.SetTokenValue(glanglexer.TplOpenToken, "<:")
		i.SetTokenValue(glanglexer.TplCloseToken, ":>")
	}
	// template XXX<Modifier> struct {}
	// are to be removed, they really just template expressions.
	for _, i := range tplTypes {
		if !fileDef.Remove(i) {
			panic("r")
		}
		i.SetTokenValue(glanglexer.TplOpenToken, "<:")
		i.SetTokenValue(glanglexer.TplCloseToken, ":>")
	}
	// <Modifier> func ()
	// and
	// func(receiver<...>)...
	// are to be removed, they really just template expressions.
	// it also attaches the method to their type.
	for _, i := range tplFuncs {
		if !fileDef.Remove(i) {
			panic("r")
		}
		i.SetTokenValue(glanglexer.TplOpenToken, "<:")
		i.SetTokenValue(glanglexer.TplCloseToken, ":>")
		attachMethod(i)
	}
	// <define> func XXX ()
	// are to be removed because those funcs are injected into the template instances
	for _, i := range defFuncs {
		if !fileDef.Remove(i) {
			panic("r")
		}
		i.SetTokenValue(glanglexer.TplOpenToken, "<:")
		i.SetTokenValue(glanglexer.TplCloseToken, ":>")
		defineFunc = append(defineFunc, i)
	}
	// regular go fund method are attached to ehir type.
	for _, i := range funcs {
		attachImplMethod(i)
	}

	/* At that moment the file is processed,
	all the template/type/struct/interface/func/ect declarations
	are well known.
	*/

	tplTypesTpl := map[string]*template.Template{}
	tplTypesFuncs := map[string]interface{}{}
	implTplFuncs := map[string]interface{}{}
	implTplData := map[string]interface{}{}
	implTplResults := map[string][]*glang.StructDecl{}

	/* The various elemnt of the source file
	are now injected as template element.
	<define> XXX		=> becomes a template.TemplateFunc
	template XXX<Modifier> 	=> becomes a template where the content = struct+methods decl
	<templat expr> func (receiver)(...) 	=> is associated to its type
	type XXX implements<Mutator()> 	=> becomes a template of <Mutator()> only,
																			- it defines template types as functions
																			- it defines regular struct into data member
	*/

	for _, i := range defFuncs {
		name := i.GetName()
		content := i.String()
		tplTypesFuncs[name] = func() error {
			fmt.Println("template func content is")
			fmt.Println(content)
			return nil
		}
	}

	currentImplName := "" // the current implement being resolved in the origin file
	for _, i := range tplTypes {
		// create a template.Template of the content of the template type.
		tplName := fmt.Sprintf("%v%v", "tplType", i.GetSlugName())
		tplContent := ""
		// the template declares a type like this
		// template XXXX struct{}
		// it is needed to replace the template keyword by a type.
		if y := i.GetToken(glanglexer.TemplateToken); y != nil {
			// y.SetType(glanglexer.TypeToken)
			y.SetValue("type")
		}
		tplContent += i.String()
		for _, m := range i.Methods {
			tplContent += m.String()
			if m.GetModifier() != nil { // test if there is a front modifier like <range $m :=...>
				tplContent += "<:end:>" // close the template expression, quick and dirty, but just works :)
			}
		}
		t, err2 := template.New(tplName).Funcs(tplTypesFuncs).Delims("<:", ":>").Parse(tplContent)
		if err2 != nil {
			panic(err2)
		}
		tplTypesTpl[tplName] = t

		// install the template type as a func for an implement declaration
		name := i.GetSlugName()
		implTplFuncs[name] = func(origin *glang.StructDecl) (*glang.StructDecl, error) {
			// the provided argument becomes the template root dot {{.}}
			tpl := tplTypesTpl["tplType"+name]
			var buf bytes.Buffer
			err2 := tpl.Execute(&buf, origin)
			if err2 == nil {
				// we shall parse it
				pkgName := pkgDecl.GetName()
				newFileDef, err3 := InterpretString(pkgName, buf.String())
				if err3 != nil {
					panic(err3)
				}
				pkgD := newFileDef.FindPackageDecl()
				newFileDef.Remove(pkgD)
				newStruct := newFileDef.FindStructsTypes()[0] // a type for a type
				// should it be added to the current template data ?
				// note, it is expected the type gets added to the package repository.
				// structTypes = append(structTypes, newStruct)

				// dont forget to attach its method.
				for _, f := range newFileDef.FindFuncs() {
					newStruct.AddMethod(f)
				}

				// add it to the currentImpl being processed.
				if _, ok := implTplResults[currentImplName]; ok == false {
					implTplResults[currentImplName] = []*glang.StructDecl{}
				}
				implTplResults[currentImplName] = append(implTplResults[currentImplName], newStruct)

				return newStruct, nil
			}
			return origin, err
		}
	}
	for _, i := range structTypes {
		// declare regular structs as data protperties
		implTplData[i.GetName()] = i
	}
	for _, i := range implTypes {
		// foreach implements<...> declaration,
		// create a template of <...>, where
		// functions are type mutators(origin),
		// and data are regular structs
		fmt.Println("----------------------------------")
		fmt.Println("----------------------------------")
		currentImplName = i.GetName()
		tplContent := i.ImplementTemplate.String()
		t, err2 := template.New("gigo").Funcs(implTplFuncs).Delims("<:", ":>").Parse(tplContent)
		if err2 != nil {
			panic(err2)
		}
		err3 := t.Execute(ioutil.Discard, implTplData)
		if err3 != nil {
			panic(err3)
		}
		// finalzie the implements instruction into a regular struct
		i.SetTokenValue(glanglexer.ImplementsToken, "struct")
		i.RemoveT(glanglexer.TplOpenToken) // get ride of the template mutations
		if len(implTplResults[currentImplName]) > 0 {
			xx := implTplResults[currentImplName][len(implTplResults[currentImplName])-1]

			tok := genericinterperter.NewTokenWithPos(lexer.Token{Type: genericlexer.WordToken, Value: xx.GetName()}, 0, 0)
			nl := genericinterperter.NewTokenWithPos(lexer.Token{Type: glanglexer.NlToken, Value: "\n"}, 0, 0)
			ws := genericinterperter.NewTokenWithPos(lexer.Token{Type: genericlexer.WsToken, Value: "\t"}, 0, 0)

			ID := glang.NewIdentifierDecl(tok)
			ID.AddExpr(tok)
			i.GetBlock().Underlying = append(i.GetBlock().Underlying, ID)
			i.GetBlock().InsertAt(1, nl)
			i.GetBlock().InsertAt(2, ws)
			i.GetBlock().InsertAt(3, ID)
		}
		currentImplName = ""
	}

	// fmt.Println((fileDef.String()))

	tplContent := fileDef.String()
	// executee the file content.
	data := &Tomate{
		placeholders:   placeholders,
		implTypes:      implTypes,
		implTplResults: implTplResults,
	}

	t, err2 := template.New("gigo").Funcs(implTplFuncs).Delims("<:", ":>").Parse(tplContent)
	if err2 != nil {
		panic(err2)
	}
	err3 := t.Execute(os.Stdout, data)
	// err3 := t.Execute(ioutil.Discard, data)
	if err3 != nil {
		panic(err3)
	}
	// fmt.Println(funcs[100])

	// genericinterperter.Dump(fileDef, 0)
	// // fmt.Println(fileDef.Dump()[:])
	// fmt.Println(len(fileDef.String()))
	// // fmt.Println(len(fileDef.String()[:950]))
	// fmt.Println(len(fileDef.Defs))
	// fmt.Println(len(fileDef.Tokens.Tokens))
	// fmt.Println(fileDef.Dump())
	// fmt.Println(fileDef.Defs[0].String()[:20])

}

type Tomate struct {
	placeholders   map[string]*glang.ImplementDecl
	implTypes      []*glang.ImplementDecl
	implTplResults map[string][]*glang.StructDecl
}

func (t *Tomate) GetResult(pl string) string {
	if impl, ok := t.placeholders[pl]; ok {
		for _, i := range t.implTypes {
			if i.GetName() == impl.GetName() {
				str := "\n"
				for _, x := range t.implTplResults[i.GetName()] {
					str += x.String() + "\n"
					for _, m := range x.Methods {
						str += m.String() + "\n"
					}
				}
				return str
			}
		}
	}
	return "not found"
}

func InterpretReader(name string, r io.Reader) (*glang.FileDecl, error) {

	l := lexer.New(r, (gigolexer.New()).StartHere)
	l.ErrorHandler = func(e string) {}

	// namer := genericinterperter.TokenerName(gigolexer.TokenName)

	reader := genericinterperter.PositionnedTokenReader(l.NextToken)
	// reader = genericinterperter.PrettyPrint(reader, namer)

	interpret := glanginterpreter.NewGigoInterpreter()
	fileDef := interpret.ProcessFile(name, reader)

	return fileDef, nil
}

func InterpretFile(fileName string) (*glang.FileDecl, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return InterpretReader(fileName, f)
}

func InterpretString(pkgName, content string) (*glang.FileDecl, error) {

	content = fmt.Sprintf("package %v\n\n%v", pkgName, content)

	var buf bytes.Buffer
	buf.WriteString(content)
	return InterpretReader("random", &buf)
}

var plToken lexer.TokenType = -200
var plIndex = 0

func placeholderToken(of genericinterperter.Tokener) (string, *genericinterperter.TokenWithPos) {
	name := fmt.Sprintf("placholder%v", plIndex)
	tok := lexer.Token{
		Type:  plToken,
		Value: fmt.Sprintf("<:.GetResult \"%v\":>", name),
	}
	plIndex++
	return name, genericinterperter.NewTokenWithPos(tok, of.GetPos().Line, of.GetPos().Pos)
}
