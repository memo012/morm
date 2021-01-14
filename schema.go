package session

import (
	"go/ast"
	"reflect"
	"strings"
	"sync"
)

//struct 标签解析结果
type Filed struct {
	Name        string // 字段名
	i           int    // 位置
	Type        string // 字段类型
	TableColumn string // 对应数据库表列名
	Tag         string // 约束条件
}

// Schema 主要包含被映射的字段(Fields)
type Schema struct {
	Fields     []*Filed // 字段属性组合
	FieldNames []string // 字段名称
	FieldMap   map[string]*Filed // key:字段名  value：字段属性
}

var (
	structMutex sync.RWMutex
	structCache = make(map[reflect.Type]*Schema)
)

func StructForType(t reflect.Type) *Schema {
	structMutex.RLock()
	st, found := structCache[t]
	structMutex.RUnlock()
	if found {
		return st
	}

	structMutex.Lock()
	defer structMutex.Unlock()
	st, found = structCache[t]
	if found {
		return st
	}

	st = &Schema{FieldMap: make(map[string]*Filed)}
	dataTypeOf(t, st)
	structCache[t] = st
	return st
}

func dataTypeOf(types reflect.Type, schema *Schema) {
	// 遍历所有字段
	for i := 0; i < types.NumField(); i++ {
		p := types.Field(i)
		// 忽略匿名字段和私有字段
		if p.Anonymous || !ast.IsExported(p.Name) {
			continue
		}
		field := &Filed{
			Name: p.Name,
			i:    i,
		}
		var tag = field.Name
		field.TableColumn = field.Name
		if tg, ok := p.Tag.Lookup("torm"); ok {
			tag = tg
		}

		tagArr := strings.Split(tag, ",")
		if len(tagArr) > 0 {
			if tagArr[0] == "-" {
				continue
			}
			if len(tagArr[0]) > 0 {
				field.TableColumn = tagArr[0]
			}
			if len(tagArr) > 1 && len(tagArr[1]) > 0 {
				field.Type = tagArr[1]
			}
		}

		schema.Fields = append(schema.Fields, field)
		schema.FieldMap[p.Name] = field
		schema.FieldNames = append(schema.FieldNames, p.Name)
	}
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

func (s *Schema) UpdateParam(dest interface{}) map[string]interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	m := make(map[string]interface{})
	for _, field := range s.Fields {
		m[field.TableColumn] =  destValue.FieldByName(field.Name).Interface()
	}
	return m
}


