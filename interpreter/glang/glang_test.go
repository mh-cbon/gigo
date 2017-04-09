package glang

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
	genericlexer "github.com/mh-cbon/gigo/lexer/generic"
	gigolexer "github.com/mh-cbon/gigo/lexer/gigo"
	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
	"github.com/mh-cbon/gigo/struct/glang"
	lexer "github.com/mh-cbon/state-lexer"
)

func TestReadVarNameIdentifier(t *testing.T) {
	content := `var1
othervar
interfacevar
a_var
templated<:varx>
templated<:vary>
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadVarName(false, true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block, "var1")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarName(true, true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block, "othervar")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarName(false, true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block, "interfacevar")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarName(false, true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block, "a_var")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarName(true, true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block, "templated<:varx>")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarName(false, true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block, "templated")
	interpret.ReadBlock(glanglexer.TplOpenToken, glanglexer.GreaterToken)
	interpret.GetMany(glanglexer.NlToken)
}

func TestFailReadVarNameIdentifier(t *testing.T) {
	content := `var
interface
123
{
templated<:varx>
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadVarName(false, true, true)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(glanglexer.VarToken))
	mustNil(t, block)
	interpret.GetMany(glanglexer.VarToken, glanglexer.NlToken)

	block, err = interpret.ReadVarName(false, true, true)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(glanglexer.InterfaceToken))
	mustNil(t, block)
	interpret.GetMany(glanglexer.InterfaceToken, glanglexer.NlToken)

	block, err = interpret.ReadVarName(false, true, true)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(genericlexer.WordToken))
	mustNil(t, block)
	interpret.GetMany(genericlexer.WordToken, glanglexer.NlToken)

	block, err = interpret.ReadVarName(false, true, true)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(glanglexer.BraceOpenToken))
	mustNil(t, block)
	interpret.GetMany(glanglexer.BraceOpenToken, glanglexer.NlToken)

}

func TestReadTypeIdentifier(t *testing.T) {
	content := `string
int
int8
int16
int32
int64
uint
uint8
uint16
uint32
uint64
float
float32
float64
interface{}
interface { }
struct{}
struct { }
struct { name string }
struct { T }
struct { *T }
struct {
	*T
}
struct {
	*T
	name string
}
struct {
	*T
	name []string
}
struct {
	*T
	name2 struct{}
}
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "string")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "int")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "int8")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "int16")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "int32")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "int64")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "uint")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "uint8")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "uint16")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "uint32")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "uint64")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "float")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "float32")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "float64")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "interface{}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "interface { }")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct{}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct { }")
	blockEq(t, block.First(), "{ }")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct { name string }")
	blockEq(t, block.First(), "{ name string }")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct { T }")
	blockEq(t, block.First(), "{ T }")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct { *T }")
	blockEq(t, block.First(), "{ *T }")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct {\n\t*T\n}")
	blockEq(t, block.First(), "{\n\t*T\n}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct {\n\t*T\n\tname string\n}")
	blockEq(t, block.First(), "{\n\t*T\n\tname string\n}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct {\n\t*T\n\tname []string\n}")
	blockEq(t, block.First(), "{\n\t*T\n\tname []string\n}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	StringEq(t, block, "struct {\n\t*T\n\tname2 struct{}\n}")
	blockEq(t, block.First(), "{\n\t*T\n\tname2 struct{}\n}")
	interpret.GetMany(glanglexer.NlToken)
}

func TestFailReadTypeIdentifier(t *testing.T) {
	content := `name
[
*name
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	mustNotNil(t, interpret.Peek(genericlexer.WordToken))
	mustNil(t, block)
	interpret.GetMany(genericlexer.WordToken, glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	mustNotNil(t, interpret.Peek(glanglexer.BracketOpenToken))
	mustNil(t, block)
	interpret.GetMany(glanglexer.BracketOpenToken, glanglexer.NlToken)

	block, err = interpret.ReadTypeIdentifier(false)
	mustNotErr(t, err)
	mustNotNil(t, interpret.Peek(glanglexer.MulToken))
	mustNil(t, block)
	interpret.GetMany(glanglexer.MulToken, genericlexer.WordToken, glanglexer.NlToken)
}

func TestReadTypeName(t *testing.T) {
	content := `T
*T
[]T
[]*T
[][]*T
[][]*T<:cccc>
[]string
[]*string
string
[]t.t
item.<:$a>
<:$a>
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadTypeName(false, true)
	mustNotErr(t, err)
	StringEq(t, block, "T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(false, true)
	mustNotErr(t, err)
	StringEq(t, block, "*T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(false, true)
	mustNotErr(t, err)
	StringEq(t, block, "[]T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(false, true)
	mustNotErr(t, err)
	StringEq(t, block, "[]*T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(false, true)
	mustNotErr(t, err)
	StringEq(t, block, "[][]*T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(true, true)
	mustNotErr(t, err)
	StringEq(t, block, "[][]*T<:cccc>")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "[]string")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "[]*string")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(true, true)
	mustNotErr(t, err)
	identifierNameEq(t, block.First(), "string")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(true, true)
	mustNotErr(t, err)
	StringEq(t, block, "[]t.t")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(true, true)
	mustNotErr(t, err)
	StringEq(t, block, "item.<:$a>")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeName(true, true)
	mustNotErr(t, err)
	StringEq(t, block, "<:$a>")
	interpret.GetMany(glanglexer.NlToken)

}

func TestFailReadTypeName(t *testing.T) {
	content := `_
(
{
[
[][
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadTypeName(false, true)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(genericlexer.WordToken))
	mustNil(t, block)
	interpret.GetMany(genericlexer.WordToken, glanglexer.NlToken)

	block, err = interpret.ReadTypeName(false, true)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(glanglexer.ParenOpenToken))
	mustNil(t, block)
	interpret.GetMany(glanglexer.ParenOpenToken, glanglexer.NlToken)

	block, err = interpret.ReadTypeName(false, true)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(glanglexer.BraceOpenToken))
	mustNil(t, block)
	interpret.GetMany(glanglexer.BraceOpenToken, glanglexer.NlToken)

	block, err = interpret.ReadTypeName(false, true)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(glanglexer.BracketOpenToken))
	mustNil(t, block)
	interpret.GetMany(glanglexer.BracketOpenToken, glanglexer.NlToken)

}

func TestReadTypeValue(t *testing.T) {
	content := `T{}
T { }
T{x: y}
T{x, y}
&T{}
[]T{T{x: y}}
[]*T{&T{x: y}}
&[]T{}
&[]*T{&T{x: y}}
&[]*T{&T{x: y}, &T{x: y}}
z[]
[]t.t{}
item.<:$a>
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "T{}")
	StringEq(t, block.First(), "T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "T { }")
	StringEq(t, block.First(), "T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "T{x: y}")
	StringEq(t, block.First(), "T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "T{x, y}")
	StringEq(t, block.First(), "T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "&T{}")
	StringEq(t, block.First(), "&T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "[]T{T{x: y}}")
	StringEq(t, block.First(), "[]T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "[]*T{&T{x: y}}")
	StringEq(t, block.First(), "[]*T")
	StringEq(t, block.GetExprs()[0], "[]*T")
	StringEq(t, block.GetExprs()[1], "{")
	StringEq(t, block.GetExprs()[2], "&")
	StringEq(t, block.GetExprs()[3], "T")
	StringEq(t, block.GetExprs()[4], "{")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "&[]T{}")
	StringEq(t, block.First(), "&[]T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "&[]*T{&T{x: y}}")
	StringEq(t, block.First(), "&[]*T")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "&[]*T{&T{x: y}, &T{x: y}}")
	StringEq(t, block.First(), "&[]*T")
	interpret.GetMany(glanglexer.NlToken)

	interpret.blockscope.AddVar("z")
	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "z[]")
	StringEq(t, block.First(), "z")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(false)
	mustNotErr(t, err)
	StringEq(t, block, "[]t.t{}")
	StringEq(t, block.First(), "[]t.t")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadTypeValue(true)
	mustNotErr(t, err)
	StringEq(t, block, "item.<:$a>")
	interpret.GetMany(glanglexer.NlToken)
}

func TestReadExpressionBlock(t *testing.T) {
	content := `somevarname
typewhatever{}
typewhatever { }
somecallexpr(x, z)
somecallexpr (x, z)
some.call.expr()
some.call.expr ()
someother[:i]
someother [:i]
T{}
T{what, ever}
i++
i--
[]t{}
func(){}
item.<:$a>
-1
+1
1+1
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "somevarname")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "typewhatever{}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "typewhatever { }")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	cexpr := block.First().(*glang.CallExpr)
	identifierNameEq(t, cexpr.ID, "somecallexpr")
	lenEq(t, 2, len(cexpr.Params.Params))
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	cexpr = block.First().(*glang.CallExpr)
	identifierNameEq(t, cexpr.ID, "somecallexpr")
	lenEq(t, 2, len(cexpr.Params.Params))
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	cexpr = block.First().(*glang.CallExpr)
	identifierNameEq(t, cexpr.ID, "some.call.expr")
	lenEq(t, 0, len(cexpr.Params.Params))
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	cexpr = block.First().(*glang.CallExpr)
	identifierNameEq(t, cexpr.ID, "some.call.expr")
	lenEq(t, 0, len(cexpr.Params.Params))
	interpret.GetMany(glanglexer.NlToken)

	interpret.blockscope.AddVar("someother")
	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "someother[:i]")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "someother [:i]")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "T{}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "T{what, ever}")
	interpret.GetMany(glanglexer.NlToken)

	interpret.blockscope.AddVar("i")
	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "i++")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "i--")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "[]t{}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(false, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "func(){}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "item.<:$a>")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "-1")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "+1")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadExpressionBlock(true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "1+1")
	// Dump(block)
	interpret.GetMany(glanglexer.NlToken)
}

func TestReadBinaryExpr(t *testing.T) {
	content := `ok
ok && false
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadBinaryExpressionBlock(false, glanglexer.NlToken)
	mustNotErr(t, err)
	StringEq(t, block, "ok")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadBinaryExpressionBlock(false, glanglexer.NlToken)
	mustNotErr(t, err)
	StringEq(t, block, "ok")
	interpret.GetMany(glanglexer.NlToken)

}

func TestReadIfStmt(t *testing.T) {
	content := `if true {}
if i>5 { }
if i>5 && false { }
if i>5 && false==true { }
if i:=5;i<5 { }
if item.<:$a> == item.<:$a> {}
if item.<:$a> == x<:$a> {}
if item.<:$a> == <:$a>x {}
if x, ok := t.(*PackageDecl); ok {}
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadIfStmt(false)
	mustNotErr(t, err)
	condEq(t, block, "true")
	bodyEq(t, block, "{}")

	interpret.GetMany(glanglexer.NlToken)
	block, err = interpret.ReadIfStmt(false)
	mustNotErr(t, err)
	condEq(t, block, "i>5")
	bodyEq(t, block, "{ }")

	interpret.GetMany(glanglexer.NlToken)
	block, err = interpret.ReadIfStmt(false)
	mustNotErr(t, err)
	condEq(t, block, "i>5 && false")
	bodyEq(t, block, "{ }")

	interpret.GetMany(glanglexer.NlToken)
	block, err = interpret.ReadIfStmt(false)
	mustNotErr(t, err)
	condEq(t, block, "i>5 && false==true")
	bodyEq(t, block, "{ }")

	interpret.GetMany(glanglexer.NlToken)
	block, err = interpret.ReadIfStmt(false)
	mustNotErr(t, err)
	initEq(t, block, "i:=5")
	condEq(t, block, "i<5")
	bodyEq(t, block, "{ }")

	interpret.GetMany(glanglexer.NlToken)
	interpret.blockscope.AddVar("item")
	block, err = interpret.ReadIfStmt(true)
	mustNotErr(t, err)
	condEq(t, block, "item.<:$a> == item.<:$a>")
	bodyEq(t, block, "{}")

	interpret.GetMany(glanglexer.NlToken)
	interpret.blockscope.AddVar("x<:$a>")
	block, err = interpret.ReadIfStmt(true)
	mustNotErr(t, err)
	condEq(t, block, "item.<:$a> == x<:$a>")
	bodyEq(t, block, "{}")

	interpret.GetMany(glanglexer.NlToken)
	interpret.blockscope.AddVar("<:$a>x")
	block, err = interpret.ReadIfStmt(true)
	mustNotErr(t, err)
	condEq(t, block, "item.<:$a> == <:$a>x")
	bodyEq(t, block, "{}")

	interpret.GetMany(glanglexer.NlToken)
	block, err = interpret.ReadIfStmt(true)
	mustNotErr(t, err)
	condEq(t, block, "ok")
	StringEq(t, block, "if x, ok := t.(*PackageDecl); ok {}")
}

func TestFailReadIfStmt(t *testing.T) {
	content := `if true;i<5 { }
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadIfStmt(false)
	mustErr(t, err, block)
	mustNotNil(t, interpret.Peek(glanglexer.TrueToken))
	mustNil(t, block)
	interpret.PeekUntil(glanglexer.NlToken)
	interpret.Read(glanglexer.NlToken)
}

func TestReadElseIfStmt(t *testing.T) {
	content := `else if true {}
else
if true {}
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadElseStmt(false)
	mustNotErr(t, err)
	condEq(t, block, "true")
	bodyEq(t, block, "{}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadElseStmt(false)
	mustNotErr(t, err)
	condEq(t, block, "true")
	bodyEq(t, block, "{}")
	interpret.GetMany(glanglexer.NlToken)
}

func TestReadElseStmt(t *testing.T) {
	content := `else {}
else{}
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadElseStmt(false)
	mustNotErr(t, err)
	bodyEq(t, block, "{}")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadElseStmt(false)
	mustNotErr(t, err)
	bodyEq(t, block, "{}")
	interpret.GetMany(glanglexer.NlToken)
}

func TestAssignExpr(t *testing.T) {
	content := `x := "r"
y := 5
z = somevar
u = somecall()
u = someother[:1]
a, b := someother[:1], "rr"
a := func(a,b){return "z"}
`
	interpret := makeRawInterpreter(content)
	interpret.blockscope.AddVar("someother")

	block, err := interpret.ReadAssignExpr(false, true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	lenEq(t, 1, len(block.IDs))
	lenEq(t, 1, len(block.Values))
	identifierNameEq(t, block.IDs[0], "x")
	StringEq(t, block.Values[0], `"r"`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadAssignExpr(false, true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	lenEq(t, 1, len(block.IDs))
	lenEq(t, 1, len(block.Values))
	identifierNameEq(t, block.IDs[0], "y")
	StringEq(t, block.Values[0], `5`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadAssignExpr(false, true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	lenEq(t, 1, len(block.IDs))
	lenEq(t, 1, len(block.Values))
	identifierNameEq(t, block.IDs[0], "z")
	StringEq(t, block.Values[0], `somevar`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadAssignExpr(false, true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	lenEq(t, 1, len(block.IDs))
	lenEq(t, 1, len(block.Values))
	identifierNameEq(t, block.IDs[0], "u")
	StringEq(t, block.Values[0], `somecall()`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadAssignExpr(false, true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	lenEq(t, 1, len(block.IDs))
	lenEq(t, 1, len(block.Values))
	identifierNameEq(t, block.IDs[0], "u")
	StringEq(t, block.Values[0], `someother[:1]`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadAssignExpr(false, true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	lenEq(t, 2, len(block.IDs))
	lenEq(t, 2, len(block.Values))
	identifierNameEq(t, block.IDs[0], "a")
	StringEq(t, block.Values[0], `someother[:1]`)
	identifierNameEq(t, block.IDs[1], "b")
	StringEq(t, block.Values[1], `"rr"`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadAssignExpr(false, true, glanglexer.SemiColonToken)
	mustNotErr(t, err)
	StringEq(t, block, "a := func(a,b){return \"z\"}")
	interpret.GetMany(glanglexer.NlToken)
	// Dump(block)
}

func TestForExpr(t *testing.T) {
	content := `for{}
for true {}
for call() {}
for expr.call() {}
for i:=0;i<5;i++ {}
for range some {}
for range []sometype{} {}
for x,y := range []sometype{} {}
for range someother.w {}
for _, t := range f.Tokens {}
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for{}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for true {}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for call() {}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for expr.call() {}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for i:=0;i<5;i++ {}`)
	interpret.GetMany(glanglexer.NlToken)

	interpret.blockscope.AddVar("some")
	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for range some {}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for range []sometype{} {}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for x,y := range []sometype{} {}`)
	interpret.GetMany(glanglexer.NlToken)

	interpret.blockscope.AddVar("someother")
	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for range someother.w {}`)
	interpret.GetMany(glanglexer.NlToken)

	interpret.blockscope.AddVar("f")
	block, err = interpret.ReadForBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `for _, t := range f.Tokens {}`)
	interpret.GetMany(glanglexer.NlToken)
}

func TestReadParenExpr(t *testing.T) {
	content := `()
(true)
(true, false)
([]struct{name string}{
	name:"",
}, call(), "text")
(true,)
(true,
false)
(true true,)
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadParenExprBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `()`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenExprBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `(true)`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenExprBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, `(true, false)`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenExprBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, "([]struct{name string}{\n\tname:\"\",\n}, call(), \"text\")")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenExprBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, "(true,)")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenExprBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, "(true,\nfalse)")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenExprBlock(false)
	mustNotErr(t, err)
	StringEq(t, block, "(true true,)")
	interpret.GetMany(glanglexer.NlToken)

	// Dump(block)
}

func TestReadParenDecl(t *testing.T) {
	content := `()
(zz error)
(zz string)
(zz []T)
(zz []T, xx string)
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	mustNotErr(t, err)
	StringEq(t, block, `()`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	mustNotErr(t, err)
	StringEq(t, block, `(zz error)`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	mustNotErr(t, err)
	StringEq(t, block, `(zz string)`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	mustNotErr(t, err)
	StringEq(t, block, `(zz []T)`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	mustNotErr(t, err)
	StringEq(t, block, `(zz []T, xx string)`)
	interpret.GetMany(glanglexer.NlToken)

	// Dump(block)
}

func TestReadVarDecl(t *testing.T) {
	content := `var x string = "eee"
var x string = []tt{}
var x string = []t.t{}
var x string = call.expr()
var (
	x string = "e"
	x int = 3
)
var ret []*PackageDecl
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadVarDecl(false)
	mustNotErr(t, err)
	StringEq(t, block, `var x string = "eee"`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarDecl(false)
	mustNotErr(t, err)
	StringEq(t, block, `var x string = []tt{}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarDecl(false)
	mustNotErr(t, err)
	StringEq(t, block, `var x string = []t.t{}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarDecl(false)
	mustNotErr(t, err)
	StringEq(t, block, `var x string = call.expr()`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarDecl(false)
	mustNotErr(t, err)
	StringEq(t, block, "var (\n\tx string = \"e\"\n\tx int = 3\n)")
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadVarDecl(false)
	mustNotErr(t, err)
	StringEq(t, block, "var ret []*PackageDecl")
	interpret.GetMany(glanglexer.NlToken)

	// Dump(block)
}

func TestReadFunc(t *testing.T) {
	content := `func (s <:.Name>Slice) Push(item <:.Name>) int {}
func (t *Todos) Hello(){fmt.Println("Hello")}
func (m Mutexed<:$.Name>) <:$m.Name>(<:$m.GetArgsBlock | joinexpr ",">) <:$m.Out>{}
func (f *ScopeDecl) GrepLine(line int) []genericinterperter.Tokener {}
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadFuncDecl(true, false)
	mustNotErr(t, err)
	StringEq(t, block, `func (s <:.Name>Slice) Push(item <:.Name>) int {}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadFuncDecl(false, false)
	mustNotErr(t, err)
	StringEq(t, block, `func (t *Todos) Hello(){fmt.Println("Hello")}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadFuncDecl(true, false)
	mustNotErr(t, err)
	StringEq(t, block, `func (m Mutexed<:$.Name>) <:$m.Name>(<:$m.GetArgsBlock | joinexpr ",">) <:$m.Out>{}`)
	interpret.GetMany(glanglexer.NlToken)

	block, err = interpret.ReadFuncDecl(false, false)
	mustNotErr(t, err)
	StringEq(t, block, `func (f *ScopeDecl) GrepLine(line int) []genericinterperter.Tokener {}`)
	interpret.GetMany(glanglexer.NlToken)
}

func TestReadTemplateExprDecl(t *testing.T) {
	content := `<:range $a := .Args> func (s <:$.Name>Slice) FindBy<:$a>(<:$a> <:$.ArgType $a>) (<:$.Name>,bool) {
	  for i, item := range s.items {
	    if item.<:$a> == <:$a> {
	      return item, true
	    }
	  }
	  return <:$.Name>{}, false
	}`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadTemplateExprDecl()
	mustNotErr(t, err)
	StringEq(t, block, content)
	interpret.GetMany(glanglexer.NlToken)

	// Dump(block)
}

func TestReadExpressionsBlock(t *testing.T) {
	content := `{
	s.items = append(s.items, item)
	return len(s.items)
}
`
	interpret := makeRawInterpreter(content)

	block, err := interpret.ReadExpressionsBlock(true, glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
	mustNotErr(t, err)
	StringEq(t, block, "{\n\ts.items = append(s.items, item)\n\treturn len(s.items)\n}")
	interpret.GetMany(glanglexer.NlToken)

	// Dump(block)
}

func StringEq(t *testing.T, x interface{}, expected string) {
	if s, ok := x.(fmt.Stringer); ok {
		swant := expected
		sgot := s.String()
		if swant != sgot {
			t.Errorf("Unexpected String() content, got=%q, want=%q", sgot, swant)
			t.FailNow()
		}
	} else {
		t.Errorf("Unexpected node type, got=%T, want=%q", x, "fmt.Stringer")
		t.FailNow()
	}
}

func lenEq(t *testing.T, iwant, igot int) {
	if iwant != igot {
		t.Errorf("Unexpected len, got=%d, want=%d", igot, iwant)
		t.FailNow()
	}
}

func identifierNameEq(t *testing.T, x interface{}, expected string) {
	if ID, ok := x.(*glang.IdentifierDecl); ok {
		swant := expected
		sgot := ID.String()
		if swant != sgot {
			t.Errorf("Unexpected ID name, got=%q, want=%q", sgot, swant)
			t.FailNow()
		}
	} else {
		t.Errorf("Unexpected node type, got=%T, want=%q\n in %v", x, "*glang.IdentifierDecl", x)
		t.FailNow()
	}
}

func blockEq(t *testing.T, x interface{}, expected string) {
	if ID, ok := x.(glang.BlockBodyer); ok {
		swant := expected
		sgot := ID.GetBlock().String()
		if swant != sgot {
			t.Errorf("Unexpected Block, got=%q, want=%q", sgot, swant)
			t.FailNow()
		}
	} else {
		t.Errorf("Unexpected node type, got=%T, want=%q", x, "*glang.IdentifierDecl")
		t.FailNow()
	}
}

func bodyEq(t *testing.T, x glang.Bodyer, expected string) {
	swant := expected
	sgot := x.GetBody().String()
	if swant != sgot {
		t.Errorf("Unexpected Body, got=%q, want=%q", sgot, swant)
		t.FailNow()
	}
}

func condEq(t *testing.T, x glang.Conder, expected string) {
	swant := expected
	sgot := x.GetCond().String()
	if swant != sgot {
		t.Errorf("Unexpected Condition, got=%q, want=%q", sgot, swant)
		t.FailNow()
	}
}

func initEq(t *testing.T, x glang.Initer, expected string) {
	swant := expected
	sgot := x.GetInit().String()
	if swant != sgot {
		t.Errorf("Unexpected init, got=%q, want=%q", sgot, swant)
		t.FailNow()
	}
}

func mustNotErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
		t.FailNow()
	}
}

func mustNil(t *testing.T, x interface{}) {
	if x != nil && !reflect.ValueOf(x).IsNil() {
		t.Errorf("wanted <nil> got=%#v\n", x)
		t.FailNow()
	}
}

func mustNotNil(t *testing.T, x interface{}) {
	if x == nil {
		t.Errorf("wanted <not nil> got=%q\n", x)
		t.FailNow()
	}
}

func mustErr(t *testing.T, err error, some ...interface{}) {
	if err == nil {
		t.Errorf("wanted <err> got=%q\nin: %v\n", err, some)
		t.FailNow()
	}
}

func makeRawInterpreter(content string) *GigoInterpreter {
	var buf bytes.Buffer
	buf.WriteString(content)
	reader := makeLexerReader(&buf)
	// reader = prettyPrinterLexer(reader)

	interpret := NewGigoInterpreter(reader)
	interpret.Scope = &glang.StrDecl{Src: content}
	return interpret
}

func TestOneFunc(t *testing.T) {

	str := `func tomate() {
		var expr string = "some"
		expr2 := "other"
}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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

	// Dump(fn.Body)
	// os.Exit(1)
}

func TestOneFuncReceiver(t *testing.T) {

	str := `func (r *Receiver) tomate() {

		}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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

	// Dump(d, 0)
}

func TestOneFuncRepeater(t *testing.T) {

	str := `<:whatever>func (r *Receiver<:whatever>) tomate<:whatever>(s []string) {
	    for i, items := range s.items {
	      if item.<:$a> == "" {
	        return item, true
	      }
	    }
	    return <:$.Name>{}, false
		}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
		return
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

	// Dump(d)
}

func TestOneFuncTemplate(t *testing.T) {

	str := `// create new Method Push of type .
func (s <:.Name>Slice) Push(item <:.Name>) int {
  s.items = append(s.items, item)
  return len(s.items)
}
`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
		return
	}
	funcs := d.FindTemplateFuncs()

	got := len(funcs)
	wanted := 1
	if wanted != got {
		t.Errorf("unexpected func len wanted=%v, got=%v", wanted, got)
	}
	fn := funcs[0]
	sgot := fn.GetName()
	swanted := " Push" // something to be adjusted here
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
	swanted = "s"
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}
	sgot = receiver.Type.String()
	swanted = " <:.Name>Slice" // something to be adjusted here
	if swanted != sgot {
		t.Errorf("unexpected func name wanted=%q, got=%q", swanted, sgot)
	}

	// Dump(d)
}

func TestOneExpr(t *testing.T) {

	str := `
func Push (){
  s.items = append(s.items, item)
  return len(s.items)
}

func (s <:.Name>Slice) Index(search <:.Name>) int {
  for i, item := range s.items {
    if item == search {
      return i
    }
  }
  return -1
}

func (s <:.Name>Slice) RemoveAt(i index) int {
	s.items = append(s.items[:i], s.items[i+1:]...)
}

func (s <:.Name>Slice) Remove(item <:.Name>) int {
  if i:= s.Index(item); i > -1 {
    s.RemoveAt(i)
    return i
  }
  return -1
}
`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
		return
	}
	// funcs := d.FindFuncs()
	tfuncs := d.FindTemplateFuncs()

	// Dump(funcs[0].GetBody())
	// Dump(tfuncs[1].GetBody())
	// Dump(d)
	fmt.Printf("%T\n", tfuncs[0])
	// fmt.Println(tfuncs[1])
	// fmt.Println(tfuncs[2])
}

func TestOneStruct(t *testing.T) {

	str := `type tomate struct {}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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
}

func TestOneStructTemplate(t *testing.T) {

	str := `type tomate struct {
		poireau<:Slice(.Todo)>
		*poireau<:Mutexed .>
}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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

	// Dump(d, 0)
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
		t.Errorf("%+v\n", err)
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
	swanted = "string"
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
	swanted = "int"
	if swanted != sgot {
		t.Errorf("unexpected propB value wanted=%q, got=%q", swanted, sgot)
	}
	// Dump(d, 0)
}

func TestOneTemplate(t *testing.T) {

	str := `template Mutexed<:.Name> struct {
	  lock *sync.Mutex
	  embed <:.Name>
	}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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
	swanted = "*sync.Mutex"
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
	swanted = "<:.Name>"
	if swanted != sgot {
		t.Errorf("unexpected propB value wanted=%q, got=%q", swanted, sgot)
	}
	// Dump(d, 0)
}

func TestOneInterface(t *testing.T) {

	str := `type todosProvider interface {
	  Push(Todo) int
	  Remove(Todo) int
	}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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
	// Dump(d, 0)
}

func TestOneImplements(t *testing.T) {

	str := `
	type Todos implements<:Mutexed (Slice .Todo "Name")> {
	  // it reads as a mutexed list of todo.
	}`
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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
	// Dump(d, 0)
}

func TestOnePackageDecl(t *testing.T) {

	str := ``
	str = fmt.Sprintf("package %v\n\n%v", "tomate", str)
	d, err := interpretStringWithPkgDecl("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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

	// Dump(d, 0)
}

func TestOneVarDecl(t *testing.T) {

	str := `
var x = "content"
var y string = "content1"
var (
	z = "tomate"
	v string = "tomate1"
)
var z = []somewhat{}
var d = struct{Name string}{Name: ""}
`
	str = fmt.Sprintf("package %v\n\n%v", "tomate", str)
	d, err := interpretString("tomate", str)
	if err != nil {
		t.Errorf("%#v\n", err)
		t.Errorf("%+v\n", err)
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
	swanted = "[]somewhat{}"
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
	swanted = "struct{Name string}{Name: \"\"}"
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
		t.Errorf("%+v\n", err)
	}
	// Dump(d)

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

// func TestOneBrokenConst(t *testing.T) {
// 	str := `const `
// 	_, err := interpretString("tomate", str)
// 	if err == nil {
// 		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
// 	}
// 	// fmt.Printf("%#v", err)
// 	// fmt.Printf("%+v", err)
//
// 	str = `const interface = ""`
// 	_, err = interpretString("tomate", str)
// 	if err == nil {
// 		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
// 	}
//
// 	str = `const struct = ""`
// 	_, err = interpretString("tomate", str)
// 	if err == nil {
// 		t.Errorf("unexpected err wanted=%v, got=%v", "<notnil>", err)
// 	}
// }

// func TestNoPackageDecl(t *testing.T) {
//
// 	str := ``
// 	d, err := interpretString("tomate", str)
// 	if err != nil {
// 		t.Errorf("%#v\n", err)
// 		t.Errorf("%+v\n", err)
// 	}
//
// 	pkgs := d.FindPackagesDecl()
// 	got := len(pkgs)
// 	wanted := 0
// 	if wanted != got {
// 		t.Errorf("unexpected packages len wanted=%v, got=%v", wanted, got)
// 	}
//
// 	// Dump(d, 0)
// }

func interpretString(pkgName, content string) (*glang.StrDecl, error) {

	var buf bytes.Buffer
	buf.WriteString(content)
	reader := makeLexerReader(&buf)
	// reader = prettyPrinterLexer(reader)
	// reader = protected(reader)

	interpret := NewGigoInterpreter(reader)
	return interpret.ProcessStr(content)
}

func interpretStringWithPkgDecl(pkgName, content string) (*glang.StrDecl, error) {

	var buf bytes.Buffer
	buf.WriteString(content)
	reader := makeLexerReader(&buf)
	//reader = prettyPrinterLexer(reader)
	// reader = genericinterperter.NewReadProtected(reader)

	interpret := NewGigoInterpreter(reader)
	return interpret.ProcessStrWithPkgDecl(content)
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
