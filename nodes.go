package junklang

import (
	"fmt"
)

// Abstract Syntax Tree
// Implementerer advanceret syntax support til aritmetik

type Node interface {
	Call(*Scope) any
}

// Node til at deklarere ny variabel til en abstrakt værdi
type DeclareNode struct {
	name  string
	value Node
}

func (n *DeclareNode) Call(s *Scope) any {
	val := n.value.Call(s)
	// Gem variablen i scopet
	s.Set(n.name, val)
	return val
}

// Node til at printe til terminalen
type PrintNode struct {
	value Node
}

func (n *PrintNode) Call(s *Scope) any {
	val := n.value.Call(s)
	fmt.Println(val)
	return nil
}

// Node til at repræsentere tal
type NumberNode struct {
	value float64
}

func (n *NumberNode) Call(_ *Scope) any {
	return n.value
}

type ReferenceNode struct {
	name string
}

func (n *ReferenceNode) Call(s *Scope) any {
	return s.Get(n.name)
}

// Node til aritmetik
type OperationNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (n *OperationNode) Call(s *Scope) any {
	left := n.Left.Call(s).(float64)
	right := n.Right.Call(s).(float64)

	switch n.Operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right
	default:
		panic("Forkert operator: " + n.Operator)
	}
}

// Node til at definere funktioner
type FuncDeclareNode struct {
	name   string
	params []string
	body   []Node
}

func (n *FuncDeclareNode) Call(s *Scope) any {
	s.SetFunction(n.name, n)
	return nil
}

// Node til at kalde functioner
type FuncCallNode struct {
	name string
	args []Node
}

func (n *FuncCallNode) Call(s *Scope) any {
	fn := s.GetFunction(n.name).(*FuncDeclareNode)

	funcScope := NewScope(s)

	for i, param := range fn.params {
		funcScope.Set(param, n.args[i].Call(s))
	}

	var res any
	for _, node := range fn.body {
		res = node.Call(funcScope)
		if _, isReturn := node.(*ReturnNode); isReturn {
			return res
		}
	}

	return res
}

// Node til at returnere fra funktioner
type ReturnNode struct {
	value Node
}

func (n *ReturnNode) Call(scope *Scope) interface{} {
	return n.value.Call(scope)
}
