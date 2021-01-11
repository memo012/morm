package session

type Field struct {
	Name string // 字段名
	Type string // 字段类型
	tablecolumn string // 对应数据库表列名
	Tag string // 额外约束条件
}

type StructSchema struct {
	FieldMap map[string]*Field
	Fields []*Field // 字段详情
}

func DataTypeOf()  {

}


func Parse(dest interface{}) {

}


