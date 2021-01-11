package session

import (
	"github.com/memo012/morm/log"
	"testing"
)

type Users struct {
	Name string `torm:"name,varchar"`
	Age  int    `torm:"age,int"`
}

func TestClause_InsertStruct(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
		Age:  1,
	}
	clause := NewClause()
	clause = clause.SetTableName("memo").
		InsertStruct(user)
	log.Info(clause.sql)
	log.Info(clause.params)

	//	sql := "INSERT INTO memo (Name,Age) VALUES (?,?)"
}

func TestClause_Condition(t *testing.T) {
	clause := NewClause()
	clause = clause.SetTableName("memo").
		AndEqual("name", "迈莫coding").
		OrEqual("age", 5)
	log.Info(clause.condition)
	log.Info(clause.params)
}

func TestClause_UpdateStruct(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
	}
	clause := NewClause()
	clause = clause.SetTableName("memo").
		UpdateStruct(user)
	log.Info(clause.sqlType[Update])
	log.Info(clause.paramsType[Update])
}