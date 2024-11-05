package vm

import (
	"fmt"
)

func DisassembleChunk(chunk *Chunk, name string) {
	fmt.Printf("===%s===\n", name)
	var offset uint = 0
	for offset < chunk.Count {
		offset = disassembleInstruction(chunk, offset)
	}
}

func disassembleInstruction(chunk *Chunk, offset uint) uint {
	fmt.Printf("%04d ", offset)

	if offset > 0 && chunk.Lines[offset] == chunk.Lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", chunk.Lines[offset])
	}

	instruction := chunk.Code[offset]

	switch instruction {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", chunk, offset)
	case OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	case OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func simpleInstruction(name string, offset uint) uint {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, chunk *Chunk, offset uint) uint {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.Constants.values[constant])
	fmt.Printf("'\n")
	return offset + 2
}
