package field

import "time"

// FieldType 定义字段类型
type FieldType int

const (
	StringType FieldType = iota
	BoolType
	IntType
	Int8Type
	Int16Type
	Int32Type
	Int64Type
	UintType
	Uint8Type
	Uint16Type
	Uint32Type
	Uint64Type
	Float32Type
	Float64Type
	TimeType
	DurationType
	AnyType
)

// Field 定义日志字段接口
type Field interface {
	Key() string
	Value() interface{}
	Type() FieldType
}

// 定义具体的字段类型
type baseField struct {
	key   string
	value interface{}
	typ   FieldType
}

func (f baseField) Key() string        { return f.key }
func (f baseField) Value() interface{} { return f.value }
func (f baseField) Type() FieldType    { return f.typ }

// 定义日志字段类型函数
func String(key string, value string) Field {
	return baseField{key: key, value: value, typ: StringType}
}

func Bool(key string, value bool) Field {
	return baseField{key: key, value: value, typ: BoolType}
}

func Int(key string, value int) Field {
	return baseField{key: key, value: value, typ: IntType}
}

func Int8(key string, value int8) Field {
	return baseField{key: key, value: value, typ: Int8Type}
}

func Int16(key string, value int16) Field {
	return baseField{key: key, value: value, typ: Int16Type}
}

func Int32(key string, value int32) Field {
	return baseField{key: key, value: value, typ: Int32Type}
}

func Int64(key string, value int64) Field {
	return baseField{key: key, value: value, typ: Int64Type}
}

func Uint(key string, value uint) Field {
	return baseField{key: key, value: value, typ: UintType}
}

func Uint8(key string, value uint8) Field {
	return baseField{key: key, value: value, typ: Uint8Type}
}

func Uint16(key string, value uint16) Field {
	return baseField{key: key, value: value, typ: Uint16Type}
}

func Uint32(key string, value uint32) Field {
	return baseField{key: key, value: value, typ: Uint32Type}
}

func Uint64(key string, value uint64) Field {
	return baseField{key: key, value: value, typ: Uint64Type}
}

func Float32(key string, value float32) Field {
	return baseField{key: key, value: value, typ: Float32Type}
}

func Float64(key string, value float64) Field {
	return baseField{key: key, value: value, typ: Float64Type}
}

func Time(key string, value time.Time) Field {
	return baseField{key: key, value: value, typ: TimeType}
}

func Duration(key string, value time.Duration) Field {
	return baseField{key: key, value: value, typ: DurationType}
}

func Any(key string, value interface{}) Field {
	return baseField{key: key, value: value, typ: AnyType}
}
