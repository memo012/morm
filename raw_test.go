package session

import (
	"database/sql"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func TestMain(m *testing.M) {
	db, _ = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/po?charset=utf8mb4")
	code := m.Run()
	_ = db.Close()
	os.Exit(code)
}

func New() *Session {
	return NewSession(db)
}

func TestSession_QueryRow(t *testing.T) {
	s := New()
	var userName string
	var age int
	s = s.Raw("select user_name,age from user where user_name = ?", "迈莫")
	res := s.QueryRow()
	if err := res.Scan(&userName, &age); err != nil {
		t.Fatal("failed to query db", err)
	}
	log.Info("userName--", userName)
	log.Info("age--", age)
}

func TestSession_Exec(t *testing.T) {
	s := New()
	key := "迈莫"
	s = s.Raw("insert into user(user_name, age) values(?, ?)", key, 21)
	_, err := s.Exec()
	if err != nil {
		t.Fatal("failed to insert db", err)
	}
}

func TestSession_Query(t *testing.T) {
	s := New()
	var userName string
	var age int
	s = s.Raw("select user_name, age from user")
	rows, err := s.Query()
	if err != nil {
		t.Fatal("fialed to query db", err)
	}
	for rows.Next() {
		err = rows.Scan(&userName, &age)
		if err != nil {
			t.Fatal("fialed to query db", err)
		}
		log.Info("userName--", userName)
		log.Info("age--", age)
	}
}
