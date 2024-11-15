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

// ValueData represents the data associated with a Value
type ValueData interface {
	asBool() bool
	asNumber() float64
	asNil()
}

// region Boolean

type Boolean struct {
	value bool
}

// Implement the ValueData type
func (b *Boolean) asBool() bool {
	return b.value
}

func (b *Boolean) asNumber() float64 {
	panic("Can't coerce a bool to float")
}

func (b *Boolean) asNil() {
	panic("Can't coerce bool to nil")
}

// endregion Boolean

// region Number

type Number struct {
	value float64
}

func (n *Number) asBool() bool {
	panic("Can't coerce number to bool")
}

func (n *Number) asNumber() float64 {
	return n.value
}

func (n *Number) asNil() {
	panic("Can't coerce nil to float")
}

// endregion Number

// region Nil

type Nil struct{}

func (n *Nil) asBool() bool {
	panic("Can't coerce nil to bool")
}

func (n *Nil) asNumber() float64 {
	panic("Can't coerce nil to float")
}

func (n *Nil) asNil() {
	return
}

// endregion Nil

// Value represents data in lox
type Value struct {
	// Type of the Value
	typeof ValueType
	// Data associated with a Value
	data ValueData
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
		data: &Boolean{
			value: boolean,
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
		data: &Number{
			value: number,
		},
	}
}

func valAsBool(value Value) bool {
	if value.typeof != VAL_BOOL {
		panic("Tried to interpret an invalid value data bool")
	}
	return value.data.asBool()
}

func valAsNumber(value Value) float64 {
	if value.typeof != VAL_NUMBER {
		panic("Tried to interpret an invalid value data number")
	}
	return value.data.asNumber()
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
