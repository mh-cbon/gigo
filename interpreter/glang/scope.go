package glang

// Scope points to a ScopeBlock,
// a ScopeBlock can have multiple parents, one children, infinite depth.
// Scope can
//- create ScopeBlock,
//- give current ScopeBlock,
//- return to parent ScopeBlock
//- It ca ntell if a variable exists in current block, or its ancestors.
type Scope struct {
	parents []*ScopeBlock
	current *ScopeBlock
}

// NewScope is a ctor.
func NewScope() *Scope {
	ret := &Scope{}
	ret.Enter()
	return ret
}

// Current returns current scope or nil.
func (s *Scope) Current() *ScopeBlock {
	return s.current
}

// HasVar returns true if the var is declared in this scope or parents.
func (s *Scope) HasVar(name string) bool {
	if s.current.HasVar(name) {
		return true
	}
	for i := len(s.parents) - 1; i >= 0; i-- {
		if s.parents[i].HasVar(name) {
			return true
		}
	}
	return false
}

// AllVars...
func (s *Scope) AllVars() []string {
	ret := []string{}
	ret = append(ret, s.current.Vars()...)
	for i := len(s.parents) - 1; i >= 0; i-- {
		ret = append(ret, s.parents[i].Vars()...)
	}
	return ret
}

// AddVar to the current scope.
func (s *Scope) AddVar(names ...string) {
	for _, name := range names {
		s.current.AddVar(name)
	}
}

// Leave the current scope.
func (s *Scope) Leave() {
	y := len(s.parents)
	s.current = nil
	if y > 0 {
		s.current = s.parents[y-1]
		s.parents = s.parents[:y-1]
	}
}

// Enter a new scope.
func (s *Scope) Enter() {
	if s.current != nil {
		s.parents = append(s.parents, s.current)
	}
	s.current = &ScopeBlock{vars: map[string]bool{}}
}

// ScopeBlock represents a block {} as a scope,
// it keep tracks of variable name
// it can tell if a variable is defined
type ScopeBlock struct {
	vars map[string]bool // super simplist until more is needed
}

// AddVar to the block scope.
func (s *ScopeBlock) AddVar(name string) {
	s.vars[name] = true
}

// Vars...
func (s *ScopeBlock) Vars() []string {
	ret := []string{}
	for k := range s.vars {
		ret = append(ret, k)
	}
	return ret
}

// HasVar returns true if the var is declared in this block scope.
func (s *ScopeBlock) HasVar(name string) bool {
	_, ok := s.vars[name]
	return ok
}
