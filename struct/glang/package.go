package glang

import (
	"fmt"

	genericinterperter "github.com/mh-cbon/gigo/interpreter/generic"
)

type Package struct {
	Name  string
	Files []genericinterperter.ScopeReceiver
}

func (p *Package) String() string {
	return fmt.Sprintf("%v", p.Name)
}

type SimplePackageRepository struct {
	Packages []*Package
}

func (s *SimplePackageRepository) AddToPackage(name string, scope genericinterperter.ScopeReceiver) {
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
