package dto

import (
	"github.com/gansidui/gose/search"
)

// 页面链接信息
type ShowPageInfo struct {
	QueryString string // 搜索串
	Page        int    // 当前链接的页码
	Start       int    // 当前页的结果起始位置
}

// 页码信息
type PaginationInfo struct {
	PageTotal          int             // 总共有多少页
	PerPageArticlesNum int             // 每页多少篇文章
	PrevPageStart      int             // 前一页起始位置
	NextPageStart      int             // 后一页起始位置
	HasPrevPage        bool            // 是否有上一页
	ShowPages          []*ShowPageInfo // 中间显示哪些页码链接
	HasNextPage        bool            // 是否有下一页
	QueryString        string          // 搜索串
}

// 文章信息, 在 search 包中
// type ArticleInfo struct {
// 	Title   string // 标题
// 	Summary string // 摘要
// 	Url     string // 原始url
// 	Path    string // 本地路径
// }

// 结果页信息
type ResultPageInfo struct {
	Articles    *[]search.ArticleInfo
	Pagination  *PaginationInfo
	ResultTotal int
	UsedTime    string
}
