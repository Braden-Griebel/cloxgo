package compiler

import (
	"fmt"
	"github.com/Braden-Griebel/cloxgo/vm"
	"math"
	"strconv"
)

// Parser represents the parser and compiler combined
type Parser struct {
	// Current token
	current Token
	// Previous token
	previous Token
	// Scanner reading the tokens from the soure code
	scanner *Scanner
	// Chunk being compiled
	compilingChunk *vm.Chunk
	// Whether an error occurred during parsing
	hadError bool
	// Whether the Parser/Compiler is in panic mode
	panicMode bool
}

func Compile(source string, chunk *vm.Chunk) bool {
	scanner := initScanner(&source)
	parser := Parser{scanner: scanner, compilingChunk: chunk}
	parser.advance()
	parser.expression()
	parser.consume(TOKEN_EOF, "Expect end of expression.")
	parser.endCompiler()
	return !parser.hadError
}

// region Expression Parsing
func (parser *Parser) expression() {
	parser.parsePrecedence(PREC_ASSIGNMENT)
}

func (parser *Parser) number() {
	value, _ := strconv.ParseFloat(
		string(
			parser.scanner.code[parser.previous.start:parser.previous.start+parser.previous.length]
			), 64,
		)
	parser.emitConstant(vm.Value(value))
}

func (parser *Parser) emitConstant(value vm.Value) {
	parser.emitBytes(vm.OP_CONSTANT, vm.OpCode(parser.makeConstant(value)))
}

func (parser *Parser) makeConstant(value vm.Value) byte {
	constant := vm.AddConstant(parser.currentChunk(), value)
	if constant > math.MaxUint8 {
		parser.error("Too many constants in one chunk.")
		return 0
	}
	return constant
}

func (parser *Parser) grouping() {
	parser.expression()
	parser.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (parser *Parser) unary()  {
	operatorType := parser.previous.tokenType

	parser.parsePrecedence(PREC_UNARY)

	switch operatorType {
	case TOKEN_MINUS:
		parser.emitByte(vm.OP_NEGATE)
	default:
		return
	}
}

type Precedence uint

const (
	PREC_NONE Precedence = iota
	PREC_ASSIGNMENT  // =
	PREC_OR          // or
	PREC_AND         // and
	PREC_EQUALITY    // == !=
	PREC_COMPARISON  // < > <= >=
	PREC_TERM        // + -
	PREC_FACTOR      // * /
	PREC_UNARY       // ! -
	PREC_CALL        // . ()
	PREC_PRIMARY
)

func (parser *Parser) parsePrecedence(precedence Precedence) {}

func (parser *Parser) binary() {
	operatorType:=parser.previous.tokenType
	rule := parser.getRule(operatorType)
	parser.parsePrecedence(Precedence(rule.precedence+1))

	switch operatorType {
	case TOKEN_PLUS:
		parser.emitByte(vm.OP_ADD)
	case TOKEN_MINUS:
		parser.emitByte(vm.OP_SUBTRACT)
	case TOKEN_STAR:
		parser.emitByte(vm.OP_MULTIPLY)
	case TOKEN_SLASH:
		parser.emitByte(vm.OP_DIVIDE)
	default:
		return
	}
}

func (parser *Parser) getRule(operatorType TokenType) ParseRule {

}

type ParseRule struct {
	prefix func()
	infix func()
	precedence Precedence
}

// endregion Expression Parsing

// region Error Handling
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

// endregion Error Handling

// region Helper Functions
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

func (parser *Parser) consume(tokenType TokenType, message string) {
	if parser.current.tokenType == tokenType {
		parser.advance()
		return
	}

	parser.errorAtCurrent(message)
}

func (parser *Parser) emitByte(instruction vm.OpCode) {
	vm.WriteChunk(parser.currentChunk(), instruction, uint(parser.previous.line))
}

func (parser *Parser) emitBytes(instruction1 vm.OpCode, instruction2 vm.OpCode) {
	parser.emitByte(instruction1)
	parser.emitByte(instruction2)
}

func (parser *Parser) currentChunk() *vm.Chunk {
	return parser.compilingChunk
}

func (parser *Parser) endCompiler() {
	parser.emitReturn()
}

func (parser *Parser) emitReturn() {
	parser.emitByte(vm.OP_RETURN)
}

// endregion Helper Functions
