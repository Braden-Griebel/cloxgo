package compiler

import (
	"fmt"
	"github.com/Braden-Griebel/cloxgo/vm"
)

type Parser struct {
	current        Token
	previous       Token
	scanner        *Scanner
	compilingChunk *vm.Chunk
	hadError       bool
	panicMode      bool
}

func Compile(source string, chunk *vm.Chunk) bool {
	scanner := initScanner(&source)
	parser := Parser{scanner: scanner, compilingChunk: chunk}
	parser.advance()
	parser.expression()
	parser.consume(TOKEN_EOF, "Expect end of expression.")
	return !parser.hadError
}

func (parser *Parser) advance() {
	parser.previous = parser.current
	for {
		parser.current = parser.scanner.scanToken()
		if parser.current.tokenType != TOKEN_ERROR {
			break
		}

		parser.errorAtCurrent(string(parser.scanner.code[parser.current.start:]))
	}
}

func (parser *Parser) errorAtCurrent(message string) {
	parser.errorAt(&parser.current, message)
}

func (parser *Parser) error(message string) {
	parser.errorAt(&parser.previous, message)
}

func (parser *Parser) errorAt(token *Token, message string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true
	_ = fmt.Errorf("[line %d] Error", token.line)

	if token.tokenType == TOKEN_EOF {
		_ = fmt.Errorf(" at end")
	} else if token.tokenType == TOKEN_ERROR {
		// Pass
	} else {
		_ = fmt.Errorf(" at '%s'", string(parser.scanner.code[token.start:token.start+token.length]))
	}

	_ = fmt.Errorf(": %s\n", message)
	parser.hadError = true
}

func (parser *Parser) consume(tokenType TokenType, message string) {
	if parser.current.tokenType == tokenType {
		parser.advance()
		return
	}

	parser.errorAtCurrent(message)
}

func (parser *Parser) emitByte(instruction vm.OpCode) {
	vm.WriteChunk(currentChunk(), instruction, parser.previous.line)
}

func (parser *Parser) currentChunk() *vm.Chunk {
	return parser.compilingChunk
}
