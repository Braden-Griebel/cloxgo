package vm

import (
	"fmt"
	"github.com/Braden-Griebel/cloxgo/compiler"
)

const DEBUG_TRACE_EXECUTION bool = false

// Maximum Size of the Stack
const STACK_MAX uint = 256

type VM struct {
	chunk    *Chunk
	ip       uint
	stack    [STACK_MAX]Value
	stackTop uint
}

type InterpretResult byte

const (
	// No Errors
	INTERPRET_OK InterpretResult = iota
	// Error during compilation step
	INTERPRET_COMPILE_ERROR
	// Error during runtime
	INTERPRET_RUNTIME_ERROR
)

func InitVM() VM {
	return VM{}
}

func (machine *VM) Interpret(source string) InterpretResult {
	var chunk Chunk

	if !compiler.Compile(source, &chunk) {
		return INTERPRET_COMPILE_ERROR
	}

	machine.chunk = &chunk
	machine.ip = 0

	result := machine.run()

	return result
}

// Stack Functions
func (machine *VM) pushValue(value Value) {
	machine.stack[machine.stackTop] = value
	machine.stackTop++
}

func (machine *VM) popValue() Value {
	machine.stackTop--
	return machine.stack[machine.stackTop]
}

func (machine *VM) readByte() OpCode {
	instruction := machine.chunk.Code[machine.ip]
	machine.ip += 1
	return instruction
}

func (machine *VM) readConstant() Value {
	return machine.chunk.Constants.values[machine.readByte()]
}

func (machine *VM) binaryOp(f func(Value, Value) Value) {
	a := machine.popValue()
	b := machine.popValue()
	machine.pushValue(f(a, b))

}

func (machine *VM) run() InterpretResult {
	for {
		if DEBUG_TRACE_EXECUTION {
			disassembleInstruction(machine.chunk, machine.ip)
			for slot := uint(0); slot < machine.stackTop; slot++ {
				fmt.Printf("[ ")
				printValue(machine.stack[slot])
				fmt.Printf(" ] ")
			}
			fmt.Printf("\n")
		}
		var instruction OpCode
		instruction = machine.readByte()
		switch instruction {
		case OP_CONSTANT:
			constant := machine.readConstant()
			machine.pushValue(constant)
		case OP_ADD:
			machine.binaryOp(add)
		case OP_SUBTRACT:
			machine.binaryOp(subtract)
		case OP_MULTIPLY:
			machine.binaryOp(multiply)
		case OP_DIVIDE:
			machine.binaryOp(divide)
		case OP_NEGATE:
			machine.pushValue(-machine.popValue())
		case OP_RETURN:
			printValue(machine.popValue())
			fmt.Printf("\n")
			return INTERPRET_OK
		}

	}
}

// Functions Passed to Binary
func add(a Value, b Value) Value {
	return a + b
}
func subtract(a Value, b Value) Value {
	return a - b
}
func multiply(a Value, b Value) Value {
	return a * b
}
func divide(a Value, b Value) Value {
	return a / b
}
