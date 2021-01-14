package session

import (
	"github.com/memo012/morm/log"
	"testing"
)

type Users struct {
	Name string `torm:"user_name,varchar"`
	Age  int    `torm:"age,int"`
}

func TestClause_InsertStruct(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
		Age:  1,
	}
	clause := newClause()
	clause = clause.SetTableName("memo").
		insertStruct(user)
	log.Info(clause.sql)
	log.Info(clause.params)

	//	sql := "INSERT INTO memo (Name,Age) VALUES (?,?)"
}

func TestClause_Condition(t *testing.T) {
	clause := newClause()
	clause = clause.SetTableName("memo").
		andEqual("name", "迈莫coding").
		orEqual("age", 5).
		selectField("name,age")
	log.Info(clause.condition)
	log.Info(clause.params)
	log.Info(clause.cselect)
}

func TestClause_UpdateStruct(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
	}
	clause := newClause()
	clause = clause.SetTableName("memo").
		updateStruct(user)
	log.Info(clause.sqlType[Update])
	log.Info(clause.paramsType[Update])
}