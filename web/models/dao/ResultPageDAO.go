package dao

import (
	"fmt"
	"github.com/gansidui/gose/search"
	"github.com/gansidui/gose/web/models/dto"
	"strconv"
	"time"
)

// 根据搜索串和页码得到结果页信息
func GetResultPageInfo(searchString, curPage string) (resultPage dto.ResultPageInfo, success bool) {
	curPageInt, err := strconv.Atoi(curPage)
	if err != nil {
		return resultPage, false
	}

	// 搜索，得到Articles
	start := time.Now()
	var perPageArticlesNum int = 10
	result, total := search.GetSearchResult(searchString, perPageArticlesNum*(curPageInt-1), perPageArticlesNum*curPageInt-1)

	resultPage.Articles = &result
	resultPage.ResultTotal = total
	resultPage.UsedTime = fmt.Sprintf("%v", time.Since(start))

	// 设置Pagination
	var pageTotal int = (total-1)/perPageArticlesNum + 1
	var pagination dto.PaginationInfo
	pagination.PageTotal = pageTotal
	pagination.PerPageArticlesNum = perPageArticlesNum
	pagination.PrevPage = curPageInt - 1
	pagination.NextPage = curPageInt + 1
	pagination.HasPrevPage = false
	pagination.HasNextPage = false
	pagination.Word = searchString

	if curPageInt > 1 {
		pagination.HasPrevPage = true
	}
	if curPageInt < pageTotal {
		pagination.HasNextPage = true
	}
	// 显示最多10个页码链接
	for i, p := 0, curPageInt/10*10+1; i < 10 && p <= pageTotal; i, p = i+1, p+1 {
		pagination.ShowPages = append(pagination.ShowPages, &dto.ShowPageInfo{Word: searchString, Page: p})
	}

	resultPage.Pagination = &pagination

	return resultPage, true
}
