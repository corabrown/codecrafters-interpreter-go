package scan

import "fmt"

type TokenType string

const (
	// single character tokens
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	STAR        TokenType = "STAR"
	SLASH       TokenType = "SLASH"

	// one or two character tokens
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"

	//literals
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"
	IDENTIFIER TokenType = "IDENTIFIER"

	// keywords
	AND    TokenType = "AND"
	CLASS  TokenType = "CLASS"
	ELSE   TokenType = "ELSE"
	FALSE  TokenType = "FALSE"
	FUN    TokenType = "FUN"
	FOR    TokenType = "FOR"
	IF     TokenType = "IF"
	NIL    TokenType = "NIL"
	OR     TokenType = "OR"
	PRINT  TokenType = "PRINT"
	RETURN TokenType = "RETURN"
	SUPER  TokenType = "SUPER"
	THIS   TokenType = "THIS"
	TRUE   TokenType = "TRUE"
	VAR    TokenType = "VAR"
	WHILE  TokenType = "WHILE"

	EOF TokenType = "EOF"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   Literal // todo: what to make this?
	Line      int
}

func (t Token) ToString() string {
	literal := "null"
	if t.Literal != nil {
		literal = t.Literal.ToString()
	}
	return fmt.Sprintf("%v %v %v\n", t.TokenType, t.Lexeme, literal)
}

type Literal interface {
	ToString() string
}

type StringLiteral struct {
	Val string
}

func (s StringLiteral) ToString() string {
	return s.Val
}

type NumberLiteral struct {
	Val float64
}

func (n NumberLiteral) ToString() string {
	if n.Val == float64(int(n.Val)) {
		return fmt.Sprintf("%.1f", n.Val)
	}
	return fmt.Sprintf("%v", n.Val)
}

type BooleanLiteral struct {
	Val bool
}

func (b BooleanLiteral) ToString() string {
	return fmt.Sprintf("%v", b.Val)
}

type NullLiteral struct {}

func (n NullLiteral) ToString() string {
	return "nil"
}