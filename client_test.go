package session

import (
	"context"
	"github.com/memo012/morm/log"
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
	client := New()
	client.Insert(context.Background(), statement)
}

func TestSession_FindOne(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("user").
		AndEqual("user_name", "迈莫").
		Select("user_name,age")
	client := New()
	user := &User{}
	_ = client.FindOne(context.Background(), statement, user)
	log.Info(user)
}

func TestSession_FindAll(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("user").
		Select("user_name,age")
	client := New()
	var user []User
	_ = client.FindAll(context.Background(), statement, &user)
	log.Info(user)
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
		Name: "迈莫",
		Age:  1,
	}
	statement := NewStatement()
	statement = statement.SetTableName("user").
		UpdateStruct(user).
		AndEqual("user_name", "迈莫")
	client := New()
	client.Update(context.Background(), statement)
}
