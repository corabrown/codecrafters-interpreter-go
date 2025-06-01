package scan

import "fmt"

type TokenType string

const (
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
	EOF         TokenType = "EOF"
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   *string // todo: what to make this?
	Line      int
}

func (t Token) toString() string {
	literal := "null"
	if t.Literal != nil {
		literal = *t.Literal
	}
	return fmt.Sprintf("%v %v %v", t.TokenType, t.Lexeme, literal)
}
