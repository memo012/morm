package session

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

var db *sql.DB

func TestMain(m *testing.M) {
	db, _ = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/vblog?charset=utf8mb4")
	code := m.Run()
	_ = db.Close()
	os.Exit(code)
}


func NewSession() *Session {
	return NewClient(db)
}

func TestSession_QueryRow(t *testing.T) {
	s := NewSession()
	var id string
	var visitor int
	s = s.Raw("select id, visitor from web where id = ?", "95843sjdfjl4")
	res := s.QueryRow()
	if err := res.Scan(&id, &visitor); err != nil {
		t.Fatal("failed to query db", err)
	}
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	key := rand.Intn(100)
	fmt.Println(key)
	s = s.Raw("insert into web(id, visitor) values(?, ?)", strconv.Itoa(key), 45)
	_, err := s.Exec()
	if err != nil {
		t.Fatal("failed to insert db", err)
	}
}

func TestSession_Query(t *testing.T) {
	s := NewSession()
	var id string
	var visitor int
	s = s.Raw("select id, visitor from web")
	rows, err := s.Query()
	if err != nil {
		t.Fatal("fialed to query db", err)
	}
	for rows.Next() {
		err = rows.Scan(&id, &visitor)
		if err != nil {
			t.Fatal("fialed to query db", err)
		}
	}
}
