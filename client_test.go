package session

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Newclient() (client *Client, err error) {
	setting := Settings{
		DriverName: "mysql",
		User:       "root",
		Password:   "12345678",
		Database:   "po",
		Host:       "127.0.0.1:3306",
		Options:    map[string]string{"charset": "utf8mb4"},
	}
	return NewClient(setting)
}

func TestSession_Insert(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
		Age:  1,
	}
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		InsertStruct(user)
	client, _ := Newclient()
	client.Insert(context.Background(), statement)
}

func TestSession_FindOne(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("user").
		AndEqual("user_name", "迈莫").
		Select("user_name,age")
	client, err := Newclient()
	if err != nil {
		log.Error(err)
		return
	}
	user := &User{}
	_ = client.FindOne(context.Background(), statement, user)
	log.Info(user)
}

func TestSession_FindAll(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("user").
		Select("user_name,age")
	client, _ := Newclient()
	var user []User
	_ = client.FindAll(context.Background(), statement, &user)
	log.Info(user)
}

func TestSession_Delete(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		AndEqual("name", "迈莫coding")
	client, _ := Newclient()
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
	client, _ := Newclient()
	client.Update(context.Background(), statement)
}
