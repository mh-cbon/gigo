package glang

import (
	"fmt"

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
	packages   PackageProvider
	blockscope *Scope
}

// NewGigoInterpreter makes a new interpreter
func NewGigoInterpreter(r genericinterperter.TokenerReaderOK) *GigoInterpreter {
	return &GigoInterpreter{
		Interpreter: *genericinterperter.NewInterpreter(r),
		packages:    &glang.SimplePackageRepository{},
		blockscope:  NewScope(),
	}
}

// ProcessFile processes given reader of tokens with the filepath as name.
func (I *GigoInterpreter) ProcessFile(file string) (*glang.FileDecl, error) {
	fileDef := &glang.FileDecl{Name: file}

	I.Scope = fileDef

	return fileDef, I.Process(true)
}

// ProcessStr processes given reader of tokens of the content string.
func (I *GigoInterpreter) ProcessStr(content string) (*glang.StrDecl, error) {
	strDef := &glang.StrDecl{Src: content}

	I.Scope = strDef

	return strDef, I.Process(false)
}

// ProcessStrWithPkgDecl processes given reader of tokens of the content string. It expects to have package decl.
func (I *GigoInterpreter) ProcessStrWithPkgDecl(content string) (*glang.StrDecl, error) {
	strDef := &glang.StrDecl{Src: content}

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

			name, err := I.ReadVarName(false, false, false)
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
			block, err := I.ReadTemplateExprDecl()
			if err != nil {
				return err
			}
			I.Scope.AddExpr(block)

		} else if tok := I.Peek(glanglexer.FuncToken); tok != nil {

			I.KeepPreviousComment()
			I.Scope.AddExprs(I.Emit())

			I.ReadMany(
				glanglexer.NlToken,
				genericlexer.CommentLineToken,
				genericlexer.CommentBlockToken,
				genericlexer.WsToken)

			funcDecl, err := I.ReadFuncDecl(true, false)
			if err != nil {
				return err
			}
			I.Scope.AddExpr(funcDecl)

		} else if tok := I.Peek(glanglexer.VarToken); tok != nil {

			varDecl, err := I.ReadVarDecl(true)
			if err != nil {
				return err
			}
			I.Scope.AddExpr(varDecl)

		} else if tok := I.Peek(glanglexer.ConstToken); tok != nil {

			constDecl, err := I.ReadConstDecl(true)
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
				ID, err := I.ReadVarName(templated, false, true)
				if err != nil {
					x, err2 := I.ReadTypeName(templated, false)
					if err2 == nil && x != nil {
						ret.AddUnderlying(x)
						ret.AddExpr(x)
						continue
					}
					return nil, err
				}

				ws := I.GetMany(genericlexer.WsToken)

				IDType, err := I.ReadTypeName(templated, true)
				if IDType == nil {
					x := glang.NewExpressionDecl()
					x.AddExpr(ID)
					ret.AddUnderlying(x)
					ret.AddExpr(x)
					ret.AddExprs(ws)
				} else {
					if err != nil {
						return nil, err
					}
					ret.Add(ID, IDType)
					ret.AddExpr(ID)
					ret.AddExprs(ws)
					ret.AddExpr(IDType)
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
		return nil, I.Debug("unexpected token", open)
	}
	count := 1

	ret = glang.NewAssignsBlockDecl()
	ret.AddExprs(I.Emit())

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

			assignment := glang.NewAssignDecl()
			left, err := I.ReadVarName(false, true, false)
			if err != nil {
				return nil, err
			}
			assignment.Left = left
			assignment.AddExpr(left)

			I.ReadMany(genericlexer.WsToken)
			assignment.AddExprs(I.Emit())

			if I.Peek(glanglexer.NlToken) == nil {

				eq := I.Read(glanglexer.AssignToken)
				if eq == nil {
					leftType, err2 := I.ReadTypeName(false, true)
					if err2 != nil {
						return nil, err2
					}
					if leftType != nil {
						assignment.LeftType = leftType
						assignment.AddExpr(leftType)
					}
					I.ReadMany(genericlexer.WsToken)

					eq = I.Read(glanglexer.AssignToken)
				}

				if !canOmitRight && eq == nil {
					return nil, I.Debug("unexpected token", glanglexer.AssignToken)
				}

				if eq == nil {
					I.ReadMany(genericlexer.WsToken, glanglexer.NlToken)
					assignment.AddExprs(I.Emit())

				} else {

					assignment.Assign = eq
					I.ReadMany(genericlexer.WsToken)
					assignment.AddExprs(I.Emit())

					var right *glang.ExpressionDecl
					right, err = I.ReadExpressionBlock(true, glanglexer.SemiColonToken)
					if err != nil {
						return nil, err
					}
					if right != nil {
						assignment.Right = right
						assignment.AddExpr(right)
					}
				}
			}

			assignment.AddExprs(I.Emit())
			ret.Assigns = append(ret.Assigns, assignment)
			ret.AddExpr(assignment)
		}
	}
	if I.Read(close) == nil {
		return nil, I.Debug("unexpected token", close)
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

			if ID, err := I.ReadVarName(false, false, false); ID != nil {
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
		return nil, I.Debug("Failed to read ReadParenDecl", open)
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

		} else if I.Read(glanglexer.CommaToken) != nil {

		} else {
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())

			I.ReadMany(glanglexer.NlToken, genericlexer.CommentBlockToken)

			ID, err := I.ReadVarName(templated, false, false)
			if err != nil {
				return nil, err
			}

			I.ReadMany(genericlexer.WsToken, genericlexer.CommentBlockToken)

			IDType, err := I.ReadTypeName(templated, true)

			var propdecl *glang.PropDecl
			if IDType == nil {
				x := glang.NewExpressionDecl()
				x.AddExpr(ID)
				propdecl = ret.AddT(x)
				ret.AddExpr(propdecl)
			} else {
				if err != nil {
					return ret, err
				}
				propdecl = ret.Add(ID, IDType)
				ret.AddExpr(propdecl)
			}
			propdecl.AddExprs(I.Emit())

			I.ReadMany(genericlexer.WsToken,
				genericlexer.CommentBlockToken,
				glanglexer.CommaToken,
				glanglexer.NlToken)
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
			r := I.Next()
			if r == nil {
				break
			}
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
		outTok, err := I.ReadTypeIdentifier(false)
		if err != nil {
			return nil, err
		}
		if outTok != nil {
			out := glang.NewPropsBlockDecl()
			out.AddT(outTok)
			out.AddExpr(outTok)
			ret.Out = out
			ret.AddExprs(I.Emit())
			// its a non paren Out like func p(...) error {}
			// panic("unhandled so far")
		}
	}

	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadExpressionsBlock ...
func (I *GigoInterpreter) ReadExpressionsBlock(templated bool, open lexer.TokenType, close lexer.TokenType) (*glang.BodyBlockDecl, error) {

	var ret *glang.BodyBlockDecl

	openTok := I.Read(open)
	if openTok == nil {
		return nil, I.Debug("require token", open)
	}

	count := 1
	ret = glang.NewBodyBlockDecl()
	ret.Open = openTok

	I.blockscope.Enter()
	defer I.blockscope.Leave()

	for {
		I.ReadMany(
			genericlexer.WsToken,
			glanglexer.NlToken,
			genericlexer.CommentLineToken,
			genericlexer.CommentBlockToken,
		)
		ret.AddExprs(I.Emit())

		if openTok := I.Read(open); openTok != nil {
			count++

		} else if closeTok := I.Read(close); closeTok != nil {
			count--
			if count == 0 {
				ret.Close = closeTok
				break
			}

		} else {

			if I.Peek(glanglexer.VarToken) != nil {
				varDecl, err := I.ReadVarDecl(templated)
				if err != nil {
					return nil, err
				}
				ret.AddExpr(varDecl)
			} else if I.Peek(glanglexer.ReturnToken) != nil {
				r, err := I.ReadReturnDecl(templated)
				if err != nil {
					return nil, err
				}
				ret.AddExpr(r)

			} else if I.Peek(genericlexer.WordToken) != nil {
				expr, err := I.ReadExpressionBlock(templated, glanglexer.NlToken, close)
				if err != nil {
					return nil, err
				}
				ret.AddExpr(expr)

			} else if I.Peek(glanglexer.ForToken) != nil {
				expr, err := I.ReadForBlock(templated)
				if err != nil {
					return nil, err
				}
				ret.AddExpr(expr)

			} else if I.Peek(glanglexer.IfToken) != nil {
				expr, err := I.ReadIfStmt(templated)
				if err != nil {
					return nil, err
				}
				ret.AddExpr(expr)

			} else if I.Next() == nil {
				return nil, I.Debug("require token", close)
				break
			}
			// <-time.After(time.Second * 2)

		}
	}
	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadVarName ...
func (I *GigoInterpreter) ReadVarName(templated, allowunderscore, allowdot bool) (*glang.IdentifierDecl, error) {
	var ret *glang.IdentifierDecl

	if p := I.Peek(genericlexer.WordToken); p != nil {
		v := p.GetValue()
		f := v[0]
		if !((f >= 'a' && f <= 'z') || (f >= 'A' && f <= 'Z') || allowunderscore && v == "_") {
			return nil, I.Debug("Invalid value '"+v+"', must start with char", genericlexer.WordToken)
		}
	} else if templated {
		if I.Peek(glanglexer.TplOpenToken) == nil {
			return nil, I.Debug("unexpected token", genericlexer.WordToken, glanglexer.TplOpenToken)
		}
	} else {
		return nil, I.Debug("unexpected token", genericlexer.WordToken, glanglexer.TplOpenToken)
	}

	ret = glang.NewIdentifierDecl()
	for {
		if p := I.Read(genericlexer.WordToken); p != nil {
			continue
		} else if templated && I.Peek(glanglexer.TplOpenToken) != nil {
			ret.AddExprs(I.Emit())
			block, err := I.ReadTemplateBlock()
			if err != nil {
				return nil, err
			}
			ret.AddExpr(block)
			continue
		} else if allowdot {
			if I.Read(glanglexer.DotToken) == nil {
				break
			}
		} else {
			break
		}
	}
	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadNumber ...
func (I *GigoInterpreter) ReadNumber() (*glang.IdentifierDecl, error) {
	var ret *glang.IdentifierDecl

	first := I.Peek(genericlexer.WordToken)
	if first == nil {
		return nil, I.Debug("unexpected token", genericlexer.WordToken)
	}
	v := first.GetValue()
	f := v[0]
	if !(f >= '0' && f <= '9') {
		return nil, I.Debug("Invalid value '"+v+"', must start with number", genericlexer.WordToken)
	}

	ret = glang.NewIdentifierDecl()
	for {
		if p := I.Read(genericlexer.WordToken); p != nil {
			v = p.GetValue()
			f = v[0]
			if !(f >= '0' && f <= '9') && f != 'e' {
				I.RewindAll()
				return nil, I.Debug("Invalid value '"+v+"'", genericlexer.WordToken)
			}
			continue
		} else if p := I.Read(glanglexer.DotToken); p != nil {
			continue
		} else {
			break
		}
	}
	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadTypeIdentifier ...
func (I *GigoInterpreter) ReadTypeIdentifier(templated bool) (*glang.ExpressionDecl, error) {
	var ret *glang.ExpressionDecl

	if p := I.Read(
		glanglexer.StringToken,
		glanglexer.IntToken,
		glanglexer.Int8Token,
		glanglexer.Int16Token,
		glanglexer.Int32Token,
		glanglexer.Int64Token,
		glanglexer.UintToken,
		glanglexer.Uint8Token,
		glanglexer.Uint16Token,
		glanglexer.Uint32Token,
		glanglexer.Uint64Token,
		glanglexer.FloatToken,
		glanglexer.Float32Token,
		glanglexer.Float64Token,
	); p != nil {
		ret = glang.NewExpressionDecl()
		ID := glang.NewIdentifierDecl()
		ID.AddExprs(I.Emit())
		ret.AddExpr(ID)

	} else if p := I.Read(glanglexer.InterfaceToken); p != nil {
		ret = glang.NewExpressionDecl()
		ID := glang.NewIdentifierDecl()
		I.ReadMany(genericlexer.WsToken, glanglexer.NlToken)
		I.ReadMany(glanglexer.BraceOpenToken)
		I.ReadMany(genericlexer.WsToken, glanglexer.NlToken)
		I.ReadMany(glanglexer.BraceCloseToken)
		ID.AddExprs(I.Emit())
		ret.AddExpr(ID)

	} else if p := I.Peek(glanglexer.StructToken); p != nil {
		ret = glang.NewExpressionDecl()
		block, err := I.ReadStructDecl(templated)
		if err != nil {
			return nil, err
		}
		ret.AddExpr(block)
	}

	return ret, nil
}

// ReadIdent ...
func (I *GigoInterpreter) ReadIdent(templated bool, allowunderscore bool) (*glang.ExpressionDecl, error) {
	var ret *glang.ExpressionDecl

	x, err := I.ReadVarName(templated, allowunderscore, true)
	if err != nil {
		return nil, err
	}

	ret = glang.NewExpressionDecl()
	ret.AddExpr(x)

	return ret, nil
}

// ReadTypeName ...
func (I *GigoInterpreter) ReadTypeName(templated bool, brackets bool) (*glang.ExpressionDecl, error) {

	var ret *glang.ExpressionDecl

	if brackets && I.Peek(glanglexer.BracketOpenToken) != nil {
		for {
			I.Read(glanglexer.BracketOpenToken)
			I.ReadMany(genericlexer.WsToken)
			I.Read(glanglexer.BracketCloseToken)
			I.ReadMany(genericlexer.WsToken)
			if I.Peek(glanglexer.BracketOpenToken) == nil {
				break
			}
		}
	}
	I.Read(glanglexer.MulToken)
	I.ReadMany(genericlexer.WsToken)

	name, err := I.ReadTypeIdentifier(templated)

	if name == nil && err == nil {
		x, err2 := I.ReadIdent(templated, false)
		err = err2
		if x != nil {
			ret = x
			x.PrependExprs(I.Emit())
		}
	} else if name != nil {
		ret = name
		name.PrependExprs(I.Emit())
	}

	if err != nil {
		I.RewindAll()
		return nil, err
	}

	return ret, nil
}

// ReadTypeValue ...
func (I *GigoInterpreter) ReadTypeValue(templated bool) (*glang.ExpressionDecl, error) {
	var ret *glang.ExpressionDecl

	I.Read(glanglexer.AndToken)

	doBraces := false

	ret = glang.NewExpressionDecl()
	ID, err := I.ReadTypeName(templated, true)
	if err != nil {
		return nil, err
	}
	if ID != nil {
		ret = ID
		// ret.AddExpr(ID)

		if x, ok := ID.First().(*glang.IdentifierDecl); ok {
			doBraces = !I.blockscope.HasVar(x.GetVarName())
		} else if _, ok := ID.First().(*glang.StructDecl); ok {
			doBraces = true
		} else {
			// maight need some check, not sure what will go here.
		}
	}

	I.ReadMany(genericlexer.WsToken)
	if doBraces /* && I.Peek(glanglexer.BraceOpenToken) != nil*/ {
		I.ReadBlock(glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
		ret.AddExprs(I.Emit())
	} else if I.Peek(glanglexer.BracketOpenToken) != nil {
		I.ReadBlock(glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
		ret.AddExprs(I.Emit())
	} else if I.Peek(glanglexer.IncToken, glanglexer.DecToken) != nil {
		I.Read(glanglexer.IncToken, glanglexer.DecToken)
		ret.AddExprs(I.Emit())
	} else {
		I.RewindAll()
	}

	return ret, nil
}

// ReadIfStmt ...
func (I *GigoInterpreter) ReadIfStmt(templated bool) (*glang.IfStmt, error) {
	var ret *glang.IfStmt

	if I.Read(glanglexer.IfToken) == nil {
		return nil, I.Debug("unexpected token", glanglexer.IfToken)
	}

	ret = glang.NewIfStmt()

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	whatis := I.PeekUntil(glanglexer.SemiColonToken,
		glanglexer.TypeAssignToken, glanglexer.BraceOpenToken)
	I.RewindAll()

	if whatis.GetType() == glanglexer.TypeAssignToken {
		I.RewindAll()
		init, err := I.ReadAssignExpr(templated, false, glanglexer.SemiColonToken)
		if err != nil {
			return nil, err
		}
		ret.Init = init
		ret.AddExpr(init)
		I.ReadMany(genericlexer.WsToken)
		I.Read(glanglexer.SemiColonToken)
		ret.AddExprs(I.Emit())
		I.blockscope.AddVar(init.CollectVarNames()...)

	} else if whatis.GetType() == glanglexer.SemiColonToken {
		// weird
		return nil, I.Debug("Not an assignment")
	}

	I.Read(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	block, err := I.ReadBinaryExpressionBlock(templated, glanglexer.BraceOpenToken)
	if err != nil {
		return nil, err
	}
	ret.AddExpr(block)
	ret.Cond = block

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	body, err := I.ReadExpressionsBlock(templated,
		glanglexer.BraceOpenToken,
		glanglexer.BraceCloseToken)
	if err != nil {
		return nil, err
	}
	ret.Body = body
	ret.AddExpr(body)

	I.ReadMany(
		genericlexer.WsToken,
		genericlexer.CommentBlockToken,
		genericlexer.CommentLineToken,
	)
	ret.AddExprs(I.Emit())

	if I.Read(glanglexer.ElseToken) != nil {
		block, err := I.ReadElseStmt(templated)
		if err != nil {
			return nil, err
		}
		ret.AddExpr(block)
	}

	return ret, nil
}

// ReadElseStmt ...
func (I *GigoInterpreter) ReadElseStmt(templated bool) (*glang.ElseStmt, error) {
	var ret *glang.ElseStmt

	if I.Read(glanglexer.ElseToken) == nil {
		return nil, I.Debug("unexpected token", glanglexer.ElseToken)
	}

	ret = glang.NewElseStmt()

	I.ReadMany(genericlexer.WsToken, glanglexer.NlToken)

	if I.Peek(glanglexer.IfToken) != nil {
		block, err := I.ReadIfStmt(templated)
		if err != nil {
			return nil, err
		}
		ret = &glang.ElseStmt{IfStmt: *block}
	} else {
		ret.AddExprs(I.Emit())
		body, err := I.ReadExpressionsBlock(templated,
			glanglexer.BraceOpenToken,
			glanglexer.BraceCloseToken)
		if err != nil {
			return nil, err
		}
		ret.Body = body
		ret.AddExpr(body)
	}

	return ret, nil
}

// ReadForBlock ...
func (I *GigoInterpreter) ReadForBlock(templated bool) (*glang.ForStmt, error) {
	var ret *glang.ForStmt

	if I.Read(glanglexer.ForToken) == nil {
		return nil, I.Debug("unexpected token", glanglexer.ForToken)
	}

	ret = glang.NewForStmt()
	I.blockscope.Enter()
	defer I.blockscope.Leave()

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	whatis := I.PeekUntil(glanglexer.RangeToken, glanglexer.TypeAssignToken, glanglexer.BraceOpenToken)

	if whatis == nil {
		return nil, I.Debug("unexpected token", glanglexer.RangeToken, glanglexer.TypeAssignToken, glanglexer.BraceOpenToken)
	}

	if whatis.GetType() == glanglexer.RangeToken {
		I.Read(glanglexer.RangeToken)
		I.ReadMany(genericlexer.WsToken)
		ret.AddExprs(I.Emit())

		block, err := I.ReadExpressionBlock(templated, glanglexer.BraceOpenToken)
		if err != nil {
			return nil, err
		}
		ret.Post = block
		ret.AddExpr(block)

		I.ReadMany(genericlexer.WsToken)
		ret.AddExprs(I.Emit())

	} else if whatis.GetType() == glanglexer.TypeAssignToken {

		I.Read(glanglexer.TypeAssignToken)
		I.ReadMany(genericlexer.WsToken)

		if I.Read(glanglexer.RangeToken) != nil {

			I.RewindAll()
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())

			init, err := I.ReadAssignExpr(templated, false, glanglexer.RangeToken, glanglexer.SemiColonToken)
			if err != nil {
				return nil, err
			}
			ret.Init = init
			ret.AddExpr(init)
			I.blockscope.AddVar(init.CollectVarNames()...)

			I.ReadMany(genericlexer.WsToken)
			I.Read(glanglexer.RangeToken)
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())

			block, err := I.ReadExpressionBlock(templated, glanglexer.BraceOpenToken)
			if err != nil {
				return nil, err
			}
			ret.Post = block
			ret.AddExpr(block)

			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())

		} else {

			I.RewindAll()
			init, err := I.ReadAssignExpr(templated, false, glanglexer.SemiColonToken)
			if err != nil {
				return nil, err
			}
			ret.Init = init
			ret.AddExpr(init)
			I.blockscope.AddVar(init.CollectVarNames()...)

			I.ReadMany(genericlexer.WsToken)
			I.Read(glanglexer.SemiColonToken)
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())

			cond, err := I.ReadExpression()
			if err != nil {
				return nil, err
			}
			ret.Cond = cond
			ret.AddExpr(cond)
			I.ReadMany(genericlexer.WsToken)
			I.Read(glanglexer.SemiColonToken)
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())

			post, err := I.ReadExpressionBlock(templated, glanglexer.BraceOpenToken)
			if err != nil {
				return nil, err
			}
			ret.Post = post
			ret.AddExpr(post)
			I.ReadMany(genericlexer.WsToken)
			I.Read(glanglexer.SemiColonToken)
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())
		}

	} else if whatis.GetType() == glanglexer.BraceOpenToken {
		I.RewindAll()
		block, err := I.ReadExpressionBlock(templated, glanglexer.BraceOpenToken)
		if err == nil {
			ret.AddExpr(block)
		}
	}

	I.ReadMany(genericlexer.WsToken, glanglexer.NlToken)
	ret.AddExprs(I.Emit())

	body, err := I.ReadExpressionsBlock(templated, glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
	if err != nil {
		return nil, err
	}
	ret.Body = body
	ret.AddExpr(body)
	ret.AddExprs(I.Emit())

	return ret, nil
}

// ReadExpressionBlock ...
func (I *GigoInterpreter) ReadExpressionBlock(templated bool, until ...lexer.TokenType) (*glang.ExpressionDecl, error) {
	var ret *glang.ExpressionDecl

	if len(I.Current()) > 0 {
		fmt.Printf("I.Current %q\n", I.Current())
		fmt.Println("I.Current", I.Current())
		panic("nop")
	}

	ret = glang.NewExpressionDecl()

	for {

		if I.Peek(
			glanglexer.StructToken,
		) != nil {
			block, err := I.ReadStructDecl(templated)
			if err != nil {
				return nil, err
			}
			ret.AddExpr(block)
			// read the body
			I.Read(genericlexer.WsToken)
			I.ReadBlock(glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
			ret.AddExprs(I.Emit())
			break
		}

		if I.Read(
			genericlexer.TextToken,
			glanglexer.TrueToken,
			glanglexer.FalseToken,
		) != nil {
			ret.AddExprs(I.Emit())
			break
		}

		I.Read(
			genericlexer.WsToken,
		)

		if I.Read(
			glanglexer.SubToken,
			glanglexer.MulToken,
			glanglexer.AddToken,
			glanglexer.RemToken,
			glanglexer.QuoToken,
		) != nil {
			ret.AddExprs(I.Emit())
		}

		I.Read(
			genericlexer.WsToken,
		)

		if I.Peek(glanglexer.FuncToken) != nil {
			block, err := I.ReadFuncDecl(templated, true)
			if err != nil {
				return nil, err
			}
			ret.AddExpr(block)
		}

		if I.Peek(glanglexer.BracketOpenToken) != nil {
			I.ReadBlock(glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
			ret.AddExprs(I.Emit())
		}

		if I.Peek(glanglexer.TplOpenToken, genericlexer.WordToken) != nil {

			v, err := I.ReadVarName(templated, true, true)
			if err != nil {
				n, err2 := I.ReadNumber()
				if n != nil {
					ret.AddExpr(n)
					continue
				} else {
					fmt.Println(err)
					fmt.Println(I.PeekN(5))
					fmt.Println(err2)
					fmt.Println(templated)
					panic(err)
					break
				}
			}

			I.Read(
				genericlexer.WsToken,
			)

			doBraces := !I.blockscope.HasVar(v.GetVarName())

			if I.Read(glanglexer.IncToken, glanglexer.DecToken) != nil {
				ret.AddExpr(v)
				ret.AddExprs(I.Emit())

			} else if I.Peek(glanglexer.BraceOpenToken) != nil {
				ret.AddExpr(v)
				if doBraces {
					I.ReadBlock(glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
					ret.AddExprs(I.Emit())
				}
				I.ReadMany(genericlexer.WsToken)

			} else if I.Peek(glanglexer.ParenOpenToken) != nil {

				callexpr := glang.NewCallExpr()
				callexpr.ID = v
				callexpr.AddExpr(v)
				callexpr.AddExprs(I.Emit())

				block, err := I.ReadParenExprBlock(templated)
				if err != nil {
					return nil, err
				}
				callexpr.Params = block
				callexpr.AddExpr(block)

				ret.AddExpr(callexpr)

			} else if I.Peek(glanglexer.BracketOpenToken) != nil {
				ret.AddExpr(v)
				I.ReadBlock(glanglexer.BracketOpenToken, glanglexer.BracketCloseToken)
				ret.AddExprs(I.Emit())

			} else if I.Peek(
				glanglexer.TypeAssignToken,
				glanglexer.AssignToken,
			) != nil {
				ret.AddExpr(v)
				I.Read(
					glanglexer.TypeAssignToken,
					glanglexer.AssignToken,
				)
				ret.AddExprs(I.Emit())
			} else {
				ret.AddExpr(v)
			}
		}

		if I.Peek(until...) != nil {
			break
		}

		if I.Peek(glanglexer.SemiColonToken) != nil {
			break
		}

		if I.Peek(glanglexer.NlToken) != nil {
			break
		}
	}

	return ret, nil
}

/*
ReadBinaryExpressionBlock reads a binary expression such as

	x == y && false == true || "whatever"
*/
func (I *GigoInterpreter) ReadBinaryExpressionBlock(templated bool, until lexer.TokenType) (*glang.BinaryExpr, error) {
	var ret *glang.BinaryExpr

	ret = glang.NewBinaryExpr()

	for {

		exl, err := I.ReadExpressionBlock(templated,
			until,
			glanglexer.GreaterToken,
			glanglexer.SmeqToken,
			glanglexer.GteqToken,
			glanglexer.NeqToken,
			glanglexer.EqToken,
			glanglexer.SmallerToken,
			glanglexer.GreaterToken,
		)
		if err != nil {
			return nil, err
		}
		ret.AddExpr(exl)
		ret.Left = exl

		I.ReadMany(genericlexer.WsToken)

		if op := I.Read(glanglexer.GreaterToken,
			glanglexer.SmeqToken,
			glanglexer.GteqToken,
			glanglexer.NeqToken,
			glanglexer.EqToken,
			glanglexer.SmallerToken,
			glanglexer.GreaterToken,
		); op != nil {
			ret.AddExprs(I.Emit())
			ret.Op = op
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())

			exr, err := I.ReadExpressionBlock(templated,
				glanglexer.BraceOpenToken,
				glanglexer.LAndToken,
				glanglexer.LOrToken,
			)
			if err != nil {
				return nil, err
			}
			ret.AddExpr(exr)
			ret.Right = exr
		}

		if I.Peek(until) != nil {
			I.RewindAll()
			break
		}

		I.ReadMany(genericlexer.WsToken)
		if I.Read(glanglexer.LAndToken, glanglexer.LOrToken) != nil {
			I.ReadMany(genericlexer.WsToken)
			I.Read(glanglexer.NlToken)
			ret.AddExprs(I.Emit())
			continue
		}

		if I.Peek(until, glanglexer.NlToken) != nil {
			I.RewindAll()
			break
		}
	}
	return ret, nil
}

// ReadAssignExpr reads a block of expressions.
// returns an error if none is found.
// The next token to analyze must be of type open,
// the block must end with a token of type close.
// In between data are read as a golang assign,
// a = x
// a, b = x
// a, b = x, u
func (I *GigoInterpreter) ReadAssignExpr(
	templated bool,
	allowdot bool,
	until ...lexer.TokenType,
) (*glang.AssignExpr, error) {

	var ret *glang.AssignExpr

	ret = glang.NewAssignExpr()
	w := "ids"
	defer I.blockscope.AddVar(ret.CollectVarNames()...)

	for {
		if I.Peek(glanglexer.NlToken) != nil {
			break
		}
		if p := I.Peek(until...); p != nil {
			break
		}
		if I.Read(glanglexer.TypeAssignToken, glanglexer.AssignToken) != nil {
			w = "values"
			I.ReadMany(genericlexer.WsToken)
			ret.AddExprs(I.Emit())
			continue
		}
		if w == "ids" {
			ID, err := I.ReadVarName(templated, true, allowdot)
			if err != nil {
				return nil, err
			}
			ret.AddID(ID)
			ret.AddExpr(ID)

			I.ReadMany(genericlexer.WsToken)
			if I.Read(glanglexer.CommaToken) != nil {
				I.ReadMany(genericlexer.WsToken)
				if I.Read(glanglexer.NlToken) != nil {
					ret.AddExprs(I.Emit())
					continue
				}
			}
			ret.AddExprs(I.Emit())

		} else {
			block, err := I.ReadExpressionBlock(templated, glanglexer.CommaToken)
			if err != nil {
				return nil, err
			}
			ret.AddValue(block)
			ret.AddExpr(block)

			I.ReadMany(genericlexer.WsToken)
			if I.Read(glanglexer.CommaToken) != nil {
				I.ReadMany(genericlexer.WsToken)
				if I.Read(glanglexer.NlToken) != nil {
					ret.AddExprs(I.Emit())
					continue
				}
			}
			ret.AddExprs(I.Emit())
		}
	}
	return ret, nil
}

// ReadCallExpr reads a block of expressions.
// returns an error if none is found.
// The next token to analyze must be of type open,
// the block must end with a token of type close.
// In between data are read as a golang block of expressions,
// x(y(x, t))
// => (y(x, t))
// => (x, t)
func (I *GigoInterpreter) ReadCallExpr(
	templated bool,
) (*glang.CallExpr, error) {

	var ret *glang.CallExpr

	ret = glang.NewCallExpr()
	ret.ID = glang.NewIdentifierDecl()
	ret.AddExpr(ret.ID)

	for {
		if I.Peek(glanglexer.ParenOpenToken) != nil {
			ret.ID.AddExprs(I.Emit())
			block, err := I.ReadParenExprBlock(templated)
			if err != nil {
				return nil, err
			}
			ret.Params = block
			ret.AddExpr(block)
			break
		}
		if I.Read(genericlexer.WordToken, glanglexer.DotToken) != nil {
			ret.ID.AddExprs(I.Emit())
		}
		if templated && I.Peek(glanglexer.TplOpenToken) != nil {
			block, err := I.ReadTemplateBlock()
			if err != nil {
				return nil, err
			}
			ret.ID.AddExpr(block)
		}
		I.ReadMany(genericlexer.WsToken)
		ret.AddExprs(I.Emit())
		if I.Read(glanglexer.ParenCloseToken) != nil {
			break
		}
	}
	return ret, nil
}

// ReadParenExprBlock reads a block of expressions.
// returns an error if none is found.
// The next token to analyze must be of type open,
// the block must end with a token of type close.
// In between data are read as a golang block of expressions,
// x(y(x, t))
// => (y(x, t))
// => (x, t)
func (I *GigoInterpreter) ReadParenExprBlock(
	templated bool,
) (*glang.CallExprBlock, error) {

	var ret *glang.CallExprBlock

	open := glanglexer.ParenOpenToken
	close := glanglexer.ParenCloseToken

	openTok := I.Read(open)
	if openTok == nil {
		return nil, I.Debug("Failed to ReadParenExprBlock", open)
	}
	count := 1

	ret = glang.NewCallExprBlock()

	for {
		if openTok := I.Read(open); openTok != nil {
			count++

		} else if closeTok := I.Read(close); closeTok != nil {
			count--
			if count == 0 {
				I.Rewind()
				break
			}

		} else if I.Read(glanglexer.ElipseToken) != nil {

		} else if I.Read(glanglexer.CommaToken) != nil {
			I.Read(genericlexer.WsToken)
			I.Read(glanglexer.NlToken)

		} else {
			ret.AddExprs(I.Emit())

			block, err := I.ReadExpressionBlock(templated,
				glanglexer.ParenCloseToken,
				glanglexer.CommaToken,
				glanglexer.ElipseToken,
			)
			if err != nil {
				return nil, err
			}
			ret.AddParam(block)
			ret.AddExprs(I.Emit())
			ret.AddExpr(block)
			I.ReadMany(genericlexer.WsToken)
		}
	}
	if I.Read(close) == nil {
		return nil, I.Debug("unexpected token", close)
	}
	ret.AddExprs(I.Emit())
	return ret, nil
}

// ReadFuncDecl reads a func with or without receiver.
// the next token must be a funcToken
// returns an error if none is found.
// func (r) AddExpr(expr Tokener) out { body }
func (I *GigoInterpreter) ReadFuncDecl(templated bool, isLiteral bool) (*glang.FuncDecl, error) {

	funcTok := I.Read(glanglexer.FuncToken)
	if funcTok == nil {
		return nil, I.Debug("Unexpected token", glanglexer.FuncToken)
	}
	ret := glang.NewFuncDecl()
	I.blockscope.Enter()
	defer I.blockscope.Leave()

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	if isLiteral == false {
		receiver, err := I.ReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		if receiver != nil {
			if err != nil {
				return nil, err
			}
			ret.Receiver = receiver
			receiver.AddExprs(I.Emit())
			ret.AddExpr(receiver)
			I.blockscope.AddVar(receiver.CollectVarNames()...)
			I.ReadMany(genericlexer.WsToken)
		}
	}

	ID, _ := I.ReadVarName(templated, false, false)
	if ID != nil {
		ret.Name = ID
		ret.AddExpr(ID)
	}

	I.ReadMany(genericlexer.WsToken)
	params, err := I.ReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
	if err != nil {
		return nil, err
	}
	params.AddExprs(I.Emit())
	ret.Params = params
	ret.AddExpr(params)
	I.blockscope.AddVar(params.CollectVarNames()...)

	I.ReadMany(genericlexer.WsToken)
	bodyStart := I.Peek(glanglexer.BraceOpenToken)
	if bodyStart == nil {
		if I.Peek(glanglexer.TplOpenToken) != nil {
			block, err2 := I.ReadVarName(templated, false, false)
			if err2 != nil {
				return nil, err2
			}
			if block != nil {
				ret.AddExpr(block)
			}
		} else {
			out, _ := I.ReadParenDecl(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
			if out != nil {
				params.AddExprs(I.Emit())
				ret.Out = out
				ret.AddExpr(out)
				I.blockscope.AddVar(out.CollectVarNames()...)
			} else {
				outTok, _ := I.ReadTypeName(templated, true)
				if outTok != nil {
					out := glang.NewPropsBlockDecl()
					out.AddT(outTok)
					out.AddExpr(outTok)
					out.AddExprs(I.Emit())
					ret.AddExpr(out)
					ret.Out = out
				}
			}
		}
	}

	I.ReadMany(genericlexer.WsToken)
	body, err := I.ReadExpressionsBlock(templated, glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
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

	if I.Peek(glanglexer.BraceOpenToken) == nil {
		return nil, I.Debug("unexpected token", glanglexer.BraceOpenToken)
	}

	ret := glang.NewStructDecl()
	ret.AddExprs(I.Emit())

	block, err := I.ReadPropsBlock(templated, glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
	if err != nil {
		return nil, err
	}
	ret.AddExpr(block)
	ret.Block = block
	return ret, nil
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

	ID, err := I.ReadVarName(true, false, false)
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

// ReadTemplateExprDecl ...
func (I *GigoInterpreter) ReadTemplateExprDecl() (*glang.TemplateFuncDecl, error) {

	fn := glang.NewTemplateFuncDecl()
	I.KeepPreviousComment()

	I.ReadMany(
		glanglexer.NlToken,
		genericlexer.CommentLineToken,
		genericlexer.CommentBlockToken,
		genericlexer.WsToken)
	fn.AddExprs(I.Emit())

	// smthig to improve here
	block, err := I.ReadTemplateBlock()
	if err != nil {
		return nil, err
	}
	block.AddExprs(I.Emit())
	fn.Modifier = block
	fn.AddExpr(block)

	I.ReadMany(genericlexer.WsToken)
	fn.AddExprs(I.Emit())

	nFunc, err := I.ReadFuncDecl(true, false)
	if err != nil {
		return nil, err
	}
	fn.Func = nFunc
	fn.AddExpr(fn.Func)

	return fn, nil
}

// ReadVarDecl reads a var declaration.
// the next token must be a VarToken
// returns an error if none is found.
// var (
// 		x = ""
// )
func (I *GigoInterpreter) ReadVarDecl(templated bool) (*glang.VarDecl, error) {

	varTok := I.Read(glanglexer.VarToken)
	if varTok == nil {
		return nil, I.Debug("unexpected token", glanglexer.VarToken)
	}

	I.ReadMany(genericlexer.WsToken)
	ret := glang.NewVarDecl()
	ret.AddExprs(I.Emit())
	defer I.blockscope.AddVar(ret.CollectVarNames()...)

	if I.Peek(glanglexer.ParenOpenToken) == nil {

		assignment := glang.NewAssignDecl()
		assignment.AddExprs(I.Emit())

		left, err2 := I.ReadVarName(templated, true, false)
		if err2 != nil {
			return nil, err2
		}
		assignment.AddExpr(left)

		I.ReadMany(genericlexer.WsToken)
		assignment.AddExprs(I.Emit())

		var err error
		var leftType genericinterperter.Tokener
		eq := I.Read(glanglexer.AssignToken)
		if eq == nil {
			leftType, err = I.ReadTypeName(templated, true)
			if leftType != nil && err != nil { // weird, need both var ?
				return ret, err
			}
			assignment.AddExpr(leftType)
			I.ReadMany(genericlexer.WsToken)
			eq = I.Read(glanglexer.AssignToken)
		}
		if eq != nil {
			I.ReadMany(genericlexer.WsToken)
			assignment.AddExprs(I.Emit())

			right, err := I.ReadExpressionBlock(templated)
			if err != nil {
				return nil, err
			}
			assignment.Left = left
			assignment.LeftType = leftType
			assignment.Assign = eq
			assignment.Right = right
			if right != nil {
				assignment.AddExpr(right)
			}
		}
		ret.AddAssignment(assignment)
		ret.AddExpr(assignment)

	} else {
		block, err := I.ReadAssignsBlock(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		if err != nil {
			return nil, err
		}
		ret.AddExpr(block)
		ret.AddAssignment(block)
	}

	return ret, nil
}

// ReadConstDecl reads a const declaration.
// the next token must be a ConstToken
// returns an error if none is found.
// const (
// 		x = ""
// )
func (I *GigoInterpreter) ReadConstDecl(templated bool) (*glang.ConstDecl, error) {

	constTok := I.Read(glanglexer.ConstToken)
	if constTok == nil {
		return nil, I.Debug("unexpected token", glanglexer.ConstToken)
	}

	ret := glang.NewConstDecl()
	defer I.blockscope.AddVar(ret.CollectVarNames()...)

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	if I.Peek(glanglexer.ParenOpenToken) == nil {

		assignment := glang.NewAssignDecl()
		assignment.AddExprs(I.Emit())

		left, err2 := I.ReadVarName(templated, false, false)
		if err2 != nil {
			return nil, err2
		}
		assignment.AddExpr(left)

		I.ReadMany(genericlexer.WsToken)
		assignment.AddExprs(I.Emit())

		var err error
		var leftType genericinterperter.Tokener
		eq := I.Read(glanglexer.AssignToken)
		if eq == nil {
			leftType, _ = I.ReadTypeName(templated, true)
			if leftType != nil {
				assignment.AddExpr(leftType)
				I.ReadMany(genericlexer.WsToken)
			}
			eq = I.Read(glanglexer.AssignToken)
		}
		if eq == nil {
			return nil, I.Debug("unexpected token", glanglexer.AssignToken)
		}
		I.ReadMany(genericlexer.WsToken)
		assignment.AddExprs(I.Emit())

		right, err := I.ReadExpressionBlock(templated)
		if err != nil {
			return nil, err
		}
		assignment.Left = left
		assignment.LeftType = leftType
		assignment.Assign = eq
		assignment.Right = right
		if right != nil {
			assignment.AddExpr(right)
		}
		ret.AddAssignment(assignment)
		ret.AddExpr(assignment)

	} else {
		block, err := I.ReadAssignsBlock(templated, glanglexer.ParenOpenToken, glanglexer.ParenCloseToken)
		if err != nil {
			return nil, err
		}
		// block.AddExprs(I.Emit())
		ret.AddExpr(block)
		ret.AddAssignment(block)
	}

	return ret, nil
}

// ReadAssignDecl reads an assign declaration.
// name type eq value
func (I *GigoInterpreter) ReadAssignDecl() (*glang.AssignDecl, error) {

	assignment := glang.NewAssignDecl()
	assignment.AddExprs(I.Emit())
	left, err := I.ReadVarName(false, true, false)
	if err != nil {
		return nil, err
	}
	assignment.AddExpr(left)

	I.ReadMany(genericlexer.WsToken)
	assignment.AddExprs(I.Emit())

	var err2 error
	var leftType genericinterperter.Tokener
	eq := I.Read(glanglexer.TypeAssignToken)
	if eq == nil {
		leftType, err2 = I.ReadTypeName(false, true)
		if leftType != nil && err2 != nil { // weird, need both var ?
			return nil, err2
		}
		assignment.AddExpr(leftType)
		I.ReadMany(genericlexer.WsToken)
		eq = I.Read(glanglexer.TypeAssignToken)
	}
	if eq == nil {
		return nil, I.Debug("unexpected token", glanglexer.TypeAssignToken)
	}
	I.ReadMany(genericlexer.WsToken)
	assignment.AddExprs(I.Emit())

	right, err2 := I.ReadExpression()
	if err2 != nil {
		return nil, err2
	}
	assignment.Left = left
	assignment.LeftType = leftType
	assignment.Assign = eq
	assignment.Right = right
	if right != nil {
		assignment.AddExpr(right)
	}

	return assignment, nil
}

// ReadReturnDecl reads a return declaration.
func (I *GigoInterpreter) ReadReturnDecl(templated bool) (*glang.ReturnDecl, error) {

	var ret *glang.ReturnDecl
	tok := I.Read(glanglexer.ReturnToken)
	if tok == nil {
		return nil, I.Debug("unexpected token", glanglexer.ReturnToken)
	}
	ret = glang.NewReturnDecl()
	ret.AddExprs(I.Emit())

	for {
		I.ReadMany(
			genericlexer.WsToken,
			genericlexer.CommentLineToken,
			genericlexer.CommentBlockToken,
		)
		ret.AddExprs(I.Emit())

		// should be a value identifier
		ID, err := I.ReadExpressionBlock(templated, glanglexer.CommaToken)
		if err != nil {
			fmt.Printf("%+v", err)
			fmt.Printf("%#v", err)
			panic(err)
			break
		}
		ret.AddExpr(ID)

		I.ReadMany(
			genericlexer.WsToken,
			genericlexer.CommentLineToken,
			genericlexer.CommentBlockToken,
		)

		if I.Read(glanglexer.CommaToken) != nil {
			I.Read(glanglexer.NlToken)
		} else if I.Read(glanglexer.NlToken) != nil {
			break
		} else if I.Peek(glanglexer.BraceCloseToken) != nil {
			break
		}
		ret.AddExprs(I.Emit())
	}
	ret.AddExprs(I.Emit())

	return ret, nil
}

// ReadTemplateBlock reads a template block.
// <:...>
func (I *GigoInterpreter) ReadTemplateBlock() (*glang.BodyBlockDecl, error) {
	if p := I.Peek(glanglexer.TplOpenToken); p != nil {
		ret, err := I.ReadBodyBlock(glanglexer.TplOpenToken, glanglexer.GreaterToken)
		if err != nil {
			return nil, err
		}
		ret.AddExprs(I.Emit())
		ret.Open.SetType(glanglexer.TplOpenToken)
		ret.Close.SetType(glanglexer.TplCloseToken)
		return ret, nil
	}
	return nil, I.Debug("unexpeced token", glanglexer.TplOpenToken)
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

	ID, err := I.ReadVarName(true, false, false)
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

	implTemplate, err := I.ReadTemplateBlock()
	if err != nil {
		return nil, err
	}
	implTemplate.AddExprs(I.Emit())
	ret.AddExpr(implTemplate)
	ret.ImplementTemplate = implTemplate

	I.ReadMany(genericlexer.WsToken)
	ret.AddExprs(I.Emit())

	block, err := I.ReadPropsBlock(true, glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
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

	block, err := I.ReadSignsBlock(glanglexer.BraceOpenToken, glanglexer.BraceCloseToken)
	if err != nil {
		return nil, err
	}
	ret.AddExpr(block)
	ret.Block = block

	return ret, nil
}
