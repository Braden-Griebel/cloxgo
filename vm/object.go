package vm

type ObjType byte

const (
	STRING_TYPE ObjType = iota
)

// ObjData represents the data associated with an Obj
type ObjData interface {
	asString() *string
}

// region string
type StringObj struct {
	value *string
}

func (s *StringObj) asString() *string {
	return s.value
}

// endregion string

// Obj represents an object in lox, such as a string, function, etc.
type Obj struct {
	// Type of the Object
	typeof ObjType
	// Data associated with the Object
	data ObjData
	// The next object in the object list
	next *Obj
}

func dataToObj(data interface{}) *Obj {
	var newObj Obj
	switch data.(type) {
	case string:
		stringData := data.(string)
		newObj = Obj{
			typeof: STRING_TYPE,
			data: &StringObj{
				value: &stringData,
			},
		}
	default:
		panic("Unable to create object from data")
	}
	return &newObj
}

func isString(obj *Obj) bool {
	return obj.typeof == STRING_TYPE
}
