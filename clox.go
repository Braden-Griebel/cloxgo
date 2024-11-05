package main

import "github.com/Braden-Griebel/cloxgo/vm"

func main() {
	machine := vm.InitVM()
	chunk := vm.InitChunk()

	constant := vm.AddConstant(&chunk, 1.2)
	vm.WriteChunk(&chunk, vm.OP_CONSTANT, 123)
	vm.WriteChunk(&chunk, vm.OpCode(constant), 123)

	constant = vm.AddConstant(&chunk, 3.4)
	vm.WriteChunk(&chunk, vm.OP_CONSTANT, 123)
	vm.WriteChunk(&chunk, vm.OpCode(constant), 123)

	vm.WriteChunk(&chunk, vm.OP_ADD, 123)

	constant = vm.AddConstant(&chunk, 5.6)
	vm.WriteChunk(&chunk, vm.OP_CONSTANT, 123)
	vm.WriteChunk(&chunk, vm.OpCode(constant), 123)

	vm.WriteChunk(&chunk, vm.OP_DIVIDE, 123)
	vm.WriteChunk(&chunk, vm.OP_NEGATE, 123)

	vm.WriteChunk(&chunk, vm.OP_RETURN, 123)

	vm.DisassembleChunk(&chunk, "Simple Return")
	machine.Interpret(&chunk)
}
