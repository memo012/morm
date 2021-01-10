package session

import (
	"reflect"
	"strings"
)

type Type int

type Clause struct {
	cselect    string
	cset       string
	tablename  string
	limit      int32
	offset     int32
	sql        string
	params     []interface{}
	sqlType    map[Type]string
	paramsType map[Type][]interface{}
}

const (
	INSERT Type = iota
	VALUES
)

// NewClause 初始化
func NewClause() *Clause {
	return &Clause{
		cselect:    "*",
		limit:      -1,
		offset:     -1,
		sqlType:    make(map[Type]string),
		paramsType: make(map[Type][]interface{}),
	}
}

// SetTableName 设置表名
func (c *Clause) SetTableName(tableName string) *Clause {
	c.tablename = tableName
	return c
}

//
func (c *Clause) InsertStruct(vars interface{}) *Clause {
	types := reflect.TypeOf(vars)
	if types.Kind() == reflect.Ptr {
		types = types.Elem()
	}
	if types.Kind() != reflect.Struct {
		return c
	}
	// 数据映射
	schema := StructForType(types)
	// 构建SQL语句
	c.Set(INSERT, c.tablename, schema.FieldNames)
	recordValues := make([]interface{}, 0)
	recordValues = append(recordValues, schema.RecordValues(vars))
	c.Set(VALUES, recordValues...)
	// SQL语句拼接
	c.Build(INSERT, VALUES)
	return c
}

// 通过关键字构建sql语句
func (c *Clause) Set(name Type, param ...interface{}) {
	sql, vars := generators[name](param...)
	c.sqlType[name] = sql
	c.paramsType[name] = vars
}

// 拼接各个SQL语句
func (c *Clause) Build(orders ...Type) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sqlType[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.paramsType[order]...)
		}
	}
	c.sql = strings.Join(sqls, " ")
	c.params = vars
}
