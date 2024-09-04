package entity

// TCVectorSearchMessage ai向量数据库查询数据
type TCVectorSearchMessage struct {
	// 1、通过doc名称查找
	FileName string
	// 2、通过filterSql查找
	FilterSql string
	// 3、通过输入文本embedding后匹配查找
	Content string

	// 通用参数
	// 返回匹配的topN
	Limit int64
}
