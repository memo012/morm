package session

import (
	"fmt"
	"reflect"
	"strings"
)

type Type int

type Operation int

type Clause struct {
	cselect    string
	cset       string
	tablename  string
	condition  string
	limit      int32
	offset     int32
	sql        string
	params     []interface{}
	sqlType    map[Type]string
	paramsType map[Type][]interface{}
}

const (
	Insert Type = iota
	Value
	Condition
)

const (
	And Operation = iota
	Or
	Lt
	Lte
	Gt
	Gte
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
	c.Set(Insert, c.tablename, schema.FieldNames)
	recordValues := make([]interface{}, 0)
	recordValues = append(recordValues, schema.RecordValues(vars))
	c.Set(Value, recordValues...)
	// SQL语句拼接
	c.Build(Insert, Value)
	return c
}

func (c *Clause) AndEqual(field string, value interface{}) *Clause {
	return c.SetCondition(Condition, "and", field, "=", value)
}

// 通过关键字构建sql语句
func (c *Clause) Set(name Type, param ...interface{}) {
	sql, vars := generators[name](param...)
	c.sqlType[name] = sql
	c.paramsType[name] = vars
}

// 查询条件组装
func (c *Clause) SetCondition(values ...interface{}) *Clause {
	sql, _ := generators[values[0]](values[1:]...)
	c.params = append(c.params, values[4])
	c.addCondition(sql, values[1].(string))
	return c
}

// 条件组成
func (c *Clause) addCondition(sql, opt string) {
	if c.condition == "" {
		c.condition = sql
	} else {
		c.condition = fmt.Sprint("(", c.condition, ")", opt, "(", sql, ")")
	}
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
