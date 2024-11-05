package vm

type OpCode byte

// Possible OpCodes
const (
	// Represents a constant value
	OP_CONSTANT OpCode = iota
	// Represents Binary Addition
	OP_ADD
	// Represents Binary Subtraction
	OP_SUBTRACT
	// Represents Binary Multiplication
	OP_MULTIPLY
	// Represents Binary Division
	OP_DIVIDE
	// Represents Unary Negation
	OP_NEGATE
	// Represents a function return
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
