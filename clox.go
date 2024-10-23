package main

func main() {
	chunk := initChunk()

	constant := addConstant(&chunk, 1.2)
	writeChunk(&chunk, OP_CONSTANT, 123)
	writeChunk(&chunk, constant, 123)

	writeChunk(&chunk, OP_RETURN, 123)

	dissasembleChunk(&chunk, "Simple Return")
}
