package compiler

import (
	"errors"
)

type Scanner struct {
	code    []rune
	start   uint
	current uint
	line    int
}

func initScanner(source *string) *Scanner {
	return &Scanner{code: []rune(*source), start: 0, current: 0, line: 0}
}

func (scanner *Scanner) scanToken() Token {
	scanner.skipWhitespace()
	scanner.start = scanner.current

	if scanner.isAtEnd() {
		return scanner.makeToken(TOKEN_EOF)
	}
	c := scanner.advance()
	if isAlpha(c) {
		return scanner.identifier()
	}
	if isDigit(c) {
		return scanner.number()
	}
	switch c {
	case '(':
		return scanner.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return scanner.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return scanner.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return scanner.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return scanner.makeToken(TOKEN_SEMICOLON)
	case ',':
		return scanner.makeToken(TOKEN_COMMA)
	case '.':
		return scanner.makeToken(TOKEN_DOT)
	case '-':
		return scanner.makeToken(TOKEN_MINUS)
	case '+':
		return scanner.makeToken(TOKEN_PLUS)
	case '/':
		return scanner.makeToken(TOKEN_SLASH)
	case '*':
		return scanner.makeToken(TOKEN_STAR)
	case '!':
		if scanner.match('=') {
			return scanner.makeToken(TOKEN_BANG_EQUAL)
		} else {
			return scanner.makeToken(TOKEN_BANG)
		}
	case '=':
		if scanner.match('=') {
			return scanner.makeToken(TOKEN_EQUAL_EQUAL)
		} else {
			return scanner.makeToken(TOKEN_EQUAL)
		}
	case '<':
		if scanner.match('=') {
			return scanner.makeToken(TOKEN_LESS_EQUAL)
		} else {
			scanner.makeToken(TOKEN_LESS)
		}
	case '>':
		if scanner.match('=') {
			return scanner.makeToken(TOKEN_GREATER_EQUAL)
		} else {
			scanner.makeToken(TOKEN_GREATER)
		}
	case '"':
		return scanner.string()

	}

	return scanner.errorToken("Unexpected token")
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= uint(len(scanner.code))
}

func (scanner *Scanner) match(expected rune) bool {
	if scanner.isAtEnd() {
		return false
	}
	if scanner.code[scanner.current] != expected {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) advance() rune {
	c := scanner.code[scanner.current]
	scanner.current++
	return c
}

func (scanner *Scanner) skipWhitespace() {
	for {
		currentChar, err := scanner.peek()
		if err != nil {
			return
		}
		switch currentChar {
		case ' ', '\r', '\t':
			_ = scanner.advance()
		case '\n':
			scanner.line++
			_ = scanner.advance()
		case '/':
			nextChar, err := scanner.peekNext()
			if err != nil {
				return
			}
			if nextChar == '/' {
				for c, _ := scanner.peek(); c != '\n' && !scanner.isAtEnd(); c, _ = scanner.peek() {
					_ = scanner.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (scanner *Scanner) string() Token {
	for c, e := scanner.peek(); e == nil && c != '"' && !scanner.isAtEnd(); c, e = scanner.peek() {
		if c == '\n' {
			scanner.line++
		}
		scanner.advance()
	}
	if scanner.isAtEnd() {
		return scanner.errorToken("Unterminated string.")
	}

	scanner.advance()
	return scanner.makeToken(TOKEN_STRING)
}

func (scanner *Scanner) number() Token {
	for {
		c, err := scanner.peek()
		if err != nil {
			// Reached end
			break
		}
		if !isDigit(c) && c != '.' {
			break
		}
		scanner.advance()
	}
	return scanner.makeToken(TOKEN_NUMBER)
}

func (scanner *Scanner) checkKeyword(start uint, length uint, rest string, tokType TokenType) TokenType {
	if scanner.current-scanner.start == start+length && // Make sure the word is the right length
		string(scanner.code[scanner.start+start:scanner.start+start+length]) == rest { // Make sure it matches rest
		return tokType
	}
	return TOKEN_IDENTIFIER
}

func (scanner *Scanner) identifierType() TokenType {
	switch scanner.code[scanner.start] {
	case 'a':
		return scanner.checkKeyword(1, 2, "nd", TOKEN_AND)
	case 'c':
		return scanner.checkKeyword(1, 4, "lass", TOKEN_CLASS)
	case 'e':
		return scanner.checkKeyword(1, 3, "lse", TOKEN_ELSE)
	case 'f':
		if scanner.current-scanner.start > 1 {
			switch scanner.code[scanner.start+1] {
			case 'a':
				return scanner.checkKeyword(2, 3, "lse", TOKEN_FALSE)
			case 'o':
				return scanner.checkKeyword(2, 1, "r", TOKEN_FOR)
			case 'u':
				return scanner.checkKeyword(2, 1, "n", TOKEN_FUN)
			}
		}
	case 'i':
		return scanner.checkKeyword(1, 1, "f", TOKEN_IF)
	case 'n':
		return scanner.checkKeyword(1, 2, "il", TOKEN_NIL)
	case 'o':
		return scanner.checkKeyword(1, 1, "r", TOKEN_OR)
	case 'p':
		return scanner.checkKeyword(1, 4, "rint", TOKEN_PRINT)
	case 'r':
		return scanner.checkKeyword(1, 5, "eturn", TOKEN_RETURN)
	case 's':
		return scanner.checkKeyword(1, 4, "uper", TOKEN_SUPER)
	case 't':
		if scanner.current-scanner.start > 1 {
			switch scanner.code[scanner.start+1] {
			case 'h':
				return scanner.checkKeyword(2, 2, "is", TOKEN_THIS)
			case 'r':
				return scanner.checkKeyword(2, 2, "ue", TOKEN_TRUE)

			}
		}
	case 'v':
		return scanner.checkKeyword(1, 2, "ar", TOKEN_VAR)
	case 'w':
		return scanner.checkKeyword(1, 4, "hile", TOKEN_WHILE)
	}
	return TOKEN_IDENTIFIER
}

func (scanner *Scanner) identifier() Token {
	for c, e := scanner.peek(); e == nil && (isAlpha(c) || isDigit(c)); c, e = scanner.peek() {
		scanner.advance()
	}
	return scanner.makeToken(scanner.identifierType())
}

func (scanner *Scanner) peek() (rune, error) {
	if scanner.current >= uint(len(scanner.code)) {
		return 0, errors.New("end of code")
	}
	return scanner.code[scanner.current], nil
}

func (scanner *Scanner) peekNext() (rune, error) {
	if scanner.current >= uint(len(scanner.code))-2 {
		return 'x', errors.New("unexpected EOF")
	}
	return scanner.code[scanner.current+1], nil
}

func (scanner *Scanner) makeToken(tokType TokenType) Token {
	return Token{
		tokenType: tokType,
		start:     scanner.start,
		length:    scanner.current - scanner.start,
		line:      scanner.line,
	}
}

func (scanner *Scanner) errorToken(msg string) Token {
	return Token{
		err:       &msg,
		tokenType: TOKEN_ERROR,
		start:     0,
		length:    uint(len(msg)),
		line:      scanner.line,
	}
}

// Other helpers
func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}
