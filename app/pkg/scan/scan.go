package scan

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/errors"
)

func Scan(fileContents string) Scanner {
	s := NewScanner(string(fileContents))
	s.scanTokens()
	for _, t := range s.tokens {
		fmt.Fprint(os.Stdout, t.toString())
	}
	for _, e := range s.errors {
		e.Report()
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
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '=':
		s.addToken(s.switchMatch('=', EQUAL_EQUAL, EQUAL), nil)
	case '!':
		s.addToken(s.switchMatch('=', BANG_EQUAL, BANG), nil)
	case '<':
		s.addToken(s.switchMatch('=', LESS_EQUAL, LESS), nil)
	case '>':
		s.addToken(s.switchMatch('=', GREATER_EQUAL, GREATER), nil)
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case ' ', '\r', '\t':
		return
	case '\n':
		s.line++
		return
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
			return
		} else if isAlpha(c) {
			s.identifier()
			return
		}

		s.addError(fmt.Sprintf("Unexpected character: %v", s.source[s.start:s.current]))
	}
}

func (s *Scanner) advance() byte {
	char := s.currentChar()
	s.current += 1
	return char
}

func (s *Scanner) currentChar() byte {
	return s.source[s.current]
}

func (s *Scanner) currentToken() string {
	endingIndex := min(s.current, len(s.source))
	return s.source[s.start:endingIndex]
}

func (s *Scanner) addToken(t TokenType, literal Literal) {
	s.tokens = append(s.tokens, Token{t, s.currentToken(), literal, s.line})
}

func (s *Scanner) addError(message string) {
	s.errors = append(s.errors, errors.NewError(s.line, message))
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if char := s.currentChar(); char != expected {
		return false
	}

	return true
}

func (s *Scanner) switchMatch(expected byte, matched, nonMatched TokenType) TokenType {
	isMatch := s.match(expected)
	if isMatch {
		s.current += 1
		return matched
	}
	return nonMatched
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.currentChar()
}

func (s *Scanner) peekNext() byte {
	if s.current+1 > len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.addError("Unterminated string.")
		return
	}
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, stringLiteral{value})
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	val, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)

	s.addToken(NUMBER, numberLiteral{val})
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	tokenType := IDENTIFIER
	if t, ok := keywords[s.currentToken()]; ok {
		tokenType = t
	}

	s.addToken(tokenType, nil)
}

func isDigit(b byte) bool {
	return (b >= '0') && (b <= '9')
}

func isAlpha(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b == '_')
}

func isAlphaNumeric(b byte) bool {
	return isAlpha(b) || isDigit(b)
}
