package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if len(fileContents) > 0 {
		Scan(string(fileContents))
	} else {
		fmt.Println("EOF  null")
	}

	lox := Lox{}
	if lox.hadError {
		os.Exit(65)
	}
}

type Error struct {
	line    int
	message string
}

func (v Error) report(where string) {
	fmt.Printf("[line %v] Error %v: %v", v.line, where, v.message)
}

type Lox struct {
	hadError bool
}

type TokenType string

const (
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	EOF         TokenType = "EOF"
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   *string // todo: what to make this?
	Line      int
}

func (t Token) toString() string {
	return fmt.Sprintf("%v %v %v", t.TokenType, t.Lexeme, t.Literal)
}

func Scan(fileContents string) {
	s := NewScanner(string(fileContents))
	s.scanTokens()
	for _, t := range s.tokens {
		fmt.Println(t.toString())
	}
}

type Scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []Token
}

func NewScanner(source string) Scanner {
	return Scanner{source: source, tokens: make([]Token, 0)}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{TokenType: EOF, Lexeme: "", Literal: nil, Line: s.line})

}

func (s *Scanner) scanToken() {
	c, ok := s.advance()
	if !ok { // todo: what should we do here?
		return
	}

	switch c {
	case "(":
		s.addToken(LEFT_PAREN)
	case ")":
		s.addToken(RIGHT_PAREN)
	}
	s.current += 1
}

func (s *Scanner) advance() (string, bool) {
	if len(s.source) < s.current {
		return "", false
	}

	return string(s.source[s.current]), true
}

func (s *Scanner) addToken(t TokenType) {
	text := s.source[s.start:s.current+1]
	s.tokens = append(s.tokens, Token{t, text, nil, s.line})
}
