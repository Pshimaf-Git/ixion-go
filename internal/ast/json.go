package ast

import (
	"encoding/json"
	"fmt"
)

// jsonProgram is an anonymous struct for JSON serialization of Program.
type jsonProgram struct {
	Statements []json.RawMessage `json:"statements"`
}

// MarshalJSON implements the json.Marshaler interface for Program.
func (p *Program) MarshalJSON() ([]byte, error) {
	jsonStatements := make([]json.RawMessage, len(p.Statements))
	for i, stmt := range p.Statements {
		raw, err := json.Marshal(stmt)
		if err != nil {
			return nil, err
		}
		jsonStatements[i] = raw
	}

	return json.Marshal(jsonProgram{Statements: jsonStatements})
}

// MarshalJSON for Expression interface
func (i *Identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string `json:"type"`
		Token string `json:"token_literal"`
		Value string `json:"value"`
	}{
		Type:  "Identifier",
		Token: i.TokenLiteral(),
		Value: i.Value,
	})
}

func (il *IntegerLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string `json:"type"`
		Token string `json:"token_literal"`
		Value int64  `json:"value"`
	}{
		Type:  "IntegerLiteral",
		Token: il.TokenLiteral(),
		Value: il.Value,
	})
}

func (sl *StringLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string `json:"type"`
		Token string `json:"token_literal"`
		Value string `json:"value"`
	}{
		Type:  "StringLiteral",
		Token: sl.TokenLiteral(),
		Value: sl.Value,
	})
}

func (tl *TypeLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string `json:"type"`
		Token string `json:"token_literal"`
		Value string `json:"value"`
	}{
		Type:  "TypeLiteral",
		Token: tl.TokenLiteral(),
		Value: tl.Value,
	})
}

func (pe *PrefixExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string     `json:"type"`
		Token    string     `json:"token_literal"`
		Operator string     `json:"operator"`
		Right    Expression `json:"right"`
	}{
		Type:     "PrefixExpression",
		Token:    pe.TokenLiteral(),
		Operator: pe.Operator,
		Right:    pe.Right,
	})
}

func (ie *InfixExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string     `json:"type"`
		Token    string     `json:"token_literal"`
		Left     Expression `json:"left"`
		Operator string     `json:"operator"`
		Right    Expression `json:"right"`
	}{
		Type:     "InfixExpression",
		Token:    ie.TokenLiteral(),
		Left:     ie.Left,
		Operator: ie.Operator,
		Right:    ie.Right,
	})
}

func (fl *FunctionLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string               `json:"type"`
		Token      string               `json:"token_literal"`
		Parameters []*FunctionParameter `json:"parameters"`
		ReturnType TypeExpression       `json:"return_type,omitempty"`
		Body       *BlockStatement      `json:"body"`
	}{
		Type:       "FunctionLiteral",
		Token:      fl.TokenLiteral(),
		Parameters: fl.Parameters,
		ReturnType: fl.ReturnType,
		Body:       fl.Body,
	})
}

func (ce *CallExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string       `json:"type"`
		Token     string       `json:"token_literal"`
		Function  Expression   `json:"function"`
		Arguments []Expression `json:"arguments"`
	}{
		Type:      "CallExpression",
		Token:     ce.TokenLiteral(),
		Function:  ce.Function,
		Arguments: ce.Arguments,
	})
}

func (ae *AssignmentExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string     `json:"type"`
		Token string     `json:"token_literal"`
		Left  Expression `json:"left"`
		Value Expression `json:"value"`
	}{
		Type:  "AssignmentExpression",
		Token: ae.TokenLiteral(),
		Left:  ae.Left,
		Value: ae.Value,
	})
}

func (fp *FunctionParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string         `json:"type"`
		Token     string         `json:"token_literal"`
		Name      *Identifier    `json:"name"`
		ParamType TypeExpression `json:"param_type,omitempty"`
	}{
		Type:      "FunctionParameter",
		Token:     fp.TokenLiteral(),
		Name:      fp.Name,
		ParamType: fp.Type,
	})
}

// MarshalJSON for Node interface (used by BlockStatement)
func (bs *BlockStatement) MarshalJSON() ([]byte, error) {
	jsonStatements := make([]json.RawMessage, len(bs.Statements))
	for i, stmt := range bs.Statements {
		raw, err := json.Marshal(stmt)
		if err != nil {
			return nil, err
		}
		jsonStatements[i] = raw
	}

	return json.Marshal(struct {
		Type       string            `json:"type"`
		Token      string            `json:"token_literal"`
		Statements []json.RawMessage `json:"statements"`
	}{
		Type:       "BlockStatement",
		Token:      bs.TokenLiteral(),
		Statements: jsonStatements,
	})
}

// Helper to marshal Expression interface
// This helper is no longer needed as MarshalJSON methods are defined for concrete types.
func marshalExpression(exp Expression) (json.RawMessage, error) {
	if exp == nil {
		return nil, nil
	}
	switch e := exp.(type) {
	case *Identifier:
		return json.Marshal(e)
	case *IntegerLiteral:
		return json.Marshal(e)
	case *StringLiteral:
		return json.Marshal(e)
	case *TypeLiteral:
		return json.Marshal(e)
	case *PrefixExpression:
		return json.Marshal(e)
	case *InfixExpression:
		return json.Marshal(e)
	case *FunctionLiteral:
		return json.Marshal(e)
	case *CallExpression:
		return json.Marshal(e)
	case *AssignmentExpression:
		return json.Marshal(e)
	case *FunctionParameter:
		return json.Marshal(e)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", exp)
	}
}

// Custom MarshalJSON for Expression interface to use marshalExpression helper
func (es *ExpressionStatement) MarshalJSON() ([]byte, error) {
	expressionJSON, err := marshalExpression(es.Expression)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Type       string          `json:"type"`
		Token      string          `json:"token_literal"`
		Expression json.RawMessage `json:"expression"`
	}{
		Type:       "ExpressionStatement",
		Token:      es.TokenLiteral(),
		Expression: expressionJSON,
	})
}

func (rs *ReturnStatement) MarshalJSON() ([]byte, error) {
	returnValueJSON, err := marshalExpression(rs.ReturnValue)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Type        string          `json:"type"`
		Token       string          `json:"token_literal"`
		ReturnValue json.RawMessage `json:"return_value"`
	}{
		Type:        "ReturnStatement",
		Token:       rs.TokenLiteral(),
		ReturnValue: returnValueJSON,
	})
}

func (ps *PrintStatement) MarshalJSON() ([]byte, error) {
	valueJSON, err := marshalExpression(ps.Value)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Type  string          `json:"type"`
		Token string          `json:"token_literal"`
		Value json.RawMessage `json:"value"`
	}{
		Type:  "PrintStatement",
		Token: ps.TokenLiteral(),
		Value: valueJSON,
	})
}

func (vs *VarStatement) MarshalJSON() ([]byte, error) {
	nameJSON, err := marshalExpression(vs.Name)
	if err != nil {
		return nil, err
	}
	valueJSON, err := marshalExpression(vs.Value)
	if err != nil {
		return nil, err
	}
	varTypeJSON, err := marshalExpression(vs.Type)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Type    string          `json:"type"`
		Token   string          `json:"token_literal"`
		Name    json.RawMessage `json:"name"`
		VarType json.RawMessage `json:"var_type,omitempty"`
		Value   json.RawMessage `json:"value"`
	}{
		Type:    "VarStatement",
		Token:   vs.TokenLiteral(),
		Name:    nameJSON,
		VarType: varTypeJSON,
		Value:   valueJSON,
	})
}

func (fd *FunctionDeclaration) MarshalJSON() ([]byte, error) {
	nameJSON, err := marshalExpression(fd.Name)
	if err != nil {
		return nil, err
	}

	paramsJSON := make([]json.RawMessage, len(fd.Parameters))
	for i, p := range fd.Parameters {
		paramJSON, err := marshalExpression(p)
		if err != nil {
			return nil, err
		}
		paramsJSON[i] = paramJSON
	}

	returnTypeJSON, err := marshalExpression(fd.ReturnType)
	if err != nil {
		return nil, err
	}

	bodyJSON, err := json.Marshal(fd.Body)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Type       string            `json:"type"`
		Token      string            `json:"token_literal"`
		Name       json.RawMessage   `json:"name"`
		Parameters []json.RawMessage `json:"parameters"`
		ReturnType json.RawMessage   `json:"return_type,omitempty"`
		Body       json.RawMessage   `json:"body"`
	}{
		Type:       "FunctionDeclaration",
		Token:      fd.TokenLiteral(),
		Name:       nameJSON,
		Parameters: paramsJSON,
		ReturnType: returnTypeJSON,
		Body:       bodyJSON,
	})
}
