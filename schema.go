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

// 对象与表结构转换
func StructForType(t reflect.Type) *Schema {
	// step1: 缓存获取
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
	// step2: 对象关系映射
	st = &Schema{FieldMap: make(map[string]*Filed)}
	dataTypeOf(t, st)

	// step3: 缓存
	structCache[t] = st
	return st
}

// 对象与表结构转换(实际工作函数)
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
		// 获得额外约束条件
		tagArr := strings.Split(tag, ",")
		if len(tagArr) > 0 {
			if tagArr[0] == "-" {
				continue
			}
			// 数据库中对应列表名称
			if len(tagArr[0]) > 0 {
				field.TableColumn = tagArr[0]
			}
			// 数据库中对应列表类型
			if len(tagArr) > 1 && len(tagArr[1]) > 0 {
				field.Type = tagArr[1]
			}
		}
		// 存储所有字段信息
		schema.Fields = append(schema.Fields, field)
		schema.FieldMap[p.Name] = field
		// 存储所有字段名称
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


