package session

import (
	"reflect"
	"testing"
)

type User struct {
	Name string `torm:"name,varchar"`
	Age  int
}

func TestStructForType(t *testing.T) {
	user := &User{}
	utypes := reflect.TypeOf(user)
	schema := StructForType(utypes.Elem())
	if len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.FieldMap["Name"].Name != "Name" {
		t.Fatal("failed to parse primary key")
	}
}
