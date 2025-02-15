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
	strings  map[string]*string
	objects  *Obj
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
	newVM := VM{}
	newVM.strings = make(map[string]*string)
	return newVM
}

func (machine *VM) FreeVM() {
	currentObject := machine.objects
	if currentObject == nil {
		return
	}
	// Loop through the objects freeing them
	// This really isn't necessary due to go's garbage collector
	// but is helpful for learning about garbage collection strategies
	for {
		// Should always be false, just here for safety
		if currentObject == nil {
			break
		}
		// Set the current object data to nil
		// this should drop the reference to the object and
		// free the memory
		// Could also just set the references to the whole object to
		// nil and let the GC collect them, but again, this
		// is just for learning
		currentObject.data = nil
		// Get the next object to work on
		nextObject := currentObject.next
		// Drop the reference to the next object from current object
		currentObject.next = nil
		// If there is no next object exit the loop
		if nextObject == nil {
			break
		}
		// If there is a next object, set that to be current and
		// go about freeing it
		currentObject = nextObject
	}
	// Now that all references within the object chain have been dropped
	machine.objects = nil
	// Empty the strings map
	machine.strings = make(map[string]*string)
	// End of function
	return
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
	// Check if the value being added is an object, if it is,
	// add it to the object linked list
	// Further, if it is a string, add it to the strings table, and intern it
	if isObj(value) {
		newObj := value.data.asObj()
		// If the object is a string, intern it
		if isString(newObj) {
			newString := newObj.data.asString()
			// Check if the newString already has an entry
			newStringPointer, ok := machine.strings[*newString]
			if ok {
				// The string exists in the hash table
				newObj.data = &StringObj{value: newStringPointer}
			} else {
				// The string doesn't yet exist in the hash table
				// Insert it
				machine.strings[*newString] = newString
				// Then set the newObj to point to it
				newObj.data = &StringObj{value: newString}
			}
		}

		newObj.next = machine.objects
		machine.objects = newObj
	}
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
	if (!isNumber(machine.peek(0)) || !isNumber(machine.peek(1))) &&
		(!isObj(machine.peek(0)) && !isObj(machine.peek(1))) {
		machine.runtimeError("Operands must be numbers or strings.")
		return INTERPRET_RUNTIME_ERROR
	}
	// If they are both object, make sure they are both strings
	if !isNumber(machine.peek(0)) || !isNumber(machine.peek(1)) {
		if !isString(machine.peek(0).data.asObj()) || !isString(machine.peek(1).data.asObj()) {
			machine.runtimeError("Operands must be numbers or strings.")
			return INTERPRET_RUNTIME_ERROR
		}
	}

	a := machine.popValue()
	b := machine.popValue()
	machine.pushValue(f(b, a))
	return INTERPRET_OK
}

func (machine *VM) run() InterpretResult {
	for {
		// Check if the ip is beyond the instructions
		if machine.ip >= uint(len(machine.chunk.Code)) {
			return INTERPRET_OK // Reached the end of the instructions
		}

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
		case OP_PRINT:
			printValue(machine.popValue())
			fmt.Print("\n")
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
	if isObj(a) && isObj(b) {
		if !isString(a.data.asObj()) || !isString(b.data.asObj()) {
			panic("Tried to add two objects which are not both strings")
		}
		aString := a.data.asObj().data.asString()
		bString := b.data.asObj().data.asString()
		newString := *aString + *bString
		return objToVal(&newString)

	} else if isNumber(a) && isNumber(b) {
		return numberToVal(a.data.asNumber() + b.data.asNumber())
	}
	panic("Tried to add non-numbers.")

}
func subtract(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to subtract non-numbers.")
	}
	return numberToVal(a.data.asNumber() - b.data.asNumber())
}
func multiply(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to multiply non-numbers.")
	}
	return numberToVal(a.data.asNumber() * b.data.asNumber())
}
func divide(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to divide non-numbers.")
	}
	return numberToVal(a.data.asNumber() / b.data.asNumber())
}

func less(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to compare non-numbers.")
	}
	return boolToVal(a.data.asNumber() < b.data.asNumber())
}

func greater(a Value, b Value) Value {
	if !isNumber(a) || !isNumber(b) {
		panic("Tried to compare non-numbers.")
	}
	return boolToVal(a.data.asNumber() > b.data.asNumber())
}
