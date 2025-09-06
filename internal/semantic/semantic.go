package semantic

import (
	"fmt"

	"ixion/internal/ast"
)

type Symbol struct {
	Name  string
	Type  string
	Scope *Scope
}

type Scope struct {
	Parent   *Scope
	Symbols  map[string]*Symbol
	Children []*Scope
}

func (s *Scope) exist(name string) bool {
	_, ok := s.Symbols[name]
	return ok
}

type Analyzer struct {
	CurrentScope *Scope
	GlobalScope  *Scope
	Errors       []error
}

func NewAnalyzer() *Analyzer {
	globalScope := &Scope{
		Symbols: make(map[string]*Symbol),
	}

	return &Analyzer{
		CurrentScope: globalScope,
		GlobalScope:  globalScope,
		Errors:       []error(nil),
	}
}

func (a *Analyzer) Analyze(program *ast.Program) []error {
	a.visitProgram(program)

	return a.Errors
}

func (a *Analyzer) err(node ast.Node, args ...any) {
	a.Errors = append(a.Errors, fmt.Errorf("%s: %s", node.TokenLiteral(), fmt.Sprint(args...)))
}

func (a *Analyzer) errf(node ast.Node, format string, args ...any) {
	a.Errors = append(a.Errors, fmt.Errorf("%s: %s", node.TokenLiteral(), fmt.Sprintf(format, args...)))
}

func (a *Analyzer) enterScope() {
	newScope := &Scope{
		Parent:  a.CurrentScope,
		Symbols: make(map[string]*Symbol),
	}

	a.CurrentScope.Children = append(a.CurrentScope.Children, newScope)
	a.CurrentScope = newScope
}

func (a *Analyzer) exitScope() {
	if a.CurrentScope.Parent != nil {
		a.CurrentScope = a.CurrentScope.Parent
	}
}

func (a *Analyzer) resolve(name string) *Symbol {
	curr := a.CurrentScope

	for curr != nil {
		if curr.exist(name) {
			return curr.Symbols[name]
		}
		curr = curr.Parent
	}

	return nil
}

func (a *Analyzer) declare(name string, _type string) bool {
	if a.CurrentScope.exist(name) {
		return false
	}

	a.CurrentScope.Symbols[name] = &Symbol{
		Name: name,
		Type: _type,
	}

	return true
}
