package junklang

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() []Node {
	var nodes []Node

	for !p.isAtEnd() {
		nodes = append(nodes, p.parseStatement())
	}

	return nodes
}

func (p *Parser) parseStatement() Node {
	token := p.advance()

	switch token.Type {
	case JUNK:
		return &PrintNode{value: p.parseExpression()}
	case BUNK:
		name := p.advance().Value
		p.consume(OPERATOR, "=")
		return &DeclareNode{name: name, value: p.parseExpression()}
	case SKUNK:
		return p.parseFunctionDeclaration()
	case DUNK:
		return &ReturnNode{value: p.parseExpression()}
	case EOF:
		return nil
	default:
		if token.Type == IDENT {
			// Er det en funktion?
			if p.peek().Type == OPERATOR && p.peek().Value == "(" {
				p.current-- // Gå tilbage til identifieren
				return p.parseFunctionCall()
			}
		}
		panic(fmt.Sprintf("uventet token: %v", token))
	}
}

func (p *Parser) parseFunctionDeclaration() Node {
	name := p.advance().Value
	p.consume(OPERATOR, "(")

	// Parse parameters
	var params []string
	if p.peek().Type != OPERATOR || p.peek().Value != ")" {
		for {
			params = append(params, p.advance().Value)
			if p.peek().Type == OPERATOR && p.peek().Value == ")" {
				break
			}
			p.consume(OPERATOR, ",")
		}
	}
	p.consume(OPERATOR, ")")

	// Parse function body
	var body []Node
	p.consume(OPERATOR, "{")
	for p.peek().Type != OPERATOR || p.peek().Value != "}" {
		body = append(body, p.parseStatement())
	}
	p.consume(OPERATOR, "}")

	return &FuncDeclareNode{name: name, params: params, body: body}
}

func (p *Parser) parseFunctionCall() Node {
	name := p.advance().Value
	p.consume(OPERATOR, "(")

	var args []Node
	if p.peek().Type != OPERATOR || p.peek().Value != ")" {
		for {
			args = append(args, p.parseExpression())
			if p.peek().Type == OPERATOR && p.peek().Value == ")" {
				break
			}
			p.consume(OPERATOR, ",")
		}
	}
	p.consume(OPERATOR, ")")

	return &FuncCallNode{name: name, args: args}
}

func (p *Parser) parseExpression() Node {
	return p.parseAdditive()
}

func (p *Parser) parseAdditive() Node {
	left := p.parseMultiplicative()

	for p.peek().Type == OPERATOR && (p.peek().Value == "+" || p.peek().Value == "-") {
		operator := p.advance().Value
		right := p.parseMultiplicative()
		left = &OperationNode{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) parseMultiplicative() Node {
	left := p.parsePrimary()

	for p.peek().Type == OPERATOR && (p.peek().Value == "*" || p.peek().Value == "/") {
		operator := p.advance().Value
		right := p.parsePrimary()
		left = &OperationNode{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) parsePrimary() Node {
	token := p.advance()

	switch token.Type {
	case NUM:
		val, _ := strconv.ParseFloat(token.Value, 64)
		return &NumberNode{value: val}
	case IDENT:
		if p.peek().Type == OPERATOR && p.peek().Value == "(" {
			p.current-- // Go back to identifier
			return p.parseFunctionCall()
		}
		return &ReferenceNode{name: token.Value}
	default:
		panic(fmt.Sprintf("Unexpected token: %v", token))
	}
}

func (p *Parser) consume(expectedType TokenType, expectedValue string) Token {
	token := p.advance()
	if token.Type != expectedType || token.Value != expectedValue {
		panic(fmt.Sprintf("Expected %v %v, got %v %v", expectedType, expectedValue, token.Type, token.Value))
	}
	return token
}

func (p *Parser) peek() Token {
	if p.isAtEnd() {
		return Token{Type: EOF}
	}
	return p.tokens[p.current]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}
