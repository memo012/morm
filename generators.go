
package session
import (
	"fmt"
	"strings"
)
// 该类旨于 完成关键字SQL语句构建工作
type generator func(values ...interface{}) (string, []interface{})
var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[Insert] = _insert
	generators[Value] = _values
	generators[Update] = _update
	generators[Condition] = _condition
	generators[Delete] = _delete
	generators[Limit] = _condition
	generators[Select] = _select
	generators[Where] = _where
}

// insert关键词
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

// values关键词
func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(",")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

// 查询条件组装
func _condition(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	sql.WriteString("`")
	sql.WriteString(values[0].(string))
	sql.WriteString("`")
	sql.WriteString(values[1].(string))
	sql.WriteString("?")
	return sql.String(), []interface{}{values[2]}
}

// update关键词
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+"=?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ",")), vars
}

// delete关键词
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

// limit关键词
func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

// select关键词
func _select(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("select %s from %s", values[0], values[1]), []interface{}{}
}

// where关键词
func _where(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("%s","WHERE"), []interface{}{}
}