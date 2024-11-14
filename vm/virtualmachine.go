package vm

import (
	"fmt"
	"os"
)

const DEBUG_PRINT_CODE bool = false
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

	if !Compile(source, &chunk) {
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

func (machine *VM) binaryOp(f func(Value, Value) Value) InterpretResult {
	if !isNumber(machine.peek(0)) || !isNumber(machine.peek(1)) {
		machine.runtimeError("Operands must be numbers.")
		return INTERPRET_RUNTIME_ERROR
	}

	a := machine.popValue()
	b := machine.popValue()
	machine.pushValue(f(b, a))
	return INTERPRET_OK
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
		case OP_NIL:
			machine.pushValue(nilToVal())
		case OP_TRUE:
			machine.pushValue(boolToVal(true))
		case OP_FALSE:
			machine.pushValue(boolToVal(false))
		case OP_EQUAL:
			a := machine.popValue()
			b := machine.popValue()
			machine.pushValue(boolToVal(valuesEqual(a, b)))
		case OP_GREATER:
			machine.binaryOp(greater)
		case OP_LESS:
			machine.binaryOp(less)
		case OP_ADD:
			res := machine.binaryOp(add)
			if res != INTERPRET_OK {
				return res
			}
		case OP_SUBTRACT:
			res := machine.binaryOp(subtract)
			if res != INTERPRET_OK {
				return res
			}
		case OP_MULTIPLY:
			res := machine.binaryOp(multiply)
			if res != INTERPRET_OK {
				return res
			}
		case OP_DIVIDE:
			res := machine.binaryOp(divide)
			if res != INTERPRET_OK {
				return res
			}
		case OP_NOT:
			machine.pushValue(boolToVal(isFalsey(machine.popValue())))
		case OP_NEGATE:
			if !isNumber(machine.peek(0)) {
				machine.runtimeError("Operand must be a number.")
				return INTERPRET_RUNTIME_ERROR
			}
			machine.pushValue(numberToVal(-valAsNumber(machine.popValue())))
		case OP_RETURN:
			printValue(machine.popValue())
			fmt.Printf("\n")
			return INTERPRET_OK
		}

	}
}

func (machine *VM) peek(position uint) Value {
	return machine.stack[machine.stackTop-1-position]
}

func (machine *VM) runtimeError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	_, _ = os.Stderr.WriteString("\n")

	instruction := machine.ip - 1
	line := machine.chunk.Lines[instruction]
	_, _ = fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
}

// Functions Passed to Binary
func add(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to add non-numbers.")
	}
	return numberToVal(a.as.number + b.as.number)
}
func subtract(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to subtract non-numbers.")
	}
	return numberToVal(a.as.number - b.as.number)
}
func multiply(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to multiply non-numbers.")
	}
	return numberToVal(a.as.number * b.as.number)
}
func divide(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to divide non-numbers.")
	}
	return numberToVal(a.as.number / b.as.number)
}

func less(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to compare non-numbers.")
	}
	return boolToVal(a.as.number < b.as.number)
}

func greater(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to compare non-numbers.")
	}
	return boolToVal(a.as.number > b.as.number)
}
