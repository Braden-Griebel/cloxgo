package main

// Possible uints
const (
	// Represents a constant value
	OP_CONSTANT uint = iota
	// Represents a function return
	OP_RETURN uint = iota
)

// Chunk is a representation of an array of uint
type Chunk struct {
	Code      []uint
	Lines     []uint
	Constants ValueArrary
	Count     uint
}

// Create a new chunk of bytecode
func initChunk() Chunk {
	c := Chunk{}
	return c
}

// Add a new byte to the chunk
func writeChunk(chunk *Chunk, codebyte uint, line uint) {
	chunk.Code = append(chunk.Code, codebyte)
	chunk.Count += 1
	chunk.Lines = append(chunk.Lines, line)
}

// Add a constant to the constant array
func addConstant(chunk *Chunk, value Value) uint {
	writeValueArray(&chunk.Constants, value)
	return chunk.Constants.count - 1
}
