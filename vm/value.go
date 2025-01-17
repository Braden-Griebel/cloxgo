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
	VAL_OBJ
)

// endregion Typing

// region Value

// ValueData represents the data associated with a Value
type ValueData interface {
	asBool() bool
	asNumber() float64
	asNil()
	asObj() *Obj
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

func (b *Boolean) asObj() *Obj {
	panic("Can't coerce bool to obj")
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

func (n *Number) asObj() *Obj {
	panic("Can't coerce number to obj")
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

func (n *Nil) asObj() *Obj {
	panic("Can't coerce obj to nil")
}

// endregion Nil

// region Object

type Object struct {
	value *Obj
}

func (o *Object) asBool() bool {
	panic("Can't coerce obj to bool")
}

func (o *Object) asNumber() float64 {
	panic("Can't coerce obj to float")
}

func (o *Object) asNil() {
	panic("Can't coerce nil to float")
}

func (o *Object) asObj() *Obj {
	return o.value
}

// endregion Object

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
	case VAL_OBJ:
		printObject(value)
	}
}

func printObject(value Value) {
	object := valAsObj(value)
	switch object.typeof {
	case STRING_TYPE:
		fmt.Printf(*object.data.asString())
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

func objToVal(obj interface{}) Value {
	return Value{
		typeof: VAL_OBJ,
		data: &Object{
			value: dataToObj(obj),
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

func valAsNil(value Value) Nil {
	return Nil{}
}

func valAsObj(value Value) *Obj {
	if value.typeof != VAL_OBJ {
		panic("Tried to interpret an invalid value data object")
	}
	return value.data.asObj()
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

func isObj(value Value) bool {
	return value.typeof == VAL_OBJ
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
	case VAL_OBJ:
		aObj := a.data.asObj()
		bObj := b.data.asObj()
		if aObj.typeof != bObj.typeof {
			return false
		}
		if isString(aObj) && isString(bObj) {
			aString := valAsObj(a).data.asString()
			bString := valAsObj(b).data.asString()
			return *aString == *bString
		}
		// No other objects implemented yet, so just return false
		return false
	default:
		return false
	}
}

// endregion Conversions
