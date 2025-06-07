package scan

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/errors"
)

type Scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []data.Token
	errors  []*errors.Error
	scanned bool
}

func NewScanner(source string) Scanner {
	return Scanner{source: source, tokens: make([]data.Token, 0), errors: make([]*errors.Error, 0), line: 1}
}

func (s *Scanner) Scan() ([]data.Token, []*errors.Error) {
	s.scanTokens()
	s.scanned = true
	return s.tokens, s.errors
}

func (s *Scanner) Scanned() bool {
	return s.scanned
}

func (s *Scanner) GetTokens() []data.Token {
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, data.Token{TokenType: data.EOF, Lexeme: "", Literal: nil, Line: s.line})
}

func (s *Scanner) ScanError() bool {
	return len(s.errors) > 0
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(data.LEFT_PAREN, nil)
	case ')':
		s.addToken(data.RIGHT_PAREN, nil)
	case '{':
		s.addToken(data.LEFT_BRACE, nil)
	case '}':
		s.addToken(data.RIGHT_BRACE, nil)
	case ',':
		s.addToken(data.COMMA, nil)
	case '.':
		s.addToken(data.DOT, nil)
	case '-':
		s.addToken(data.MINUS, nil)
	case '+':
		s.addToken(data.PLUS, nil)
	case ';':
		s.addToken(data.SEMICOLON, nil)
	case '*':
		s.addToken(data.STAR, nil)
	case '=':
		s.addToken(s.switchMatch('=', data.EQUAL_EQUAL, data.EQUAL), nil)
	case '!':
		s.addToken(s.switchMatch('=', data.BANG_EQUAL, data.BANG), nil)
	case '<':
		s.addToken(s.switchMatch('=', data.LESS_EQUAL, data.LESS), nil)
	case '>':
		s.addToken(s.switchMatch('=', data.GREATER_EQUAL, data.GREATER), nil)
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.addToken(data.SLASH, nil)
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

func (s *Scanner) addToken(t data.TokenType, literal data.Literal) {
	s.tokens = append(s.tokens, data.Token{t, s.currentToken(), literal, s.line})
}

func (s *Scanner) addError(message string) {
	s.errors = append(s.errors, errors.NewError(s.line, message, ""))
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

func (s *Scanner) switchMatch(expected byte, matched, nonMatched data.TokenType) data.TokenType {
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
	s.addToken(data.STRING, data.StringLiteral{value})
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

	s.addToken(data.NUMBER, data.NumberLiteral{val})
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	tokenType := data.IDENTIFIER
	if t, ok := data.Keywords[s.currentToken()]; ok {
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
