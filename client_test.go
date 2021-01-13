package session

import (
	"context"
	"testing"
)

func TestSession_Insert(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
		Age:  1,
	}
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		InsertStruct(user)
	client := NewClient(nil)
	client.Insert(context.Background(), statement)
}

func TestSession_FindOne(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		AndEqual("name", "迈莫coding").
		OrEqual("age", 2)
	client := NewClient(nil)
	client.FindOne(context.Background(), statement, &User{})
}