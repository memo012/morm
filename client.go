package session

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
)

type Settings struct {
	DriverName     string
	User           string
	Password       string
	Database       string
	Host           string
	Options        map[string]string
	MaxOpenConns   int
	MaxIdleConns   int
	LoggingEnabled bool
}

type Client struct {
	db      *sql.DB
	session *Session
}

func (s *Settings) DataSourceName() string {
	queryString := ""
	for key, value := range s.Options {
		queryString += key + "=" + value + "&"
	}
	ustr := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", s.User, s.Password, s.Host, s.Database, queryString)

	return ustr
}

func NewClient(settings Settings) (c *Client, err error) {
	db, err := sql.Open(settings.DriverName, settings.DataSourceName())
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	c = &Client{db: db}
	c.session = NewSession(db)
	log.Info("Connect database success")
	return
}

func (c *Client) Close() {
	if err := c.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// 新增数据API
func (s *Client) Insert(ctx context.Context, statement *Statement) (int64, error) {
	sql := statement.clause.sql
	vars := statement.clause.params
	result, err := s.session.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// 查询语句API
func (s *Client) FindOne(ctx context.Context, statement *Statement, dest interface{}) (err error) {
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
	rows := s.session.Raw(statement.clause.sql, statement.clause.params...).QueryRow()

	destType := reflect.TypeOf(dest).Elem()
	schema := StructForType(destType)
	// 获取指针指向的元素信息
	destVal := reflect.New(destType).Elem()
	// 结构体字段
	var values []interface{}
	for _, name := range schema.FieldNames {
		values = append(values, destVal.FieldByName(name).Addr().Interface())
	}
	if err := rows.Scan(values...); err != nil {
		log.Info(err)
		return err
	}
	destSlice.Set(destVal)
	return nil
}

func (s *Client) FindAll(ctx context.Context, statement *Statement, dest interface{}) (err error) {
	log.Info(reflect.TypeOf(dest).Kind())
	if reflect.TypeOf(dest).Kind() != reflect.Ptr || reflect.ValueOf(dest).IsNil() {
		return fmt.Errorf("dest is not a ptr or nil")
	}
	destSlice := reflect.ValueOf(dest).Elem()
	destType := destSlice.Type().Elem()

	// 拼接完整SQL语句
	createFindSQL(statement)

	// 进行与数据库交互
	rows, err := s.session.Raw(statement.clause.sql, statement.clause.params...).Query()
	if err != nil {
		return err
	}

	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}

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
	return nil
}

// 删除操作 API
func (s *Client) Delete(ctx context.Context, statement *Statement) (int64, error) {
	createDeleteSQL(statement)
	log.Info(statement.clause.params)
	res, err := s.session.Raw(statement.clause.sql, statement.clause.params...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 	更新操作 API
func (s *Client) Update(ctx context.Context, statement *Statement) (int64, error) {
	createUpdateSQL(statement)
	log.Info(statement.clause.params)
	res, err := s.session.Raw(statement.clause.sql, statement.clause.params...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func createUpdateSQL(statement *Statement) {
	createConditionSQL(statement)
	statement.clause.Build(Update, Where, Condition)
}

func createDeleteSQL(statement *Statement) {
	statement.clause.Set(Delete, statement.clause.tablename)
	createConditionSQL(statement)
	statement.clause.Build(Delete, Where, Condition)
}

func createFindSQL(statement *Statement) {
	statement.clause.Set(Select, statement.clause.cselect, statement.clause.tablename)
	createConditionSQL(statement)
	statement.clause.Build(Select, Where, Condition)
}

// 拼接完整SQL语句
func createConditionSQL(statement *Statement) {
	if statement.clause.condition != "" {
		statement.clause.Set(Where, "where")
		statement.clause.SetCondition(Condition, statement.clause.condition, statement.clause.params)
	}
}
