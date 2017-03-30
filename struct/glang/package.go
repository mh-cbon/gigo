package glang

import (
	"fmt"

	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
)

type Package struct {
	Name  string
	Files []ScopeReceiver
}

type ScopeReceiver interface {
	genericinterperter.ScopeReceiver
	FindPackagesDecl() []*PackageDecl
	FindImplementsTypes() []*ImplementDecl
	FindStructsTypes() []*StructDecl
	FindTemplatesTypes() []*TemplateDecl
	FindFuncs() []*FuncDecl
	FindTemplateFuncs() []FuncDeclarer
	FindDefineFuncs() []*TemplateFuncDecl
}

func (p *Package) String() string {
	return fmt.Sprintf("%v", p.Name)
}

func (p *Package) FindPackagesDecl() []*PackageDecl {
	var ret []*PackageDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindPackagesDecl()...)
	}
	return ret
}
func (p *Package) FindImplementsTypes() []*ImplementDecl {
	var ret []*ImplementDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindImplementsTypes()...)
	}
	return ret
}
func (p *Package) FindStructsTypes() []*StructDecl {
	var ret []*StructDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindStructsTypes()...)
	}
	return ret
}
func (p *Package) FindTemplatesTypes() []*TemplateDecl {
	var ret []*TemplateDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindTemplatesTypes()...)
	}
	return ret
}
func (p *Package) FindFuncs() []*FuncDecl {
	var ret []*FuncDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindFuncs()...)
	}
	return ret
}
func (p *Package) FindTemplateFuncs() []FuncDeclarer {
	var ret []FuncDeclarer
	for _, f := range p.Files {
		ret = append(ret, f.FindTemplateFuncs()...)
	}
	return ret
}
func (p *Package) FindDefineFuncs() []*TemplateFuncDecl {
	var ret []*TemplateFuncDecl
	for _, f := range p.Files {
		ret = append(ret, f.FindDefineFuncs()...)
	}
	return ret
}

type SimplePackageRepository struct {
	Packages []*Package
}

func (s *SimplePackageRepository) AddToPackage(name string, scope ScopeReceiver) {
	pkg := s.GetPackage(scope.GetName())
	pkg.Files = append(pkg.Files, scope)
}
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
