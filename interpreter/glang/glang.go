package glang

import (
	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
	genericlexer "github.com/mh-cbon/gigo/lexer/generic"
	glanglexer "github.com/mh-cbon/gigo/lexer/glang"
	glang "github.com/mh-cbon/gigo/struct/glang"
	lexer "github.com/mh-cbon/state-lexer"
)

// PackageProvider provides package.
type PackageProvider interface {
	AddToPackage(string, glang.ScopeReceiver)
}

// GigoInterpreter interprets gigo syntax
type GigoInterpreter struct {
	genericinterperter.Interpreter
	packages PackageProvider
}

// NewGigoInterpreter makes a new interpreter
func NewGigoInterpreter() *GigoInterpreter {
	return &GigoInterpreter{
		Interpreter: *genericinterperter.NewInterpreter(),
		packages:    &glang.SimplePackageRepository{},
	}
}

func (I *GigoInterpreter) getTokens(reader genericinterperter.TokenerReader) []genericinterperter.Tokener {
	var tokens []genericinterperter.Tokener
	for {
		if next := reader(); next != nil {
			tokens = append(tokens, next)
		} else {
			break
		}
	}
	return tokens
}

// ProcessFile processes given reader of tokens with the filepath as name.
func (I *GigoInterpreter) ProcessFile(file string, reader genericinterperter.TokenerReader) (*glang.FileDecl, error) {
	fileDef := &glang.FileDecl{Name: file}

	I.Tokens = I.getTokens(reader)
	I.Scope = fileDef

	return fileDef, I.Process(true)
}

// ProcessStr processes given reader of tokens of the content string.
func (I *GigoInterpreter) ProcessStr(content string, reader genericinterperter.TokenerReader) (*glang.StrDecl, error) {
	strDef := &glang.StrDecl{Src: content}

	I.Tokens = I.getTokens(reader)
	I.Scope = strDef

	return strDef, I.Process(false)
}

// ProcessStrWithPkgDecl processes given reader of tokens of the content string. It expects to have package decl.
func (I *GigoInterpreter) ProcessStrWithPkgDecl(content string, reader genericinterperter.TokenerReader) (*glang.StrDecl, error) {
	strDef := &glang.StrDecl{Src: content}

	I.Tokens = I.getTokens(reader)
	I.Scope = strDef

	return strDef, I.Process(true)
}

// Process given tokens in the given scope.
// a scope can be a file or a string.
func (I *GigoInterpreter) Process(withpkgdcl bool) error {

	pkgDecl, err := I.ReadPackageDecl()
	if withpkgdcl && err != nil {
		return err
	}
	if pkgDecl != nil {
		I.Scope.AddExpr(pkgDecl)
	}
	// I.packages.AddToPackage(decl.GetName(), I.Scope.(glang.ScopeReceiver))

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

			name, err := I.ReadIdentifierDecl(false)
			if err != nil {
				return err
			}

			I.ReadMany(genericlexer.WsToken)

			if typeTok := I.Peek(glanglexer.StructToken); typeTok != nil {

				sDecl, err := I.ReadStructDecl(false)
				if err != nil {
					return err
				}
				sDecl.Name = name
				sDecl.PrependExpr(name)
				sDecl.PrependExprs(preTokens)
				I.Scope.AddExpr(sDecl)

			} else if typeTok := I.Peek(glanglexer.InterfaceToken); typeTok != nil {

				sDecl, err := I.ReadInterfaceDecl()
				if err != nil {
					return err
				}
				sDecl.Name = name
				sDecl.PrependExpr(name)
				sDecl.PrependExprs(preTokens)
				I.Scope.AddExpr(sDecl)

			} else if tok := I.Peek(glanglexer.ImplementsToken); tok != nil {

				implDecl, err := I.ReadImplDecl()
				if err != nil {
					return err
				}
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

			tplDecl, err := I.ReadTemplateDecl()
			if err != nil {
				return err
			}
			I.Scope.AddExpr(tplDecl)

		} else if tok := I.Peek(glanglexer.TplOpenToken); tok != nil {

			// fn := &glang.TemplateFuncDecl{}
			fn := glang.NewTemplateFuncDecl()
			I.KeepPreviousComment()

			I.ReadMany(
				glanglexer.NlToken,
				genericlexer.CommentLineToken,
				genericlexer.CommentBlockToken,
				genericlexer.WsToken)
			fn.AddExprs(I.Emit())

			// smthig to improve here
			block, err := I.ReadBodyBlock(glanglexer.TplOpenToken, glanglexer.GreaterToken)
			if err != nil {
				return err
			}
			block.Open.SetType(glanglexer.TplOpenToken)
			block.Close.SetType(glanglexer.TplCloseToken)
			block.AddExprs(I.Emit())
			fn.Modifier = block
			fn.AddExpr(block)

			I.ReadMany(genericlexer.WsToken)
			fn.AddExprs(I.Emit())

			nFunc, err := I.ReadFuncDecl(true)
			if err != nil {
				return err
			}
			fn.Func = nFunc
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

			funcDecl, err := I.ReadFuncDecl(true)
			if err != nil {
				return err
			}
			I.Scope.AddExpr(funcDecl)

		} else if tok := I.Peek(glanglexer.VarToken); tok != nil {

			varDecl, err := I.ReadVarDecl()
			if err != nil {
				return err
			}
			I.Scope.AddExpr(varDecl)

		} else if tok := I.Peek(glanglexer.ConstToken); tok != nil {

			constDecl, err := I.ReadConstDecl()
			if err != nil {
				return err
			}
			I.Scope.AddExpr(constDecl)

		} else if x := I.Next(); x == nil {
			I.Scope.AddExprs(I.Emit())
			break
		}
	}
	return nil
}

// ReadPackageDecl reads the tokens until it finds a package token.
// returns an error if none is found.
// It creates a new package declaration, attach it to the scope,
// use its name to add/get a package in the current program
// and attach the current scope to it.
func (I *GigoInterpreter) ReadPackageDecl() (*glang.PackageDecl, error) {

	I.ReadMany(
		glanglexer.NlToken,
		genericlexer.CommentLineToken,
		genericlexer.CommentBlockToken,
		genericlexer.WsToken)

	tok := I.Read(glanglexer.PackageToken)
	if tok == nil {
		return nil, I.Debug("missing package decl", glanglexer.PackageToken)
	}

	decl := glang.NewPackageDecl()
	decl.AddExprs(I.Emit())
	decl.AddExprs(I.GetMany(genericlexer.WsToken))

	name := I.Get(genericlexer.WordToken)
	if name == nil {
		return nil, I.Debug("missing package name", genericlexer.WordToken)
	}
	decl.Name = name
	decl.AddExpr(name)
	decl.AddExprs(I.GetMany(genericlexer.WsToken, glanglexer.NlToken))

	return decl, nil
}

// KeepPreviousComment unreads the tokens appropriately (almost)
// to keep the comment attached to the declaration that is going to be analyzed.
// It does not Emit tokens.
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

// ReadPropsBlock reads a block of property.
// returns an error if none is found.
// The next token to analyze must be of type open,
// the block must end with a token of type close.
// In between data are read as a golang block of properties,
// type struct{ EmbedXXX; AProp string }.
func (I *GigoInterpreter) ReadPropsBlock(
	templated bool,
	open lexer.TokenType,
	close lexer.TokenType,
) (*glang.PropsBlockDecl, error) {

	var ret *glang.PropsBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return nil, I.Debug("unexpected token", open)
	}
	count := 1

	ret = glang.NewPropsBlockDecl()

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

			if I.Peek(glanglexer.PoireauPointerToken) != nil || I.Peek(glanglexer.PoireauToken) != nil {
				Poireau, err := I.ReadPoireauDecl()
				if err != nil {
					return nil, err
				}
				ret.AddPoireau(Poireau)

			} else {
				ID, err := I.ReadIdentifierDecl(templated)
				if err != nil {
					return nil, err
				}

				I.ReadMany(genericlexer.WsToken)
				IDType, err := I.ReadIdentifierDecl(templated)
				if IDType == nil {
					ret.AddUnderlying(ID)
				} else {
					if err != nil {
						return nil, err
					}
					ret.Add(ID, IDType)
				}
			}

		}
	}
	if I.Read(close) == nil {
		return nil, I.Debug("unexpected token", close)
	}
	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadAssignsBlock reads a block of asignment.
// returns an error if none is found.
// The next token to analyze must be of type open,
// the block must end with a token of type close.
// In between data are read as a golang block of properties,
// var (
//		xx="fff"
//		y="gggg"
//	)
// works with const also.
func (I *GigoInterpreter) ReadAssignsBlock(
	canOmitRight bool,
	open lexer.TokenType,
	close lexer.TokenType,
) (*glang.AssignsBlockDecl, error) {

	var ret *glang.AssignsBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return ret, I.Debug("unexpected token", open)
	}
	count := 1

	ret = glang.NewAssignsBlockDecl()

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

			left, err := I.ReadIdentifierDecl(false)
			if err != nil {
				return nil, err
			}
			assignment := glang.NewAssignDecl()
			assignment.AddExpr(left)

			I.ReadMany(genericlexer.WsToken)
			assignment.AddExprs(I.Emit())

			leftType, err := I.ReadIdentifierDecl(false)
			if leftType != nil {
				if err != nil {
					return ret, err
				}
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

					right, err = I.ReadExpression()
					if err != nil {
						return nil, err
					}
					assignment.Right = right
					assignment.AddExpr(right)
				}
			} else {
				if I.Read(glanglexer.AssignToken) == nil {
					return ret, I.Debug("unexpected token", glanglexer.AssignToken)
				}
				I.ReadMany(genericlexer.WsToken)
				right, err = I.ReadExpression()
				if err != nil {
					return nil, err
				}
				assignment.Right = right
				assignment.AddExpr(right)
			}

			assignment.AddExprs(I.Emit())
			ret.AddExpr(assignment)

		}
	}
	if I.Read(close) == nil {
		return ret, I.Debug("unexpected token", close)
	}
	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadSignsBlock reads a block of func signatures.
// returns an error if none is found.
// The next token to analyze must be of type open,
// the block must end with a token of type close.
// In between data are read as a golang block of properties,
// type interface { f() }
func (I *GigoInterpreter) ReadSignsBlock(
	open lexer.TokenType,
	close lexer.TokenType,
) (*glang.SignsBlockDecl, error) {

	var ret *glang.SignsBlockDecl
	I.ReadMany(
		glanglexer.NlToken,
		genericlexer.WsToken,
		genericlexer.CommentBlockToken,
		genericlexer.CommentLineToken,
	)

	openTok := I.Read(open)
	if openTok == nil {
		return ret, I.Debug("unexpected token", open)
	}
	count := 1

	ret = glang.NewSignsBlockDecl()

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

			if ID, err := I.ReadIdentifierDecl(false); ID != nil {
				if err != nil {
					return ret, err
				}
				I.ReadMany(genericlexer.WsToken)
				if I.Peek(glanglexer.NlToken) != nil {
					ret.AddUnderlying(ID)
					I.ReadMany(glanglexer.NlToken)
					ret.AddExprs(I.Emit())
				} else {
					sign, err := I.ReadFuncSign(ID)
					if err != nil {
						return nil, err
					}
					ret.Add(sign)
				}
			}

		}
	}
	if I.Read(close) == nil {
		return ret, I.Debug("unexpected token", close)
	}
	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadParenDecl reads a block of parameters declaration.
// returns an error if none is found.
// The next token to analyze must be of type open,
// the block must end with a token of type close.
// In between data are read as a golang block of properties,
// func (r *Recevier) ...(p string, n int) (bool, err)
func (I *GigoInterpreter) ReadParenDecl(
	templated bool, open,
	close lexer.TokenType,
) (*glang.PropsBlockDecl, error) {

	var ret *glang.PropsBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return nil, I.Debug("unexpected token", open)
	}
	count := 1

	ret = glang.NewPropsBlockDecl()

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

			ID, err := I.ReadIdentifierDecl(templated)
			if err != nil {
				return nil, err
			}

			I.ReadMany(genericlexer.WsToken, genericlexer.CommentBlockToken)

			IDType, err := I.ReadIdentifierDecl(templated)

			var propdecl *glang.PropDecl
			if IDType == nil {
				propdecl = ret.AddT(ID)
			} else {
				if err != nil {
					return ret, err
				}
				propdecl = ret.Add(ID, IDType)
			}
			propdecl.AddExprs(I.Emit())

			I.ReadMany(genericlexer.WsToken, genericlexer.CommentBlockToken, glanglexer.SemiColonToken, glanglexer.NlToken)
		}
	}
	if I.Read(close) == nil {
		return nil, I.Debug("unexpected token", close)
	}
	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadBodyBlock reads a block.
// returns an error if none is found.
// The next token to analyze must be of type open,
// the block must end with a token of type close.
// In between data is accumulated, any block is ok.
func (I *GigoInterpreter) ReadBodyBlock(open lexer.TokenType, close lexer.TokenType) (*glang.BodyBlockDecl, error) {

	var ret *glang.BodyBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return nil, I.Debug("unexpected token", open)
	}
	count := 1
	ret = glang.NewBodyBlockDecl()
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
	return ret, nil
}

// ReadFuncSign reads a func signature.
// returns an error if none is found.
// AddExpr(expr Tokener)
func (I *GigoInterpreter) ReadFuncSign(ID *glang.IdentifierDecl) (*glang.FuncDecl, error) {
	ret := glang.NewFuncDecl()
	ret.Name = ID
	ret.AddExpr(ID)

	I.ReadMany(genericlexer.WsToken)
	params, err := I.ReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	if err != nil {
		return nil, err
	}
	params.AddExprs(I.Emit())
	ret.Params = params
	ret.AddExpr(params)

	I.ReadMany(genericlexer.WsToken)
	out, _ := I.ReadParenDecl(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	if out != nil {
		out.AddExprs(I.Emit())
		ret.Out = out
		ret.AddExpr(out)
	} else {
		outTok, err := I.ReadIdentifierDecl(false)
		if outTok != nil {
			if err != nil {
				return ret, err
			}
			out := glang.NewPropsBlockDecl()
			out.AddT(outTok)
			ret.Out = out
			ret.AddExprs(I.Emit())
			// its a non paren Out like func p(...) error {}
			// panic("unhandled so far")
		}
	}

	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadFuncDecl reads a func with or without receiver.
// the next token must be a funcToken
// returns an error if none is found.
// func (r) AddExpr(expr Tokener) out { body }
func (I *GigoInterpreter) ReadFuncDecl(templated bool) (*glang.FuncDecl, error) {

	funcTok := I.Read(glanglexer.FuncToken)
	if funcTok == nil {
		return nil, I.Debug("Unexpected token", glanglexer.FuncToken)
	}
	ret := glang.NewFuncDecl()

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	receiver, err := I.ReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	if receiver != nil {
		if err != nil {
			return nil, err
		}
		ret.Receiver = receiver
		receiver.AddExprs(I.Emit())
		ret.AddExpr(receiver)
	}

	I.ReadMany(genericlexer.WsToken)
	ID, _ := I.ReadIdentifierDecl(templated)
	ret.Name = ID
	ret.AddExpr(ID)

	I.ReadMany(genericlexer.WsToken)
	params, err := I.ReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	if err != nil {
		return nil, err
	}
	params.AddExprs(I.Emit())
	ret.Params = params
	ret.AddExpr(params)

	I.ReadMany(genericlexer.WsToken)
	bodyStart := I.Peek(glanglexer.BracketOpenToken)
	if bodyStart == nil {
		out, _ := I.ReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		if out != nil {
			params.AddExprs(I.Emit())
			ret.Out = out
			ret.AddExpr(out)
		} else {
			outTok, _ := I.ReadIdentifierDecl(templated)
			if outTok != nil {
				out := glang.NewPropsBlockDecl()
				out.AddT(outTok)
				out.AddExprs(I.Emit())
				ret.AddExpr(out)
				ret.Out = out
			}
		}
	}

	I.ReadMany(genericlexer.WsToken)
	body, err := I.ReadBodyBlock(glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
	if err != nil {
		return ret, err
	}
	body.AddExprs(I.Emit())
	ret.Body = body
	ret.AddExpr(body)

	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadStructDecl reads a struct with its props.
// the next token must be a StructToken
// returns an error if none is found.
// type xx struct { block }
func (I *GigoInterpreter) ReadStructDecl(templated bool) (*glang.StructDecl, error) {

	structTok := I.Read(glanglexer.StructToken)
	if structTok == nil {
		return nil, I.Debug("unexpected token", glanglexer.StructToken)
	}
	I.ReadMany(genericlexer.WsToken)
	ret := glang.NewStructDecl()
	ret.AddExprs(I.Emit())

	block, err := I.ReadPropsBlock(true, glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
	if block != nil {
		if err != nil {
			return nil, err
		}
		ret.AddExpr(block)
		ret.Block = block
	}
	return ret, err
}

// ReadTemplateDecl reads a template with its props.
// the next token must be a TemplateToken
// returns an error if none is found.
// template xx<..> struct { block }
func (I *GigoInterpreter) ReadTemplateDecl() (*glang.TemplateDecl, error) {

	tplTok := I.Read(glanglexer.TemplateToken)
	if tplTok == nil {
		return nil, I.Debug("unexpected token", glanglexer.TemplateToken)
	}
	I.ReadMany(genericlexer.WsToken)

	ret := glang.NewTemplateDecl()
	ret.AddExprs(I.Emit())

	ID, err := I.ReadIdentifierDecl(true)
	if err != nil {
		return nil, err
	}
	ret.Name = ID

	I.ReadMany(genericlexer.WsToken)
	if I.Peek(glanglexer.StructToken) != nil {
		structDecl, err := I.ReadStructDecl(true)
		if err != nil {
			return ret, nil
		}
		structDecl.Name = ID
		structDecl.PrependExpr(ID)
		ret.AddExpr(structDecl)
		ret.Block = structDecl.Block
		ret.Methods = structDecl.Methods
	}

	return ret, nil
}

// ReadVarDecl reads a var declaration.
// the next token must be a VarToken
// returns an error if none is found.
// var (
// 		x = ""
// )
func (I *GigoInterpreter) ReadVarDecl() (*glang.VarDecl, error) {

	varTok := I.Read(glanglexer.VarToken)
	if varTok == nil {
		return nil, I.Debug("unexpected token", glanglexer.VarToken)
	}
	I.ReadMany(genericlexer.WsToken)
	ret := glang.NewVarDecl()

	if I.Peek(glanglexer.ParenOpenToken) == nil {

		left, err2 := I.ReadIdentifierDecl(false)
		if err2 != nil {
			return nil, err2
		}
		I.ReadMany(genericlexer.WsToken)
		var err error
		var leftType genericinterperter.Tokener
		if I.Peek(glanglexer.AssignToken) != nil {
			if I.Read(glanglexer.AssignToken) == nil {
				return nil, I.Debug("unexpected token", glanglexer.AssignToken)
			}
		} else {
			leftType, err = I.ReadIdentifierDecl(false)
			if leftType != nil {
				if err != nil {
					return ret, err
				}
			}
			I.ReadMany(genericlexer.WsToken)
			if I.Read(glanglexer.AssignToken) == nil {
				return nil, I.Debug("unexpected token", glanglexer.AssignToken)
			}
		}
		I.ReadMany(genericlexer.WsToken)
		right, err := I.ReadExpression()
		if err != nil {
			return nil, err
		}
		assignment := glang.NewAssignDecl()
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
		block, err := I.ReadAssignsBlock(false, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		if err != nil {
			return nil, err
		}
		// block.AddExprs(I.Emit())
		ret.AddExpr(block)
	}

	return ret, nil
}

// ReadConstDecl reads a const declaration.
// the next token must be a ConstToken
// returns an error if none is found.
// const (
// 		x = ""
// )
func (I *GigoInterpreter) ReadConstDecl() (*glang.ConstDecl, error) {

	constTok := I.Read(glanglexer.ConstToken)
	if constTok == nil {
		return nil, I.Debug("unexpected token", glanglexer.ConstToken)
	}
	I.ReadMany(genericlexer.WsToken)
	ret := glang.NewConstDecl()

	if I.Peek(glanglexer.ParenOpenToken) == nil {

		left, err2 := I.ReadIdentifierDecl(false)
		if err2 != nil {
			return nil, err2
		}
		I.ReadMany(genericlexer.WsToken)
		leftType, err := I.ReadIdentifierDecl(false)
		if leftType != nil && err != nil {
			return ret, err
		}
		I.ReadMany(genericlexer.WsToken)

		var right *glang.ExpressionDecl
		if I.Read(glanglexer.AssignToken) != nil {
			right, err = I.ReadExpression()
			if err != nil {
				return nil, err
			}
		}
		I.ReadMany(genericlexer.WsToken)

		assignment := glang.NewAssignDecl()
		assignment.Left = left
		assignment.LeftType = leftType
		assignment.Right = right
		assignment.AddExprs(I.Emit())
		ret.AddExpr(assignment)
	} else {
		ret.AddExprs(I.Emit())
		block, err := I.ReadAssignsBlock(true, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		if err != nil {
			return nil, err
		}
		// block.AddExprs(I.Emit())
		ret.AddExpr(block)
	}

	return ret, nil
}

// ReadIdentifierDecl reads an identifier.
// the next token must be a WordToken/TplOpenToken.
// returns an error if none is found.
// wordToken: nameOfTheFunc / nameofTheStruct / nameofTHeParam ect
//		it must be a word.
// template: <...>name / name<...> / name<...>name
//		it must be one to many subsequent wordToken / template.
func (I *GigoInterpreter) ReadIdentifierDecl(templated bool) (*glang.IdentifierDecl, error) {
	var ret *glang.IdentifierDecl
	var f genericinterperter.Tokener
	I.Read(glanglexer.ElipseToken)
	if templated == false {
		f = I.Read(genericlexer.WordToken)
		if f == nil {
			return nil, I.Debug("unexpected token", genericlexer.WordToken)
		}
	} else {
		for {
			if p := I.Read(genericlexer.WordToken); p != nil {
				// ok
				if f == nil {
					f = p
				}
			} else if p := I.Peek(glanglexer.TplOpenToken); p != nil {
				//ok
				if f == nil {
					f = p
				}
				block, err := I.ReadBodyBlock(glanglexer.TplOpenToken, glanglexer.GreaterToken)
				if err != nil {
					return nil, err
				}
				block.AddExprs(I.Emit())
				block.Open.SetType(glanglexer.TplOpenToken)
				block.Close.SetType(glanglexer.TplCloseToken)
				if ret == nil {
					ret = glang.NewIdentifierDecl()
				}
				ret.AddExpr(block)

			} else {
				break
			}
		}
	}
	if f != nil {
		if ret == nil {
			ret = glang.NewIdentifierDecl()
		}
		ret.AddExprs(I.Emit())
	} else {
		return nil, I.Debug("unexpected token", genericlexer.WordToken, glanglexer.TplOpenToken)
	}
	return ret, nil
}

// ReadPoireauDecl reads a poireau<M()> declaration.
// the next token must be a PoireauToken | PoireauPointerToken.
// returns an error if none is found.
// PoireauToken:
//	- poireau<Mutator>
//	- *poireau<Mutator>
func (I *GigoInterpreter) ReadPoireauDecl() (*glang.PoireauDecl, error) {
	var ret *glang.PoireauDecl
	tok := I.Read(glanglexer.PoireauToken, glanglexer.PoireauPointerToken)
	if tok == nil {
		return nil, I.Debug("unexpected token", glanglexer.PoireauToken, glanglexer.PoireauPointerToken)
	}

	ID, err := I.ReadIdentifierDecl(true)
	if err != nil {
		return nil, err
	}
	ret = glang.NewPoireauDecl()
	ret.ImplementTemplate = ID
	ret.AddExpr(ID)
	ret.AddExprs(I.Emit())

	return ret, nil
}

//ReadExpression reads tokens until Nl / semmi if found.
// returns an error if non is found.
func (I *GigoInterpreter) ReadExpression() (*glang.ExpressionDecl, error) {
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
	ret := glang.NewExpressionDecl()
	ret.AddExprs(I.Emit())
	if len(ret.GetTokens()) == 0 {
		return nil, I.Debug("Unexpected token", glanglexer.NlToken, glanglexer.SemiColonToken)
	}
	return ret, nil
}

// ReadImplDecl reads an implements with its props.
// the next token must be a ImplementsToken
// returns an error if none is found.
// type xx implements<..> { block }
func (I *GigoInterpreter) ReadImplDecl() (*glang.ImplementDecl, error) {

	implTok := I.Read(glanglexer.ImplementsToken)
	if implTok == nil {
		return nil, I.Debug("unexpected token", glanglexer.ImplementsToken)
	}
	ret := glang.NewImplementDecl()

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	implTemplate, err := I.ReadBodyBlock(glanglexer.TplOpenToken, glanglexer.GreaterToken)
	if err != nil {
		return nil, err
	}
	implTemplate.Open.SetType(glanglexer.TplOpenToken)
	implTemplate.Close.SetType(glanglexer.TplCloseToken)
	implTemplate.AddExprs(I.Emit())
	ret.AddExpr(implTemplate)
	ret.ImplementTemplate = implTemplate

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	block, err := I.ReadPropsBlock(true, glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
	if err != nil {
		return nil, err
	}
	ret.AddExpr(block)

	return ret, nil
}

// ReadInterfaceDecl reads an interface with its signs.
// the next token must be a InterfaceToken
// returns an error if none is found.
// type xx interface { block }
func (I *GigoInterpreter) ReadInterfaceDecl() (*glang.InterfaceDecl, error) {

	intfTok := I.Read(glanglexer.InterfaceToken)
	if intfTok == nil {
		return nil, I.Debug("unexpected token", glanglexer.InterfaceToken)
	}
	I.ReadMany(genericlexer.WsToken)

	ret := glang.NewInterfaceDecl()
	ret.AddExprs(I.Emit())

	block, err := I.ReadSignsBlock(glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
	if err != nil {
		return nil, err
	}
	ret.AddExpr(block)
	ret.Block = block

	return ret, nil
}
