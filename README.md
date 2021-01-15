# TORM框架
![TORM标题](https://img-blog.csdnimg.cn/20210115212316910.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxMDY2MDY2,size_16,color_FFFFFF,t_70#pic_center)

# 介绍
**TORM** 为对象-关系映射(Object-Relational Mapping，简称ORM)框架 ，是【七天系列】中的其中一篇关于ORM框架的项目，写这项目的目的主要有以下几点：

- 用最少的代码来实现一款ORM框架
- 通过TORM进而理解ORM实现原理
- 去了解框架设计的奥妙

# 代码实现过程
> 关注【迈莫coding】，查看TORM实现过程文章，代码+文章+视频(后期会录)

# 架构图
![在这里插入图片描述](https://img-blog.csdnimg.cn/20210115223300711.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxMDY2MDY2,size_16,color_FFFFFF,t_70)

# 函数调用图
![在这里插入图片描述](https://img-blog.csdnimg.cn/20210115223214200.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxMDY2MDY2,size_16,color_FFFFFF,t_70)

# 操作手册
- **Insert操作手册**
-  **Delete操作手册**
- **Update操作手册**
- **Find操作手册**

## 代码演示
```go 
package session

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Newclient() (client *Client, err error) {
	setting := Settings{
		DriverName: "mysql",
		User:       "root",
		Password:   "12345678",
		Database:   "po",
		Host:       "127.0.0.1:3306",
		Options:    map[string]string{"charset": "utf8mb4"},
	}
	return NewClient(setting)
}

// 数据库新增使用示例
func TestSession_Insert(t *testing.T) {
	user := &Users{
		Name: "迈莫coding",
		Age:  1,
	}
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		InsertStruct(user)
	client, _ := Newclient()
	client.Insert(context.Background(), statement)
}

// 数据库单条查询示例
func TestSession_FindOne(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("user").
		AndEqual("user_name", "迈莫").
		Select("user_name,age")
	client, err := Newclient()
	if err != nil {
		log.Error(err)
		return
	}
	user := &User{}
	_ = client.FindOne(context.Background(), statement, user)
	log.Info(user)
}
// 数据库多条查询示例
func TestSession_FindAll(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("user").
		Select("user_name,age")
	client, _ := Newclient()
	var user []User
	_ = client.FindAll(context.Background(), statement, &user)
	log.Info(user)
}
// 数据库删除数据示例
func TestSession_Delete(t *testing.T) {
	statement := NewStatement()
	statement = statement.SetTableName("memo").
		AndEqual("name", "迈莫coding")
	client, _ := Newclient()
	client.Delete(context.Background(), statement)
}
// 数据库数据更新示例
func TestSession_Update(t *testing.T) {
	user := &Users{
		Name: "迈莫",
		Age:  1,
	}
	statement := NewStatement()
	statement = statement.SetTableName("user").
		UpdateStruct(user).
		AndEqual("user_name", "迈莫")
	client, _ := Newclient()
	client.Update(context.Background(), statement)
}
