package glang

import (
	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
	genericlexer "github.com/mh-cbon/gigo/lexer/generic"
	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
	glang "github.com/mh-cbon/gigo/struct/glang"
	lexer "github.com/mh-cbon/state-lexer"
)

type PackageProvider interface {
	AddToPackage(string, genericinterperter.ScopeReceiver)
}

type GigoInterpreter struct {
	genericinterperter.Interpreter
	packages PackageProvider
}

func NewGigoInterpreter() *GigoInterpreter {
	return &GigoInterpreter{
		Interpreter: *genericinterperter.NewInterpreter(),
		packages:    &glang.SimplePackageRepository{},
	}
}
func (I *GigoInterpreter) ProcessFile(file string, reader genericinterperter.TokenerReader) *glang.FileDecl {
	fileDef := &glang.FileDecl{Name: file}

	var tokens []genericinterperter.Tokener
	for {
		if next := reader(); next != nil {
			tokens = append(tokens, next)
		} else {
			break
		}
	}

	I.Process(fileDef, tokens)

	return fileDef
}
func (I *GigoInterpreter) Process(
	scope genericinterperter.ScopeReceiver,
	tokens []genericinterperter.Tokener) error {

	I.Tokens = tokens
	I.Scope = scope

	I.MustDoPackageDecl()

	for {
		if I.Ended() {
			I.Scope.AddExprs(I.Emit())
			break
		}

		I.ReadMany(
			genericlexer.CommentLineToken,
			genericlexer.CommentBlockToken,
			genericlexer.WsToken,
			glanglexer.NlToken)

		if tok := I.Read(glanglexer.TypeToken); tok != nil {

			I.KeepPreviousComment()
			I.Scope.AddExprs(I.Emit())

			I.ReadMany(
				glanglexer.TypeToken,
				glanglexer.NlToken,
				genericlexer.CommentLineToken,
				genericlexer.CommentBlockToken,
				genericlexer.WsToken)
			preTokens := I.Emit()

			name := I.MustReadIdentifierDecl(false)

			I.ReadMany(genericlexer.WsToken)

			if typeTok := I.Peek(glanglexer.StructToken); typeTok != nil {

				sDecl := I.ReadStructDecl(false)
				sDecl.Name = name
				sDecl.PrependExpr(name)
				sDecl.PrependExprs(preTokens)
				I.Scope.AddExpr(sDecl)

			} else if typeTok := I.Peek(glanglexer.InterfaceToken); typeTok != nil {

				sDecl := I.ReadInterfaceDecl()
				sDecl.Name = name
				sDecl.PrependExpr(name)
				sDecl.PrependExprs(preTokens)
				I.Scope.AddExpr(sDecl)

			} else if tok := I.Peek(glanglexer.ImplementsToken); tok != nil {

				implDecl := I.ReadImplDecl()
				implDecl.Name = name
				implDecl.PrependExpr(name)
				implDecl.PrependExprs(preTokens)
				I.Scope.AddExpr(implDecl)
			}
		} else if tok := I.Peek(glanglexer.ImportToken); tok != nil {

			I.Read(glanglexer.ImportToken)
			I.ReadMany(
				glanglexer.ImportToken, // wip
				glanglexer.NlToken,
				genericlexer.CommentLineToken,
				genericlexer.CommentBlockToken,
				genericlexer.WsToken)
			I.Scope.AddExprs(I.Emit())
			I.ReadBlock(glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
			I.ReadMany(
				glanglexer.NlToken,
				genericlexer.WsToken)
			I.Scope.AddExprs(I.Emit())

		} else if tok := I.Peek(glanglexer.TemplateToken); tok != nil {

			tplDecl := I.ReadTemplateDecl()
			I.Scope.AddExpr(tplDecl)

		} else if tok := I.Peek(glanglexer.SmallerToken); tok != nil {

			// fn := &glang.TemplateFuncDecl{}
			fn := glang.NewTemplateFuncDecl(tok)
			I.KeepPreviousComment()

			I.ReadMany(
				glanglexer.NlToken,
				genericlexer.CommentLineToken,
				genericlexer.CommentBlockToken,
				genericlexer.WsToken)
			fn.AddExprs(I.Emit())

			block := I.ReadBodyBlock(glanglexer.SmallerToken, glanglexer.GreaterToken)
			block.Open.SetType(glanglexer.TplOpenToken)
			block.Close.SetType(glanglexer.TplCloseToken)
			block.AddExprs(I.Emit())
			fn.Modifier = block
			fn.AddExpr(block)
			I.ReadMany(genericlexer.WsToken)
			fn.AddExprs(I.Emit())

			fn.Func = I.ReadFuncDecl(true)
			fn.AddExpr(fn.Func)
			I.Scope.AddExpr(fn)

		} else if tok := I.Peek(glanglexer.FuncToken); tok != nil {

			I.KeepPreviousComment()
			I.Scope.AddExprs(I.Emit())

			I.ReadMany(
				glanglexer.NlToken,
				genericlexer.CommentLineToken,
				genericlexer.CommentBlockToken,
				genericlexer.WsToken)

			funcDecl := I.ReadFuncDecl(true)
			I.Scope.AddExpr(funcDecl)

		} else if tok := I.Peek(glanglexer.VarToken); tok != nil {

			varDecl := I.ReadVarDecl()
			I.Scope.AddExpr(varDecl)

		} else if tok := I.Peek(glanglexer.ConstToken); tok != nil {

			constDecl := I.ReadConstDecl()
			I.Scope.AddExpr(constDecl)

		} else if x := I.Next(); x == nil {
			I.Scope.AddExprs(I.Emit())
			break
		}
	}
	return nil
}

func (I *GigoInterpreter) MustDoPackageDecl() {

	for I.Peek(glanglexer.PackageToken) == nil {
		I.ReadMany(
			glanglexer.NlToken,
			genericlexer.CommentLineToken,
			genericlexer.CommentBlockToken,
			genericlexer.WsToken)
	}

	if tok := I.Read(glanglexer.PackageToken); tok != nil {
		decl := glang.NewPackageDecl(tok)
		decl.AddExprs(I.Emit())
		I.Scope.AddExpr(decl)
		decl.AddExprs(I.GetMany(genericlexer.WsToken))
		if name := I.Get(genericlexer.WordToken); name != nil {
			decl.Name = name
			decl.AddExpr(name)
			decl.AddExprs(I.GetMany(genericlexer.WsToken))
			decl.AddExprs(I.GetMany(glanglexer.NlToken))
			I.packages.AddToPackage(decl.GetName(), I.Scope)
		} else {
			panic("missing package name")
		}
	} else {
		panic("missing package decl")
	}
}

func (I *GigoInterpreter) KeepPreviousComment() {

	I.Rewind()
	nls := 0
	cl := 0
	for {
		I.Rewind()
		if token := I.Last(); token != nil {
			if token.GetType() == glanglexer.NlToken {
				nls++
				if nls > 1 {
					break
				}
			} else if token.GetType() == genericlexer.CommentBlockToken {
				break
			} else if token.GetType() == genericlexer.CommentLineToken {
				if cl == 0 {
					nls = 0
				}
				cl++
			}
		} else {
			break
		}
	}
}

func (I *GigoInterpreter) ReadPropsBlock(templated bool, open lexer.TokenType, close lexer.TokenType) *glang.PropsBlockDecl {

	var ret *glang.PropsBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return ret
	}
	count := 1

	ret = glang.NewPropsBlockDecl(openTok)

	for {
		I.ReadMany(
			genericlexer.WsToken,
			genericlexer.CommentBlockToken,
			genericlexer.CommentLineToken,
			glanglexer.NlToken,
		)
		if openTok := I.Read(open); openTok != nil {
			count++

		} else if closeTok := I.Read(close); closeTok != nil {
			count--
			if count == 0 {
				I.Rewind()
				break
			}

		} else {

			ret.AddExprs(I.Emit())

			ID := I.MustReadIdentifierDecl(templated)

			I.ReadMany(genericlexer.WsToken)
			IDType := I.ReadIdentifierDecl(templated)
			if IDType == nil {
				ret.AddUnderlying(IDType)
			} else {
				ret.Add(ID, IDType)
			}

		}
	}
	I.MustRead(close)
	ret.AddExprs(I.Emit())
	return ret
}

func (I *GigoInterpreter) MustReadAssignsBlock(canOmitRight bool, open, close lexer.TokenType) *glang.AssignsBlockDecl {
	r := I.ReadAssignsBlock(canOmitRight, open, close)
	if r == nil {
		panic(I.Debug("Missing assignements block: ", glanglexer.ParenOpenToken))
	}
	return r
}
func (I *GigoInterpreter) ReadAssignsBlock(canOmitRight bool, open lexer.TokenType, close lexer.TokenType) *glang.AssignsBlockDecl {

	var ret *glang.AssignsBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return ret
	}
	count := 1

	ret = glang.NewAssignsBlockDecl(openTok)

	for {
		I.ReadMany(genericlexer.WsToken,
			genericlexer.CommentBlockToken,
			genericlexer.CommentLineToken,
			glanglexer.NlToken,
		)
		if openTok := I.Read(open); openTok != nil {
			count++

		} else if closeTok := I.Read(close); closeTok != nil {
			count--
			if count == 0 {
				I.Rewind()
				break
			}

		} else {
			ret.AddExprs(I.Emit())

			left := I.MustReadIdentifierDecl(false)
			assignment := glang.NewAssignDecl(left)
			assignment.AddExpr(left)

			I.ReadMany(genericlexer.WsToken)
			assignment.AddExprs(I.Emit())

			leftType := I.ReadIdentifierDecl(false)
			if leftType != nil {
				assignment.LeftType = leftType
				assignment.AddExpr(leftType)
			}
			I.ReadMany(genericlexer.WsToken)
			assignment.AddExprs(I.Emit())

			var right *glang.ExpressionDecl
			if canOmitRight {
				if I.Read(glanglexer.AssignToken) != nil {
					assignment.AddExprs(I.Emit())
					I.ReadMany(genericlexer.WsToken)
					assignment.AddExprs(I.Emit())

					right = I.MustReadExpression()
					assignment.Right = right
					assignment.AddExpr(right)
				}
			} else {
				I.MustRead(glanglexer.AssignToken)
				I.ReadMany(genericlexer.WsToken)
				right = I.MustReadExpression()
				assignment.Right = right
				assignment.AddExpr(right)
			}

			assignment.AddExprs(I.Emit())
			ret.AddExpr(assignment)

		}
	}
	I.MustRead(close)
	ret.AddExprs(I.Emit())
	return ret
}

func (I *GigoInterpreter) ReadSignsBlock(open lexer.TokenType, close lexer.TokenType) *glang.SignsBlockDecl {

	var ret *glang.SignsBlockDecl
	I.ReadMany(
		glanglexer.NlToken,
		genericlexer.WsToken,
		genericlexer.CommentBlockToken,
		genericlexer.CommentLineToken,
	)

	openTok := I.Read(open)
	if openTok == nil {
		return ret
	}
	count := 1

	ret = glang.NewSignsBlockDecl(openTok)

	for {
		I.ReadMany(
			glanglexer.NlToken,
			genericlexer.WsToken,
			genericlexer.CommentBlockToken,
			genericlexer.CommentLineToken,
		)
		if openTok := I.Read(open); openTok != nil {
			count++

		} else if closeTok := I.Read(close); closeTok != nil {
			count--
			if count == 0 {
				I.Rewind()
				break
			}

		} else {
			ret.AddExprs(I.Emit())

			if ID := I.ReadIdentifierDecl(false); ID != nil {
				I.ReadMany(genericlexer.WsToken)
				if I.Peek(glanglexer.NlToken) != nil {
					ret.AddUnderlying(ID)
					I.ReadMany(glanglexer.NlToken)
					ret.AddExprs(I.Emit())
				} else {
					ret.Add(I.MustReadFuncSign(ID))
				}
			}

		}
	}
	I.MustRead(close)
	ret.AddExprs(I.Emit())
	return ret
}

func (I *GigoInterpreter) MustReadParenDecl(templated bool, open, close lexer.TokenType) *glang.PropsBlockDecl {
	r := I.ReadParenDecl(templated, open, close)
	if r == nil {
		panic(I.Debug("Missing parenthesis block: ", glanglexer.ParenOpenToken))
	}
	return r
}
func (I *GigoInterpreter) ReadParenDecl(templated bool, open, close lexer.TokenType) *glang.PropsBlockDecl {

	var ret *glang.PropsBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return ret
	}
	count := 1

	ret = glang.NewPropsBlockDecl(openTok)

	for {
		if openTok := I.Read(open); openTok != nil {
			count++

		} else if closeTok := I.Read(close); closeTok != nil {
			count--
			if count == 0 {
				I.Rewind()
				break
			}

		} else if semiColonTok := I.Read(glanglexer.SemiColonToken); semiColonTok != nil {

		} else {
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())

			I.ReadMany(glanglexer.NlToken, genericlexer.CommentBlockToken)

			ID := I.MustReadIdentifierDecl(templated)

			I.ReadMany(genericlexer.WsToken, genericlexer.CommentBlockToken)

			IDType := I.ReadIdentifierDecl(templated)

			var propdecl *glang.PropDecl
			if IDType == nil {
				propdecl = ret.AddT(ID)
			} else {
				propdecl = ret.Add(ID, IDType)
			}
			propdecl.AddExprs(I.Emit())

			I.ReadMany(genericlexer.WsToken, genericlexer.CommentBlockToken, glanglexer.SemiColonToken, glanglexer.NlToken)
		}
	}
	I.MustRead(close)
	ret.AddExprs(I.Emit())
	return ret
}

func (I *GigoInterpreter) ReadBodyBlock(open lexer.TokenType, close lexer.TokenType) *glang.BodyBlockDecl {

	var ret *glang.BodyBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return ret
	}
	count := 1
	ret = glang.NewBodyBlockDecl(openTok)
	ret.Open = openTok
	for {

		if openTok := I.Read(open); openTok != nil {
			count++

		} else if closeTok := I.Read(close); closeTok != nil {
			count--
			if count == 0 {
				ret.Close = closeTok
				break
			}

		} else {
			I.Next()
		}
	}
	return ret
}

func (I *GigoInterpreter) MustReadFuncSign(ID *glang.IdentifierDecl) *glang.FuncDecl {
	r := I.ReadFuncSign(ID)
	if r == nil {
		panic(I.Debug("Incomplete Func signature "+ID.String(), glanglexer.ParenOpenToken))
	}
	return r
}
func (I *GigoInterpreter) ReadFuncSign(ID *glang.IdentifierDecl) *glang.FuncDecl {
	ret := glang.NewFuncDecl(ID)
	ret.Name = ID
	ret.AddExpr(ID)

	I.ReadMany(genericlexer.WsToken)
	params := I.MustReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	params.AddExprs(I.Emit())
	ret.Params = params
	ret.AddExpr(params)

	I.ReadMany(genericlexer.WsToken)
	out := I.ReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	if out != nil {
		out.AddExprs(I.Emit())
		ret.Out = out
		ret.AddExpr(out)
	} else {
		outTok := I.ReadIdentifierDecl(false)
		if outTok != nil {
			out := glang.NewPropsBlockDecl(outTok)
			out.AddT(outTok)
			ret.Out = out
			ret.AddExprs(I.Emit())
			// its a non paren Out like func p(...) error {}
			// panic("unhandled so far")
		}
	}

	ret.AddExprs(I.Emit())
	return ret
}

func (I *GigoInterpreter) ReadFuncDecl(templated bool) *glang.FuncDecl {

	funcTok := I.MustRead(glanglexer.FuncToken)
	ret := glang.NewFuncDecl(funcTok)

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	receiver := I.ReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	if receiver != nil {
		ret.Receiver = receiver
		receiver.AddExprs(I.Emit())
		ret.AddExpr(receiver)
	}

	I.ReadMany(genericlexer.WsToken)
	ID := I.MustReadIdentifierDecl(templated)
	ret.Name = ID
	ret.AddExpr(ID)

	I.ReadMany(genericlexer.WsToken)
	params := I.MustReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	params.AddExprs(I.Emit())
	ret.Params = params
	ret.AddExpr(params)

	I.ReadMany(genericlexer.WsToken)
	bodyStart := I.Peek(glanglexer.BracketOpenToken)
	if bodyStart == nil {
		out := I.ReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		if out != nil {
			params.AddExprs(I.Emit())
			ret.Out = out
			ret.AddExpr(out)
		} else {
			outTok := I.ReadIdentifierDecl(templated)
			if outTok != nil {
				out := glang.NewPropsBlockDecl(outTok)
				out.AddT(outTok)
				out.AddExprs(I.Emit())
				ret.AddExpr(out)
				ret.Out = out
			}
		}
	}

	I.ReadMany(genericlexer.WsToken)
	body := I.ReadBodyBlock(glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
	if body != nil {
		body.AddExprs(I.Emit())
		ret.Body = body
		ret.AddExpr(body)
	} else {
		panic(I.Debug("Missing func body block: ", glanglexer.ParenOpenToken))
	}

	ret.AddExprs(I.Emit())
	return ret
}

func (I *GigoInterpreter) ReadStructDecl(templated bool) *glang.StructDecl {

	structTok := I.MustRead(glanglexer.StructToken)
	I.ReadMany(genericlexer.WsToken)

	ret := glang.NewStructDecl(structTok)
	ret.AddExprs(I.Emit())

	block := I.ReadPropsBlock(templated, glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
	ret.AddExpr(block)

	return ret
}

func (I *GigoInterpreter) ReadTemplateDecl() *glang.TemplateDecl {

	tplTok := I.MustRead(glanglexer.TemplateToken)
	I.ReadMany(genericlexer.WsToken)

	ret := glang.NewTemplateDecl(tplTok)
	ret.AddExprs(I.Emit())

	ID := I.MustReadIdentifierDecl(true)
	ret.Name = ID

	I.ReadMany(genericlexer.WsToken)
	if I.Peek(glanglexer.StructToken) != nil {
		structDecl := I.ReadStructDecl(true)
		structDecl.Name = ID
		structDecl.PrependExpr(ID)
		ret.AddExpr(structDecl)
	}

	return ret
}

func (I *GigoInterpreter) ReadVarDecl() *glang.VarDecl {

	varTok := I.MustRead(glanglexer.VarToken)
	I.ReadMany(genericlexer.WsToken)
	ret := glang.NewVarDecl(varTok)

	if I.Peek(glanglexer.ParenOpenToken) == nil {

		left := I.MustReadIdentifierDecl(false)
		I.ReadMany(genericlexer.WsToken)
		var leftType genericinterperter.Tokener
		if I.Peek(glanglexer.AssignToken) != nil {
			I.MustRead(glanglexer.AssignToken)
		} else {
			leftType = I.ReadIdentifierDecl(false)
			I.ReadMany(genericlexer.WsToken)
			I.MustRead(glanglexer.AssignToken)
		}
		I.ReadMany(genericlexer.WsToken)
		right := I.MustReadExpression()
		assignment := glang.NewAssignDecl(left)
		assignment.Left = left
		assignment.LeftType = leftType
		assignment.Right = right
		assignment.AddExpr(left)
		if leftType != nil {
			assignment.AddExpr(leftType)
		}
		if right != nil {
			assignment.AddExpr(right)
		}
		ret.AddExpr(assignment)

	} else {
		block := I.MustReadAssignsBlock(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		// block.AddExprs(I.Emit())
		ret.AddExpr(block)
	}

	return ret
}

func (I *GigoInterpreter) ReadConstDecl() *glang.ConstDecl {

	constTok := I.MustRead(glanglexer.ConstToken)
	I.ReadMany(genericlexer.WsToken)
	ret := glang.NewConstDecl(constTok)

	if I.Peek(glanglexer.ParenOpenToken) == nil {

		left := I.MustReadIdentifierDecl(false)
		I.ReadMany(genericlexer.WsToken)
		leftType := I.ReadIdentifierDecl(false)
		I.ReadMany(genericlexer.WsToken)

		var right *glang.ExpressionDecl
		if I.Read(glanglexer.AssignToken) != nil {
			right = I.MustReadExpression()
		}
		I.ReadMany(genericlexer.WsToken)

		assignment := glang.NewAssignDecl(left)
		assignment.Left = left
		assignment.LeftType = leftType
		assignment.Right = right
		assignment.AddExprs(I.Emit())
		ret.AddExpr(assignment)
	} else {
		ret.AddExprs(I.Emit())
		block := I.MustReadAssignsBlock(true, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		// block.AddExprs(I.Emit())
		ret.AddExpr(block)
	}

	return ret
}

func (I *GigoInterpreter) MustReadIdentifierDecl(templated bool) *glang.IdentifierDecl {
	r := I.ReadIdentifierDecl(templated)
	if r == nil {
		ts := []lexer.TokenType{genericlexer.WordToken}
		if templated {
			ts = append(ts, glanglexer.SmallerToken, glanglexer.GreaterToken)
		}
		panic(I.Debug("Identifier not found: ", ts...))
	}
	return r
}
func (I *GigoInterpreter) ReadIdentifierDecl(templated bool) *glang.IdentifierDecl {
	var ret *glang.IdentifierDecl
	var f genericinterperter.Tokener
	I.Read(glanglexer.ElipseToken)
	if templated == false {
		f = I.Read(genericlexer.WordToken)

	} else {
		for {
			if p := I.Read(genericlexer.WordToken); p != nil {
				// ok
				if f == nil {
					f = p
				}
			} else if p := I.Peek(glanglexer.SmallerToken); p != nil {
				//ok
				if f == nil {
					f = p
				}
				block := I.ReadBodyBlock(glanglexer.SmallerToken, glanglexer.GreaterToken)
				block.Open.SetType(glanglexer.TplOpenToken)
				block.Close.SetType(glanglexer.TplCloseToken)

			} else {
				break
			}
		}
	}
	if f != nil {
		ret = glang.NewIdentifierDecl(f)
		ret.AddExprs(I.Emit())
	}
	return ret
}
func (I *GigoInterpreter) MustReadExpression() *glang.ExpressionDecl {
	var f genericinterperter.Tokener
	for {
		if p := I.Read(glanglexer.NlToken); p != nil {
			break
		} else if p := I.Read(glanglexer.SemiColonToken); p != nil {
			break
		} else {
			k := I.Next()
			if f == nil {
				f = k
			}
		}
	}
	ret := glang.NewExpressionDecl(f)
	ret.AddExprs(I.Emit())
	return ret
}

func (I *GigoInterpreter) ReadImplDecl() *glang.ImplementDecl {

	implTok := I.MustRead(glanglexer.ImplementsToken)
	ret := glang.NewImplementDecl(implTok)

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	implTemplate := I.ReadBodyBlock(glanglexer.SmallerToken, glanglexer.GreaterToken)
	implTemplate.Open.SetType(glanglexer.TplOpenToken)
	implTemplate.Close.SetType(glanglexer.TplCloseToken)
	implTemplate.AddExprs(I.Emit())
	ret.AddExpr(implTemplate)
	ret.ImplementTemplate = implTemplate

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	block := I.ReadPropsBlock(true, glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
	ret.AddExpr(block)

	return ret
}

func (I *GigoInterpreter) ReadInterfaceDecl() *glang.InterfaceDecl {

	intfTok := I.MustRead(glanglexer.InterfaceToken)
	I.ReadMany(genericlexer.WsToken)

	ret := glang.NewInterfaceDecl(intfTok)
	ret.AddExprs(I.Emit())

	block := I.ReadSignsBlock(glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
	ret.AddExpr(block)

	return ret
}
