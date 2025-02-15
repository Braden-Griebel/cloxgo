package vm

import (
	"fmt"
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
	compilingChunk *Chunk
	// Whether an error occurred during parsing
	hadError bool
	// Whether the Parser/Compiler is in panic mode
	panicMode bool
	// Parser Rules
	rules map[TokenType]ParseRule
}

func (parser *Parser) InitRules() {
	parser.rules = map[TokenType]ParseRule{
		TOKEN_LEFT_PAREN:    {parser.grouping, nil, PREC_NONE},
		TOKEN_RIGHT_PAREN:   {nil, nil, PREC_NONE},
		TOKEN_LEFT_BRACE:    {nil, nil, PREC_NONE},
		TOKEN_RIGHT_BRACE:   {nil, nil, PREC_NONE},
		TOKEN_COMMA:         {nil, nil, PREC_NONE},
		TOKEN_DOT:           {nil, nil, PREC_NONE},
		TOKEN_MINUS:         {parser.unary, parser.binary, PREC_TERM},
		TOKEN_PLUS:          {nil, parser.binary, PREC_TERM},
		TOKEN_SEMICOLON:     {nil, nil, PREC_NONE},
		TOKEN_SLASH:         {nil, parser.binary, PREC_FACTOR},
		TOKEN_STAR:          {nil, parser.binary, PREC_FACTOR},
		TOKEN_BANG:          {parser.unary, nil, PREC_NONE},
		TOKEN_BANG_EQUAL:    {nil, nil, PREC_NONE},
		TOKEN_EQUAL:         {nil, nil, PREC_NONE},
		TOKEN_EQUAL_EQUAL:   {nil, parser.binary, PREC_EQUALITY},
		TOKEN_GREATER:       {nil, parser.binary, PREC_COMPARISON},
		TOKEN_GREATER_EQUAL: {nil, parser.binary, PREC_COMPARISON},
		TOKEN_LESS:          {nil, parser.binary, PREC_COMPARISON},
		TOKEN_LESS_EQUAL:    {nil, parser.binary, PREC_COMPARISON},
		TOKEN_IDENTIFIER:    {nil, nil, PREC_NONE},
		TOKEN_STRING:        {parser.string, nil, PREC_NONE},
		TOKEN_NUMBER:        {parser.number, nil, PREC_NONE},
		TOKEN_AND:           {nil, nil, PREC_NONE},
		TOKEN_CLASS:         {nil, nil, PREC_NONE},
		TOKEN_ELSE:          {nil, nil, PREC_NONE},
		TOKEN_FALSE:         {parser.literal, nil, PREC_NONE},
		TOKEN_FOR:           {nil, nil, PREC_NONE},
		TOKEN_FUN:           {nil, nil, PREC_NONE},
		TOKEN_IF:            {nil, nil, PREC_NONE},
		TOKEN_NIL:           {parser.literal, nil, PREC_NONE},
		TOKEN_OR:            {nil, nil, PREC_NONE},
		TOKEN_PRINT:         {nil, nil, PREC_NONE},
		TOKEN_RETURN:        {nil, nil, PREC_NONE},
		TOKEN_SUPER:         {nil, nil, PREC_NONE},
		TOKEN_THIS:          {nil, nil, PREC_NONE},
		TOKEN_TRUE:          {parser.literal, nil, PREC_NONE},
		TOKEN_VAR:           {nil, nil, PREC_NONE},
		TOKEN_WHILE:         {nil, nil, PREC_NONE},
		TOKEN_ERROR:         {nil, nil, PREC_NONE},
		TOKEN_EOF:           {nil, nil, PREC_NONE},
	}
}

func Compile(source string, chunk *Chunk) bool {
	scanner := initScanner(&source)
	parser := Parser{scanner: scanner, compilingChunk: chunk}
	parser.InitRules()
	parser.advance()

	for !parser.match(TOKEN_EOF) {
		parser.declaration()
	}

	//parser.consume(TOKEN_EOF, "Expect end of expression.")
	parser.endCompiler()
	return !parser.hadError
}

// region Declaration Parsing

func (parser *Parser) declaration() {
	parser.statement()
}

// endregion Declaration Parsing

// region Statement Parsing

func (parser *Parser) statement() {
	if parser.match(TOKEN_PRINT) {
		parser.printStatement()
	}
}

func (parser *Parser) printStatement() {
	parser.expression()
	parser.consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	parser.emitByte(OP_PRINT)
}

// endregion Statement Parsing

// region Expression Parsing
func (parser *Parser) expression() {
	parser.parsePrecedence(PREC_ASSIGNMENT)
}

func (parser *Parser) number() {
	value, _ := strconv.ParseFloat(
		string(
			parser.scanner.code[parser.previous.start:parser.previous.start+parser.previous.length],
		), 64,
	)
	parser.emitConstant(numberToVal(value))
}

func (parser *Parser) string() {
	newString := string(parser.scanner.code[parser.previous.start+1 : parser.previous.start+parser.previous.length-1])
	parser.emitConstant(objToVal(&newString))
}

func (parser *Parser) emitConstant(value Value) {
	parser.emitBytes(OP_CONSTANT, OpCode(parser.makeConstant(value)))
}

func (parser *Parser) makeConstant(value Value) byte {
	constant := AddConstant(parser.currentChunk(), value)
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

func (parser *Parser) unary() {
	operatorType := parser.previous.tokenType

	parser.parsePrecedence(PREC_UNARY)

	switch operatorType {
	case TOKEN_BANG:
		parser.emitByte(OP_NOT)
	case TOKEN_MINUS:
		parser.emitByte(OP_NEGATE)
	default:
		return
	}
}

type Precedence uint

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

func (parser *Parser) parsePrecedence(precedence Precedence) {
	parser.advance()
	prefixRule := parser.getRule(parser.previous.tokenType).prefix
	if prefixRule == nil {
		parser.error("Expect expression.")
		return
	}

	prefixRule()

	for precedence <= parser.getRule(parser.current.tokenType).precedence {
		parser.advance()
		infixRule := parser.getRule(parser.previous.tokenType).infix
		infixRule()
	}
}

func (parser *Parser) binary() {
	operatorType := parser.previous.tokenType
	rule := parser.getRule(operatorType)
	parser.parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case TOKEN_BANG_EQUAL:
		parser.emitBytes(OP_EQUAL, OP_NOT)
	case TOKEN_EQUAL_EQUAL:
		parser.emitByte(OP_EQUAL)
	case TOKEN_GREATER:
		parser.emitByte(OP_GREATER)
	case TOKEN_GREATER_EQUAL:
		parser.emitBytes(OP_LESS, OP_NOT)
	case TOKEN_LESS:
		parser.emitByte(OP_LESS)
	case TOKEN_LESS_EQUAL:
		parser.emitBytes(OP_GREATER, OP_NOT)
	case TOKEN_PLUS:
		parser.emitByte(OP_ADD)
	case TOKEN_MINUS:
		parser.emitByte(OP_SUBTRACT)
	case TOKEN_STAR:
		parser.emitByte(OP_MULTIPLY)
	case TOKEN_SLASH:
		parser.emitByte(OP_DIVIDE)
	default:
		return
	}
}

func (parser *Parser) literal() {
	switch parser.previous.tokenType {
	case TOKEN_FALSE:
		parser.emitByte(OP_FALSE)
	case TOKEN_TRUE:
		parser.emitByte(OP_TRUE)
	case TOKEN_NIL:
		parser.emitByte(OP_NIL)
	default:
		return
	}
}

func (parser *Parser) getRule(operatorType TokenType) ParseRule {
	return parser.rules[operatorType]
}

type ParseRule struct {
	prefix     func()
	infix      func()
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

func (parser *Parser) emitByte(instruction OpCode) {
	WriteChunk(parser.currentChunk(), instruction, uint(parser.previous.line))
}

func (parser *Parser) emitBytes(instruction1 OpCode, instruction2 OpCode) {
	parser.emitByte(instruction1)
	parser.emitByte(instruction2)
}

func (parser *Parser) currentChunk() *Chunk {
	return parser.compilingChunk
}

func (parser *Parser) endCompiler() {
	if DEBUG_PRINT_CODE {
		if !parser.hadError {
			DisassembleChunk(parser.currentChunk(), "code")
		}
	}
}

func (parser *Parser) emitReturn() {
	parser.emitByte(OP_RETURN)
}

func (parser *Parser) match(tokenType TokenType) bool {
	if !parser.check(tokenType) {
		return false
	}
	parser.advance()
	return true
}

func (parser *Parser) check(tokenType TokenType) bool {
	return parser.current.tokenType == tokenType
}

// endregion Helper Functions
