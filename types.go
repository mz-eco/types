package types

import (
	"fmt"
	R "reflect"
)

func Type(x interface{}) R.Type {

	switch n := x.(type) {
	case R.Value:
		return n.Type()
	case R.Type:
		return n
	default:
		return R.TypeOf(x)
	}
}

func ElemType(x interface{}) R.Type {

	switch n := x.(type) {
	case R.Value:
		return n.Type().Elem()
	case R.Type:
		return n.Elem()
	default:
		return Type(x).Elem()
	}

}

func Kind(x interface{}) R.Kind {
	return Type(x).Kind()
}

func IsFunc(x interface{}) bool {
	return Kind(x) == R.Func
}

func is(x R.Type, v interface{}) bool {

	switch n := v.(type) {
	case R.Type:
		return x == v
	case R.Value:
		return x == n.Type()
	case XKind:
		return n.Is(x)
	case R.Kind:
		return x.Kind() == n
	default:
		panic("cccc")
	}
}

type XKind int

const (
	StructPtr XKind = 0
)

var (
	TypeError = ElemType((*error)(nil))
)

func (m XKind) Is(x R.Type) bool {

	if x.Kind() != R.Ptr {
		return false
	}

	return x.Elem().Kind() == R.Struct
}

func Is(x interface{}, v ...interface{}) bool {

	n := Type(x)

	for _, k := range v {
		if is(n, k) {
			return true
		}
	}
	return false
}
func TypeIn(fn interface{}, index int) R.Type {

	var (
		fx = Type(fn)
	)

	if fx.NumIn() <= index {
		panic("func in overflow")
	}

	return fx.In(index)
}

func IsIn(fn interface{}, args ...interface{}) bool {

	var (
		fx = Type(fn)
	)

	if fx.NumIn() != len(args) {
		return false
	}

	for index := 0; index < fx.NumIn(); index++ {

		if !Is(fx.In(index), args[index]) {
			return false
		}
	}

	return true
}

func IsOut(fn interface{}, args ...interface{}) bool {

	var (
		fx = Type(fn)
	)

	if fx.NumOut() != len(args) {
		return false
	}

	for index := 0; index < fx.NumOut(); index++ {

		if !Is(fx.Out(index), args[index]) {
			return false
		}
	}

	return true
}

func New(x interface{}) R.Value {
	return R.New(Type(x))
}

func Value(x interface{}) R.Value {

	switch n := x.(type) {
	case R.Value:
		return n
	case R.Type:
		panic("type could not get value")
	default:
		return R.ValueOf(x)
	}

}

func Values(x ...interface{}) []R.Value {

	var (
		values = make([]R.Value, 0)
	)

	for _, v := range x {
		values = append(values, Value(v))
	}

	return values

}

func Interfaces(x ...R.Value) []interface{} {

	var (
		values = make([]interface{}, 0)
	)

	for _, v := range x {
		values = append(values, v.Interface())
	}

	return values
}

func CallByName(x interface{}, name string, v ...interface{}) []interface{} {

	var (
		this = Value(x)
	)

	if !Is(x, R.Struct, StructPtr) {
		panic(
			fmt.Sprintf("func %s not found.", name))
	}

	m := this.MethodByName(name)

	o := m.Call(Values(v...))

	return Interfaces(o...)
}

func HasField(x interface{}, v interface{}) bool {

	var (
		this  = Type(x)
		check = Type(v)
	)

	if this.Kind() == R.Ptr {
		this = this.Elem()
	}

	if !is(this, R.Struct) {
		panic("HasField check input must give a struct/*struct")
	}

	for index := 0; index < this.NumField(); index++ {

		if is(this.Field(index).Type, check) {
			return true
		}
	}

	return false
}

func FieldTypeByName(x interface{}, name string) (field R.Type) {

	var (
		this = Type(x)
	)

	if this.Kind() == R.Ptr {
		this = this.Elem()
	}

	if !is(this, R.Struct) {
		panic("HasField check input must give a struct/*struct")
	}

	fx, ok := this.FieldByName(name)

	if !ok {
		return
	}

	return fx.Type
}
