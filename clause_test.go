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

	//	sql := "INSERT INTO memo (Name,Age) VALUES (?,?)"

	log.Info(clause.sql)
	log.Info(clause.params)
}
