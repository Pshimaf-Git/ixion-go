package semantic

import (
	"ixion/internal/ast"
)

const (
	unknownType = "unknown type"
	funcType    = "function"

	intType    = "int"
	stringType = "string"
)

func (a *Analyzer) visitProgram(program *ast.Program) {
	for _, stmt := range program.Statements {
		a.visitStmt(stmt)
	}
}

func (a *Analyzer) visitStmt(stmt ast.Statement) {
	switch x := stmt.(type) {
	case *ast.VarStatement:
		a.visitVarStmt(x)
	case *ast.ExpressionStatement:
		a.visitExpressionStmt(x)
	case *ast.ReturnStatement:
		a.visitReturnStmt(x)
	case *ast.PrintStatement:
		a.vistPrintStmt(x)
	case *ast.FunctionDeclaration:
		a.visitFuncDecl(x)
	case *ast.BlockStatement:
		a.visitBlockStmt(x)
	}
}

func (a *Analyzer) visitVarStmt(vs *ast.VarStatement) {
	if symbol := a.resolve(vs.Name.Value); symbol != nil {
		a.errf(vs, "variable '%s' already declare", vs.Name.Value)
	}

	var varType string
	if vs.Type != nil {
		varType = vs.Type.String()
	} else if vs.Value != nil {
		varType = a.getExprType(vs.Value)
	} else {
		varType = unknownType
	}

	if !a.declare(vs.Name.Value, varType) {
		a.errf(vs, "variable '%s' already declare in these scope", vs.Name.Value)
	}

	if vs.Value != nil {
		a.visitExpression(vs.Value)
	}
}

func (a *Analyzer) visitExpressionStmt(es *ast.ExpressionStatement) {
	a.visitExpression(es.Expression)
}

func (a *Analyzer) vistPrintStmt(ps *ast.PrintStatement) {
	if ps.Value != nil {
		a.visitExpression(ps.Value)
	}
}

func (a *Analyzer) visitReturnStmt(rs *ast.ReturnStatement) {
	if rs.ReturnValue != nil {
		a.visitExpression(rs.ReturnValue)
	}
}

func (a *Analyzer) visitFuncDecl(fd *ast.FunctionDeclaration) {
	if !a.declare(fd.Name.Value, funcType) {
		a.errf(fd, "function '%s' already declare", fd.Name.Value)
	}

	a.enterScope()

	for _, param := range fd.Parameters {
		paramType := unknownType
		if param.Type != nil {
			paramType = param.Type.String()
		}
		if !a.declare(param.Name.Value, paramType) {
			a.errf(param, "parameter '%s' already declared", param.Name.Value)
		}
	}

	a.visitBlockStmt(fd.Body)

	a.exitScope()
}

func (a *Analyzer) visitBlockStmt(bs *ast.BlockStatement) {
	a.enterScope()
	for _, stmt := range bs.Statements {
		a.visitStmt(stmt)
	}
	a.exitScope()
}

func (a *Analyzer) visitExpression(expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.Identifier:
		a.visitIdentifier(e)
	case *ast.IntegerLiteral:
		// Ничего не проверяем для литералов
	case *ast.StringLiteral:
		// Ничего не проверяем для литералов
	case *ast.PrefixExpression:
		a.visitPrefixExpression(e)
	case *ast.InfixExpression:
		a.visitInfixExpression(e)
	case *ast.AssignmentExpression:
		a.visitAssignmentExpression(e)
	case *ast.CallExpression:
		a.visitCallExpression(e)
	}
}

func (a *Analyzer) visitIdentifier(id *ast.Identifier) {
	// Проверяем, объявлена ли переменная
	if symbol := a.resolve(id.Value); symbol == nil {
		a.errf(id, "undeclared variable '%s'", id.Value)
	}
}

func (a *Analyzer) visitPrefixExpression(pe *ast.PrefixExpression) {
	a.visitExpression(pe.Right)
	// TODO: Check compatibility of operand type with operator
}

func (a *Analyzer) visitInfixExpression(ie *ast.InfixExpression) {
	a.visitExpression(ie.Left)
	a.visitExpression(ie.Right)
	// TODO: Check compatibility of operand types
}

func (a *Analyzer) visitAssignmentExpression(ae *ast.AssignmentExpression) {
	if ident, ok := ae.Left.(*ast.Identifier); ok {
		if symbol := a.resolve(ident.Value); symbol == nil {
			a.errf(ae, "cannot assign to undeclared variable '%s'", ident.Value)
		}
	} else {
		a.err(ae, "left side of assignment must be an identifier")
	}

	a.visitExpression(ae.Value)
}

func (a *Analyzer) visitCallExpression(ce *ast.CallExpression) {
	if ident, ok := ce.Function.(*ast.Identifier); ok {
		if symbol := a.resolve(ident.Value); symbol == nil {
			a.errf(ce, "call to undeclared function '%s'", ident.Value)
		} else if symbol.Type != funcType {
			a.errf(ce, "'%s' is not a function", ident.Value)
		}
	}

	for _, arg := range ce.Arguments {
		a.visitExpression(arg)
	}
}

func (a *Analyzer) getExprType(expr ast.Expression) string {
	switch x := expr.(type) {
	case *ast.IntegerLiteral:
		return intType
	case *ast.StringLiteral:
		return stringType
	case *ast.Identifier:
		if symbol := a.resolve(x.Value); symbol != nil {
			return symbol.Type
		}
		return unknownType
	default:
		return unknownType
	}
}
