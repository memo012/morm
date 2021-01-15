package session

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

type User struct {
	Name string `torm:"user_name,varchar"`
	Age  int    `torm:"age,int"`
}

func TestStructForType(t *testing.T) {
	user := &User{}
	utypes := reflect.TypeOf(user)
	schema := StructForType(utypes.Elem())
	log.Info(schema.FieldNames)
	for i := 0; i < len(schema.Fields); i++ {
		log.Info("字段名称：", schema.Fields[i].Name, ";字段类型:", schema.Fields[i].Type,
			";对应数据库列名:", schema.Fields[i].TableColumn)
	}
	if len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.FieldMap["Name"].Name != "Name" {
		t.Fatal("failed to parse primary key")
	}
}
