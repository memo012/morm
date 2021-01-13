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
		OrEqual("age", 2).
		Select("name,age")
	client := NewClient(nil)
	client.FindOne(context.Background(), statement, &User{})
}

func TestSession_Delete(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		AndEqual("name", "迈莫coding")
	client := NewClient(nil)
	client.Delete(context.Background(), statement)
}

func TestSession_Update(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
		Age:  1,
	}
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		UpdateStruct(user).
		AndEqual("name", "迈莫")
	client := NewClient(nil)
	client.Update(context.Background(), statement)
}