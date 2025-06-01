package scan

import "fmt"

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
	case "{": 
		
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
	text := s.source[s.start : s.current+1]
	s.tokens = append(s.tokens, Token{t, text, nil, s.line})
}
