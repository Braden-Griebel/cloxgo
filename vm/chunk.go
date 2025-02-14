package vm

type OpCode byte

// Possible OpCodes
const (
	// OP_CONSTANT Represents a constant value
	OP_CONSTANT OpCode = iota
	// OP_NIL Represents a Nil Value
	OP_NIL
	// OP_TRUE represents a true value
	OP_TRUE
	// OP_FALSE represents a false value
	OP_FALSE
	// OP_EQUAL represents the equality operator
	OP_EQUAL
	// OP_GREATER represents the greater than operator
	OP_GREATER
	// OP_LESS represents the less than operator
	OP_LESS
	// OP_ADD Represents Binary Addition
	OP_ADD
	// OP_SUBTRACT Represents Binary Subtraction
	OP_SUBTRACT
	// OP_MULTIPLY Represents Binary Multiplication
	OP_MULTIPLY
	// OP_DIVIDE Represents Binary Division
	OP_DIVIDE
	// OP_NOT represents Unary Logical Not
	OP_NOT
	// OP_NEGATE Represents Unary Negation
	OP_NEGATE
	// OP_PRINT Prints the top of the stack
	OP_PRINT
	// OP_RETURN Represents a function return
	OP_RETURN
)

// Chunk is a representation of an array of uint
type Chunk struct {
	Code      []OpCode
	Lines     []uint
	Constants ValueArrary
	Count     uint
}

// Create a new chunk of bytecode
func InitChunk() Chunk {
	c := Chunk{}
	return c
}

// Add a new byte to the chunk
func WriteChunk(chunk *Chunk, codebyte OpCode, line uint) {
	chunk.Code = append(chunk.Code, codebyte)
	chunk.Count += 1
	chunk.Lines = append(chunk.Lines, line)
}

// Add a constant to the constant array
func AddConstant(chunk *Chunk, value Value) byte {
	writeValueArray(&chunk.Constants, value)
	return byte(chunk.Constants.count - 1)
}
