package vm

import (
	"fmt"
)

// region Typing
// ValueType represents the Type of a Value (bool, nil, number)
type ValueType byte

// Enum representing the various data types
const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
)

// endregion Typing

// region Value

// ValueAs represents the data associated with a Value
type ValueAs struct {
	// Represents a Boolean
	boolean bool
	// Represents a numerical value
	number float64
}

// Value represents data in lox
type Value struct {
	// Type of the Value
	typeof ValueType
	// Data associated with a Value
	as ValueAs
}

// endregion Value

// region Value Array

type ValueArrary struct {
	values   []Value
	capacity uint
	count    uint
}

func initValueArray() ValueArrary {
	return ValueArrary{}
}

func writeValueArray(array *ValueArrary, value Value) {
	array.values = append(array.values, value)
	array.count += 1
}

func printValue(value Value) {
	switch value.typeof {
	case VAL_BOOL:
		if valAsBool(value) {
			fmt.Printf("true")
		} else {
			fmt.Printf("false")
		}
	case VAL_NIL:
		fmt.Printf("nil")
	case VAL_NUMBER:
		fmt.Printf("%g", valAsNumber(value))
	}
}

// endregion Value Array

// region Conversions

func boolToVal(boolean bool) Value {
	return Value{
		typeof: VAL_BOOL,
		as: ValueAs{
			boolean: boolean,
		},
	}
}

func nilToVal() Value {
	return Value{
		typeof: VAL_NIL,
	}
}

func numberToVal(number float64) Value {
	return Value{
		typeof: VAL_NUMBER,
		as: ValueAs{
			number: number,
		},
	}
}

func valAsBool(value Value) bool {
	if value.typeof != VAL_BOOL {
		panic("Tried to interpret an invalid value as bool")
	}
	return value.as.boolean
}

func valAsNumber(value Value) float64 {
	if value.typeof != VAL_NUMBER {
		panic("Tried to interpret an invalid value as number")
	}
	return value.as.number
}

func valAsNil() {
	return
}

func isBool(value Value) bool {
	return value.typeof == VAL_BOOL
}

func isNil(value Value) bool {
	return value.typeof == VAL_NIL
}

func isNumber(value Value) bool {
	return value.typeof == VAL_NUMBER
}

func isFalsey(value Value) bool {
	return isNil(value) || (isBool(value) && !valAsBool(value))
}

func valuesEqual(a Value, b Value) bool {
	if a.typeof != b.typeof {
		return false
	}
	switch a.typeof {
	case VAL_BOOL:
		return valAsBool(a) == valAsBool(b)
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return valAsNumber(a) == valAsNumber(b)
	default:
		return false
	}
}

// endregion Conversions
