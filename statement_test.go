package session

import (
	"github.com/memo012/morm/log"
	"testing"
)

func TestStatement_InsertStruct(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
		Age:  1,
	}
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		InsertStruct(user)
	log.Info(statement.clause.sql)
	log.Info(statement.clause.params)
}

func TestStatement_UpdateStruct(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
	}
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		UpdateStruct(user)
	log.Info(statement.clause.sqlType[Update])
	log.Info(statement.clause.paramsType[Update])
}