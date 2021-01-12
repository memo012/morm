package session

import (
	"context"
	"fmt"
	"reflect"
)

// 新增数据API
func (s *Session) Insert(ctx context.Context, statement *Statement) (int64, error) {
	sql := statement.clause.sql
	vars := statement.clause.params
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// 查询语句API
func (s *Session) FindOne(ctx context.Context, statement *Statement, dest interface{}) (err error) {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr || reflect.ValueOf(dest).IsNil() {
		return fmt.Errorf("dest is not a ptr or nil")
	}
	v := reflect.TypeOf(dest).Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("dest is not a struct")
	}
	// 拼接完整SQL语句
	createFindSQL(statement)
	return nil
}

// 拼接完整SQL语句
func createFindSQL(statement *Statement) {
	statement.clause.Set(Select, statement.clause.cselect, statement.clause.tablename)
	if statement.clause.condition != "" {
		statement.clause.Set(Where, "where")
	}
	statement.clause.Build(Select, Where, Condition)
}
