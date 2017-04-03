package glang

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
	gigolexer "github.com/mh-cbon/gigo/lexer/gigo"
	"github.com/mh-cbon/gigo/struct/glang"
	lexer "github.com/mh-cbon/state-lexer"
)

func TestOneFunc(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}
	funcs := d.FindFuncs()

	got := len(funcs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected func len wanted=%v, got=%v", wanted, got)
	}
	fn := funcs[0]
	sgot := fn.GetName()
	swanted := "tomate"
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}

	// genericinterperter.Dump(fn.Body)
	// os.Exit(1)
}

func TestOneFuncReceiver(t *testing.T) {

	str := `func (r *Receiver) tomate() {

		}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}
	funcs := d.FindFuncs()

	got := len(funcs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected func len wanted=%v, got=%v", wanted, got)
	}
	fn := funcs[0]
	sgot := fn.GetName()
	swanted := " tomate" // something to be adjusted here
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}

	got = len(fn.Receiver.Props)
	wanted = 1
	if wanted != got {
		t.Errorf("unexpected receiver len wanted=%v, got=%v", wanted, got)
	}

	receiver := fn.Receiver.Props[0]
	sgot = receiver.GetName()
	swanted = "r"
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = receiver.Type.String()
	swanted = " *Receiver" // something to be adjusted here
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}

	// genericinterperter.Dump(d, 0)
}

func TestOneFuncRepeater(t *testing.T) {

	str := `<:whatever>func (r *Receiver<:whatever>) tomate<:whatever>() {

		}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}
	funcs := d.FindTemplateFuncs()

	got := len(funcs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected func len wanted=%v, got=%v", wanted, got)
	}
	fn := funcs[0]
	sgot := fn.GetName()
	swanted := " tomate<:whatever>" // something to be adjusted here
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}

	got = len(fn.GetReceiver().Props)
	wanted = 1
	if wanted != got {
		t.Errorf("unexpected receiver len wanted=%v, got=%v", wanted, got)
	}

	receiver := fn.GetReceiver().Props[0]
	sgot = receiver.GetName()
	swanted = "r"
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = receiver.Type.String()
	swanted = " *Receiver<:whatever>" // something to be adjusted here
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}

	// genericinterperter.Dump(d, 0)
}

func TestOneStruct(t *testing.T) {

	str := `type tomate struct {}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}
	structs := d.FindStructsTypes()

	got := len(structs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected struct len wanted=%v, got=%v", wanted, got)
	}
	st := structs[0]
	sgot := st.GetName()
	swanted := "tomate"
	if swanted != sgot {
		t.Errorf("unexpected struct name wanted=%q, got=%q", swanted, sgot)
	}

	// genericinterperter.Dump(d, 0)
}

func TestOneStructTemplate(t *testing.T) {

	str := `type tomate struct {
		poireau<:Slice(.Todo)>
		*poireau<:Mutexed .>
}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}

	structs := d.FindStructsTypes()
	got := len(structs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected func len wanted=%v, got=%v", wanted, got)
	}

	st := structs[0]
	sgot := st.GetName()
	swanted := "tomate"
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}

	poireaux := st.Block.Poireaux
	got = len(poireaux)
	wanted = 2
	if wanted != got {
		t.Errorf("unexpected poireaux len wanted=%v, got=%v", wanted, got)
	}

	poireau0 := poireaux[0]
	sgot = poireau0.String()
	swanted = "poireau<:Slice(.Todo)>"
	if swanted != sgot {
		t.Errorf("unexpected poireau0 mutation wanted=%v, got=%v", swanted, sgot)
	}
	sgot = poireau0.GetImplementTemplate()
	swanted = "poireau<:Slice(.Todo)>"
	if swanted != sgot {
		t.Errorf("unexpected poireau0 mutation wanted=%v, got=%v", swanted, sgot)
	}
	bgot := poireau0.IsPointer()
	bwanted := false
	if bwanted != bgot {
		t.Errorf("unexpected poireau0 out wanted=%v, got=%v", bwanted, bgot)
	}

	poireau1 := poireaux[1]
	sgot = poireau1.String()
	swanted = "*poireau<:Mutexed .>"
	if swanted != sgot {
		t.Errorf("unexpected poireau1 mutation wanted=%v, got=%v", swanted, sgot)
	}
	sgot = poireau1.GetImplementTemplate()
	swanted = "*poireau<:Mutexed .>"
	if swanted != sgot {
		t.Errorf("unexpected poireau1 mutation wanted=%v, got=%v", swanted, sgot)
	}
	bgot = poireau1.IsPointer()
	bwanted = true
	if bwanted != bgot {
		t.Errorf("unexpected poireau1 out wanted=%v, got=%v", bwanted, bgot)
	}

	// genericinterperter.Dump(d, 0)
}

func TestOneBrokenStruct(t *testing.T) {
	str := `type tomate struct qsdqd{}`
	_, err := interpretString("tomate", str)
	if err == nil {
		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
	}
	// fmt.Printf("%#v", err)
	// fmt.Printf("%+v", err)
}

func TestOneStructWithProps(t *testing.T) {

	str := `type tomate struct {
		A string
		B int
	}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}

	structs := d.FindStructsTypes()
	got := len(structs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected structs len wanted=%v, got=%v", wanted, got)
	}

	st := structs[0]
	sgot := st.GetName()
	swanted := "tomate"
	if swanted != sgot {
		t.Errorf("unexpected struct name wanted=%q, got=%q", swanted, sgot)
	}

	props := st.Block.Props
	got = len(props)
	wanted = 2
	if wanted != got {
		t.Errorf("unexpected props len wanted=%v, got=%v", wanted, got)
	}

	propA := props[0]
	sgot = propA.GetName()
	swanted = "A"
	if swanted != sgot {
		t.Errorf("unexpected propA name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = propA.Type.String()
	swanted = " string"
	if swanted != sgot {
		t.Errorf("unexpected propA value wanted=%q, got=%q", swanted, sgot)
	}

	propB := props[1]
	sgot = propB.GetName()
	swanted = "B"
	if swanted != sgot {
		t.Errorf("unexpected propB name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = propB.Type.String()
	swanted = " int"
	if swanted != sgot {
		t.Errorf("unexpected propB value wanted=%q, got=%q", swanted, sgot)
	}
	// genericinterperter.Dump(d, 0)
}

func TestOneTemplate(t *testing.T) {

	str := `template Mutexed<:.Name> struct {
	  lock *sync.Mutex
	  embed <:.Name>
	}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}

	structs := d.FindTemplatesTypes()
	got := len(structs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected structs len wanted=%v, got=%v", wanted, got)
	}

	st := structs[0]
	sgot := st.GetName()
	swanted := "Mutexed<:.Name>"
	if swanted != sgot {
		t.Errorf("unexpected struct name wanted=%q, got=%q", swanted, sgot)
	}

	props := st.Block.Props
	got = len(props)
	wanted = 2
	if wanted != got {
		t.Errorf("unexpected props len wanted=%v, got=%v", wanted, got)
	}

	propA := props[0]
	sgot = propA.GetName()
	swanted = "lock"
	if swanted != sgot {
		t.Errorf("unexpected propA name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = propA.Type.String()
	swanted = " *sync.Mutex"
	if swanted != sgot {
		t.Errorf("unexpected propA value wanted=%q, got=%q", swanted, sgot)
	}

	propB := props[1]
	sgot = propB.GetName()
	swanted = "embed"
	if swanted != sgot {
		t.Errorf("unexpected propB name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = propB.Type.String()
	swanted = " <:.Name>"
	if swanted != sgot {
		t.Errorf("unexpected propB value wanted=%q, got=%q", swanted, sgot)
	}
	// genericinterperter.Dump(d, 0)
}

func TestOneInterface(t *testing.T) {

	str := `type todosProvider interface {
	  Push(Todo) int
	  Remove(Todo) int
	}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}

	infs := d.FindInterfaces()
	got := len(infs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected interfaces len wanted=%v, got=%v", wanted, got)
	}

	inf := infs[0]
	sgot := inf.GetName()
	swanted := "todosProvider"
	if swanted != sgot {
		t.Errorf("unexpected interface name wanted=%q, got=%q", swanted, sgot)
	}

	signs := inf.Block.Signs
	got = len(signs)
	wanted = 2
	if wanted != got {
		t.Errorf("unexpected signs len wanted=%v, got=%v", wanted, got)
	}

	sign1 := signs[0]
	sgot = sign1.GetName()
	swanted = "Push"
	if swanted != sgot {
		t.Errorf("unexpected sign1 name wanted=%q, got=%q", swanted, sgot)
	}
	// tbd params.

	sign2 := signs[1]
	sgot = sign2.GetName()
	swanted = "Remove"
	if swanted != sgot {
		t.Errorf("unexpected sign2 name wanted=%q, got=%q", swanted, sgot)
	}
	// tbd params.
	// genericinterperter.Dump(d, 0)
}

func TestOneImplements(t *testing.T) {

	str := `
	type Todos implements<:Mutexed (Slice .Todo "Name")> {
	  // it reads as a mutexed list of todo.
	}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}

	impls := d.FindImplementsTypes()
	got := len(impls)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected implements len wanted=%v, got=%v", wanted, got)
	}

	impl := impls[0]
	sgot := impl.GetName()
	swanted := "Todos"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}

	sgot = impl.GetImplementTemplate()
	swanted = `<:Mutexed (Slice .Todo "Name")>`
	if swanted != sgot {
		t.Errorf("unexpected impl mutator wanted=%q, got=%q", swanted, sgot)
	}

	// tbd props
	// genericinterperter.Dump(d, 0)
}

func TestOnePackageDecl(t *testing.T) {

	str := ``
	str = fmt.Sprintf("package %v\n\n%v", "tomate", str)
	d, err := interpretStringWithPkgDecl("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}

	pkgs := d.FindPackagesDecl()
	got := len(pkgs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected packages len wanted=%v, got=%v", wanted, got)
	}

	pkg := pkgs[0]
	sgot := pkg.GetName()
	swanted := "tomate"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}

	// genericinterperter.Dump(d, 0)
}

func TestOneVarDecl(t *testing.T) {

	str := `
var x = "content"
var y string = "content1"
var (
	z = "tomate"
	v string = "tomate1"
)
var z = interface{}{}
var d = struct{Name string}{Name: ""}
`
	str = fmt.Sprintf("package %v\n\n%v", "tomate", str)
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}

	founds := d.FindVarDecl()
	got := len(founds)
	wanted := 5
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}

	found := founds[0]
	got = len(found.GetAssignments())
	wanted = 1
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}
	assignment := found.GetAssignments()[0]
	sgot := assignment.GetLeft()
	swanted := "x"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetAssign()
	swanted = "="
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetRight()
	swanted = "\"content\""
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}

	// -
	found = founds[1]
	got = len(found.GetAssignments())
	wanted = 1
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}
	assignment = found.GetAssignments()[0]
	sgot = assignment.GetLeft()
	swanted = "y"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetAssign()
	swanted = "="
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetRight()
	swanted = "\"content1\""
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}

	// -
	found = founds[2]
	got = len(found.GetAssignments())
	wanted = 2
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}
	assignment = found.GetAssignments()[0]
	sgot = assignment.GetLeft()
	swanted = "z"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetAssign()
	swanted = "="
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetRight()
	swanted = "\"tomate\""
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	//-
	assignment = found.GetAssignments()[1]
	sgot = assignment.GetLeft()
	swanted = "v"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetAssign()
	swanted = "="
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetRight()
	swanted = "\"tomate1\""
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}

	// -
	found = founds[3]
	got = len(found.GetAssignments())
	wanted = 1
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}
	assignment = found.GetAssignments()[0]
	sgot = assignment.GetLeft()
	swanted = "z"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetAssign()
	swanted = "="
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetRight()
	swanted = "interface" // TBD
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}

	// -
	found = founds[4]
	got = len(found.GetAssignments())
	wanted = 1
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}
	assignment = found.GetAssignments()[0]
	sgot = assignment.GetLeft()
	swanted = "d"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetAssign()
	swanted = "="
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetRight()
	swanted = "struct" // TBD
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}

}

func TestOneBrokenVar(t *testing.T) {
	str := `var `
	_, err := interpretString("tomate", str)
	if err == nil {
		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
	}
	// fmt.Printf("%#v", err)
	// fmt.Printf("%+v", err)

	str = `var interface = ""`
	_, err = interpretString("tomate", str)
	if err == nil {
		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
	}

	str = `var struct = ""`
	_, err = interpretString("tomate", str)
	if err == nil {
		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
	}
}

func TestOneConstDecl(t *testing.T) {

	str := `
const y  = "content1"
const (
	numberToken lexer.TokenType = iota
	wsToken
)
`
	str = fmt.Sprintf("package %v\n\n%v", "tomate", str)
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		// t.Errorf("%+v\n", err)
	}
	// genericinterperter.Dump(d)

	founds := d.FindConstDecl()
	got := len(founds)
	wanted := 2
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}

	found := founds[0]
	got = len(found.GetAssignments())
	wanted = 1
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}
	assignment := found.GetAssignments()[0]
	sgot := assignment.GetLeft()
	swanted := "y"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetAssign()
	swanted = "="
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetRight()
	swanted = "\"content1\""
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}

	// -
	found = founds[1]
	got = len(found.GetAssignments())
	wanted = 2
	if wanted != got {
		t.Errorf("unexpected var len wanted=%v, got=%v", wanted, got)
	}
	assignment = found.GetAssignments()[0]
	sgot = assignment.GetLeft()
	swanted = "numberToken"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetLeftType()
	swanted = "lexer.TokenType"
	if swanted != sgot {
		t.Errorf("unexpected impl name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetAssign()
	swanted = "="
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
	sgot = assignment.GetRight()
	swanted = "iota"
	if swanted != sgot {
		t.Errorf("unexpected assignment wanted=%q, got=%q", swanted, sgot)
	}
}

func TestOneBrokenConst(t *testing.T) {
	str := `const `
	_, err := interpretString("tomate", str)
	if err == nil {
		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
	}
	// fmt.Printf("%#v", err)
	// fmt.Printf("%+v", err)

	str = `const interface = ""`
	_, err = interpretString("tomate", str)
	if err == nil {
		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
	}

	str = `const struct = ""`
	_, err = interpretString("tomate", str)
	if err == nil {
		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
	}
}

func TestNoPackageDecl(t *testing.T) {

	str := ``
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
	}

	pkgs := d.FindPackagesDecl()
	got := len(pkgs)
	wanted := 0
	if wanted != got {
		t.Errorf("unexpected packages len wanted=%v, got=%v", wanted, got)
	}

	// genericinterperter.Dump(d, 0)
}

func interpretString(pkgName, content string) (*glang.StrDecl, error) {

	var buf bytes.Buffer
	buf.WriteString(content)
	reader := makeLexerReader(&buf)
	//reader = prettyPrinterLexer(reader)

	interpret := NewGigoInterpreter()
	return interpret.ProcessStr(content, reader)
}

func interpretStringWithPkgDecl(pkgName, content string) (*glang.StrDecl, error) {

	var buf bytes.Buffer
	buf.WriteString(content)
	reader := makeLexerReader(&buf)
	//reader = prettyPrinterLexer(reader)

	interpret := NewGigoInterpreter()
	return interpret.ProcessStrWithPkgDecl(content, reader)
}

func makeLexerReader(r io.Reader) func() genericinterperter.Tokener {

	l := lexer.New(r, (gigolexer.New()).StartHere)
	l.ErrorHandler = func(e string) {}

	return genericinterperter.PositionnedTokenReader(l.NextToken)
}

func prettyPrinterLexer(reader func() genericinterperter.Tokener) func() genericinterperter.Tokener {

	namer := genericinterperter.TokenerName(gigolexer.TokenName)
	reader = genericinterperter.PrettyPrint(reader, namer)

	return reader
}
