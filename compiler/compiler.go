package compiler

import "fmt"

func Compile(source string) {
	scanner := initScanner(&source)
	line := -1

	for {
		token := scanner.scanToken()
		if token.line != line {
			fmt.Printf("%4d", token.line)
			line = token.line
		} else {
			fmt.Print("   | ")
		}
		fmt.Printf("%2d '%.*s'\n", token.tokenType, token.length, token.start)

		if token.tokenType == TOKEN_EOF {
			break
		}
	}
}

type TokenType byte

const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE

	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS

	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL

	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL

	TOKEN_GREATER
	TOKEN_GREATER_EQUAL

	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE

	TOKEN_FOR
	TOKEN_FUN
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR

	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS

	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_ERROR
	TOKEN_EOF
)

type Token struct {
	tokenType TokenType
	start     uint
	length    uint
	line      int
}

func (scanner *Scanner) scanToken() Token {
	scanner.start = scanner.current

	if scanner.isAtEnd() {
		return makeToken(TOKEN_EOF)
	}

	return errorToken("Unexpected token")
}
