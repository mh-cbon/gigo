package glang

import (
	"fmt"

	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
)

// Package isa set of multilpe files.
type Package struct {
	Name  string
	Files []ScopeReceiver
}

// ScopeReceiver is an interpreted source code.
type ScopeReceiver interface {
	genericinterperter.ScopeReceiver
	FindPackagesDecl() []*PackageDecl
	FindImplementsTypes() []*ImplementDecl
	FindStructsTypes() []*StructDecl
	FindTemplatesTypes() []*TemplateDecl
	FindFuncs() []*FuncDecl
	FindTemplateFuncs() []FuncDeclarer
	FindDefineFuncs() []*TemplateFuncDecl
	FindSymbols(string) []genericinterperter.Expressioner
	String() string
}

func (p *Package) String() string {
	return fmt.Sprintf("%v", p.Name)
}

// FindPackagesDecl returns all package declarations found.
func (p *Package) FindPackagesDecl() []*PackageDecl {
	var ret []*PackageDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindPackagesDecl()...)
	}
	return ret
}

// FindImplementsTypes returns all implements declarations found.
func (p *Package) FindImplementsTypes() []*ImplementDecl {
	var ret []*ImplementDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindImplementsTypes()...)
	}
	return ret
}

// FindStructsTypes returns all struct declarations found.
func (p *Package) FindStructsTypes() []*StructDecl {
	var ret []*StructDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindStructsTypes()...)
	}
	return ret
}

// FindTemplatesTypes returns all template declarations found.
func (p *Package) FindTemplatesTypes() []*TemplateDecl {
	var ret []*TemplateDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindTemplatesTypes()...)
	}
	return ret
}

// FindFuncs returns all func declarations found.
func (p *Package) FindFuncs() []*FuncDecl {
	var ret []*FuncDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindFuncs()...)
	}
	return ret
}

// FindTemplateFuncs returns all funcs with templating declarations found.
func (p *Package) FindTemplateFuncs() []FuncDeclarer {
	var ret []FuncDeclarer
	for _, f := range p.Files {
		ret = append(ret, f.FindTemplateFuncs()...)
	}
	return ret
}

// FindDefineFuncs returns all <define> declarations found.
func (p *Package) FindDefineFuncs() []*TemplateFuncDecl {
	var ret []*TemplateFuncDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindDefineFuncs()...)
	}
	return ret
}

// SimplePackageRepository is the reference of all apckages created.
type SimplePackageRepository struct {
	Packages []*Package
}

// AddToPackage adds given scope to a package of given name.
func (s *SimplePackageRepository) AddToPackage(name string, scope ScopeReceiver) {
	pkg := s.GetPackage(scope.GetName())
	pkg.Files = append(pkg.Files, scope)
}

// GetPackage creates a new package or return existing one.
func (s *SimplePackageRepository) GetPackage(name string) *Package {
	for _, p := range s.Packages {
		if p.Name == name {
			return p
		}
	}
	p := &Package{Name: name}
	s.Packages = append(s.Packages, p)
	return p
}
