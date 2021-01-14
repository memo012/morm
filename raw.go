package session

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"strings"
)

// session类 旨于与数据库进行交互
type Session struct {
	// 数据库引擎
	db     *sql.DB
	tx     *sql.Tx
	// SQL动态参数
	sqlValues []interface{}
	// SQL语句
	sql strings.Builder
}

// CommonDB is a minimal function set of db
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// DB return tx if a tx begins otherwise return *sql.DB
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func NewSession(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlValues = nil
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlValues = append(s.sqlValues, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlValues...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	return s.DB().QueryRow(s.sql.String(), s.sqlValues...)
}

func (s *Session) Query() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlValues...); err != nil {
		log.Error(err)
	}
	return
}
