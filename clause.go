package session

import (
	"fmt"
	"reflect"
	"strings"
)

type Type int

type Operation int

// 条件组装 不外显 便于与用户API分层隔离

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
	Update
	Delete
	Limit
	Condition
)

// NewClause 初始化
func newClause() *Clause {
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
func (c *Clause) insertStruct(vars interface{}) *Clause {
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
	c.Build(Insert, Value)
	return c
}

//
func (c *Clause) updateStruct(vars interface{}) *Clause {
	types := reflect.TypeOf(vars)
	if types.Kind() == reflect.Ptr {
		types = types.Elem()
	}
	if types.Kind() != reflect.Struct {
		return c
	}
	// 数据映射
	schema := StructForType(types)
	m := make(map[string]interface{})
	m = schema.UpdateParam(vars)
	// 构建SQL语句
	c.Set(Update, c.tablename, m)
	return c
}

func (c *Clause) andEqual(field string, value interface{}) *Clause {
	return c.setCondition(Condition, "AND", field, "=", value)
}

func (c *Clause) orEqual(field string, value interface{}) *Clause {
	return c.setCondition(Condition, "OR", field, "=", value)
}

// 通过关键字构建sql语句
func (c *Clause) Set(name Type, param ...interface{}) {
	sql, vars := generators[name](param...)
	c.sqlType[name] = sql
	c.paramsType[name] = vars
}

// 查询条件组装
func (c *Clause) setCondition(values ...interface{}) *Clause {
	sql, vars := generators[values[0].(Type)](values[2:]...)
	c.params = append(c.params, vars...)
	c.addCondition(sql, values[1].(string))
	return c
}

// 条件组成
func (c *Clause) addCondition(sql, opt string) {
	if c.condition == "" {
		c.condition = sql
	} else {
		c.condition = fmt.Sprint("(", c.condition, ") ", opt, " (", sql, ")")
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
