package vm

import "fmt"

type Value float64

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
	fmt.Printf("%g", value)
}
