# 七天教程规划

接下来我会将TORM框架的设计原理分为七天来讲述，每天完成其中一个模块，七天形成orm框架基础模型，下面所看到的就是我分为七天的实现功能，如下所示：

- Day01: 序言
- Day02: database/sql基础
- Day03: 对象表结构映射
- Day04: 条件组件库
- Day05: 条件组件API
- Day06: 用户CRUD操作API
- Day07: 支持事务

# TORM框架
![TORM标题](https://img-blog.csdnimg.cn/20210115212316910.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxMDY2MDY2,size_16,color_FFFFFF,t_70#pic_center)


# 谈谈ORM框架
> 对象-关系映射（Object-Relational Mapping，简称ORM），面向对象的开发方法是当今企业级应用开发环境中的主流开发方法，关系数据库是企业级应用环境中永久存放数据的主流数据存储系统。

# 对象-数据库映射关系

|数据库|高级编程语言|
|--|--|
|表(table)|	类(class/struct)|
|记录|	对象|
|字段(field/column)|	对象属性|


# TORM由来

Go 语言中使用比较广泛 ORM 框架是 [gorm](https://github.com/jinzhu/gorm) 和 
[xorm](https://github.com/go-xorm/xorm) 。
自己在实现过程，参考了他们的实现原理，同时也参考了【极客兔兔GeeORM】的实现方式，比如关键词SQL语句的封装和事务实现，就是借鉴【GeeORM】的实现方式。

但TORM框架与他们也有不同点，比如在gorm中单条查询语句的时候，表示结构如下所示：

```go
// 获取第一个匹配记录
db.Where("name = ?", "jinzhu")
```

在TORM中的单条查询API 为：

```go
  statement = statement.SetTableName("user").
      AndEqual("user_name", "迈莫coding").
      Select("user_name,age")
```

在TORM中提供了条件组件库，让用户来进行选择对应的条件进行组装，从而组装成完成的SQL语句，进而的到结果集。



# TORM框架

TORM为对象-关系映射(Object-Relational Mapping，简称ORM)框架 ，是【七天系列】中的其中一篇关于ORM框架的项目，写这项目的目的主要有以下几点：
- 用最少的代码来实现一款ORM框架
- 通过TORM进而理解ORM实现原理
- 去了解框架设计的奥妙


# 代码实现过程
>关注【迈莫coding】，查看TORM实现过程文章，代码+文章+视频

若对Go中反射的使用不了解的话，我写了三篇关于反射的文章，给小伙伴提供参考，足以应对本项目中所使用的反射知识点。

go反射第一弹：[https://mp.weixin.qq.com/s/F8yZyqC5UwoewsX0THqy1w](https://mp.weixin.qq.com/s/F8yZyqC5UwoewsX0THqy1w)  
go反射第二弹：[https://mp.weixin.qq.com/s/lgZykTL8ls6aG0OMNSbZMw](https://mp.weixin.qq.com/s/lgZykTL8ls6aG0OMNSbZMw)  
go反射第三弹：[https://mp.weixin.qq.com/s/vFt06c9herwTrx1LTxNaKg](https://mp.weixin.qq.com/s/vFt06c9herwTrx1LTxNaKg)

# 代码目录

```go
torm
|--raw.go                  // 底层与数据库交互语句
|--raw_test.go
|--schema.go               // 对象表结构映射
|--schema_test.go
|--generators.go           // 关键词sql语句
|--clause.go               // 条件组件库
|--clause_test.go
|--statement.go            // 条件组件库操作API 
|--statement_test.go
|--client.go               // 用户CRUD操作API 
|--client_test.go
|--transaction.go          // 支持事务
|--go.mod
```


# 架构图
![在这里插入图片描述](https://img-blog.csdnimg.cn/20210115223300711.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxMDY2MDY2,size_16,color_FFFFFF,t_70)

# 函数调用图
![在这里插入图片描述](https://img-blog.csdnimg.cn/20210115223214200.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxMDY2MDY2,size_16,color_FFFFFF,t_70)

# 操作手册
- **Insert操作手册**
-  **Delete操作手册**
- **Update操作手册**
- **Find操作手册**


## Insert操作手册

向数据库中新增一条信息

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
```

## FindOne操作手册

通过指定条件查询符合的一条数据

```go
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
```

## FindAll操作手册

```go
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
```

## Delete操作手册

```go
// 数据库删除数据示例
func TestSession_Delete(t *testing.T) {
   statement := NewStatement()
   statement = statement.SetTableName("memo").
      AndEqual("name", "迈莫coding")
   client, _ := Newclient()
   client.Delete(context.Background(), statement)
}
```

## Update操作手册

```go
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
```

> 文章也会持续更新，可以微信搜索「 迈莫coding 」第一时间阅读，回复『1024』领取学习go资料。

![在这里插入图片描述](https://img-blog.csdnimg.cn/20210116225630132.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxMDY2MDY2,size_16,color_FFFFFF,t_70#pic_center)