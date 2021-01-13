package session

import (
	"context"
	"fmt"
	"github.com/memo012/morm/log"
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
	destSlice := reflect.Indirect(reflect.ValueOf(dest))
	destValue := reflect.ValueOf(dest).Elem()
	if destValue.Kind() != reflect.Struct {
		return fmt.Errorf("dest is not a struct")
	}

	// 拼接完整SQL语句
	createFindSQL(statement)
	// 进行与数据库交互
	rows, err := s.Raw(statement.clause.sql, statement.clause.params...).Query()
	if err != nil {
		return err
	}

	destType := reflect.TypeOf(dest).Elem()
	schema := StructForType(destType)

	for rows.Next() {
		// 获取指针指向的元素信息
		dest := reflect.New(destType).Elem()
		// 结构体字段
		var values []interface{}
		for _, name := range schema.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		// 赋值
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// 拼接完整SQL语句
func createFindSQL(statement *Statement) {
	statement.clause.Set(Select, statement.clause.cselect, statement.clause.tablename)
	if statement.clause.condition != "" {
		statement.clause.Set(Where, "where")
		log.Info(statement.clause.sql)
		statement.clause.SetCondition(Condition, statement.clause.condition, statement.clause.params)
	}
	statement.clause.Build(Select, Where, Condition)
}
