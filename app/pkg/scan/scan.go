package scan

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/errors"
)

func Scan(fileContents string) Scanner {
	s := NewScanner(string(fileContents))
	s.scanTokens()
	for _, t := range s.tokens {
		fmt.Fprint(os.Stdout, t.toString())
	}
	for _, e := range s.errors {
		e.Report("Unexpected character")
	}

	return s
}

type Scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []Token
	errors  []errors.Error
}

func NewScanner(source string) Scanner {
	return Scanner{source: source, tokens: make([]Token, 0), errors: make([]errors.Error, 0), line: 1}
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

func (s *Scanner) ScanError() bool {
	return len(s.errors) > 0
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case "(":
		s.addToken(LEFT_PAREN)
	case ")":
		s.addToken(RIGHT_PAREN)
	case "{":
		s.addToken(LEFT_BRACE)
	case "}":
		s.addToken(RIGHT_BRACE)
	case ",":
		s.addToken(COMMA)
	case ".":
		s.addToken(DOT)
	case "-":
		s.addToken(MINUS)
	case "+":
		s.addToken(PLUS)
	case ";":
		s.addToken(SEMICOLON)
	case "*":
		s.addToken(STAR)
	case "=":
		s.addToken(s.switchMatch("=", EQUAL_EQUAL, EQUAL))
	case "!":
		s.addToken(s.switchMatch("=", BANG_EQUAL, BANG))
	case "<":
		s.addToken(s.switchMatch("=", LESS_EQUAL, LESS))
	case ">":
		s.addToken(s.switchMatch("=", GREATER_EQUAL, GREATER))
	case "/":
		if s.match("/") {
			for !s.isAtEnd() && s.currentChar() != "\n" {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	default:
		s.addError()
	}
}

func (s *Scanner) advance() string {
	char := s.currentChar()
	s.current += 1
	return char
}

func (s *Scanner) currentChar() string {

	currentChar := string(s.source[s.current])
	return currentChar
}

func (s *Scanner) addToken(t TokenType) {
	endingIndex := min(s.current, len(s.source))
	text := s.source[s.start:endingIndex]
	s.tokens = append(s.tokens, Token{t, text, nil, s.line})
}

func (s *Scanner) addError() {
	message := s.source[s.start:s.current]
	s.errors = append(s.errors, errors.NewError(s.line, message))
}

func (s *Scanner) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}
	if char := s.currentChar(); char != expected {
		return false
	}

	return true
}

func (s *Scanner) switchMatch(expected string, matched, nonMatched TokenType) TokenType {
	isMatch := s.match(expected)
	if isMatch {
		s.current += 1
		return matched
	}
	return nonMatched
}
