package dao

import (
	"fmt"
	"github.com/gansidui/gose/search"
	"github.com/gansidui/gose/web/models/dto"
	"strconv"
	"time"
)

// 根据搜索串q, 起始位置start，每页显示数量num 得到结果页信息
func GetResultPageInfo(q, start, num string) (resultPage dto.ResultPageInfo, success bool) {
	startInt, err := strconv.Atoi(start)
	if err != nil {
		return resultPage, false
	}
	numInt, err := strconv.Atoi(num)
	if err != nil {
		return resultPage, false
	}

	// 搜索，得到Articles
	startTime := time.Now()
	result, total := search.GetSearchResult(q, startInt, startInt+numInt-1)

	resultPage.Articles = &result
	resultPage.ResultTotal = total
	resultPage.UsedTime = fmt.Sprintf("%v", time.Since(startTime))

	// 设置Pagination
	var pageTotal int = (total-1)/numInt + 1
	var curPageInt int = startInt/numInt + 1
	var pagination dto.PaginationInfo
	pagination.PageTotal = pageTotal
	pagination.PerPageArticlesNum = numInt
	pagination.PrevPageStart = startInt - numInt
	pagination.NextPageStart = startInt + numInt
	pagination.HasPrevPage = false
	pagination.HasNextPage = false
	pagination.QueryString = q

	if curPageInt > 1 {
		pagination.HasPrevPage = true
	}
	if curPageInt < pageTotal {
		pagination.HasNextPage = true
	}
	// 最多显示num个页码链接
	for i, p := 0, (curPageInt-1)/numInt*numInt+1; i < numInt && p <= pageTotal; i, p = i+1, p+1 {
		pagination.ShowPages = append(pagination.ShowPages, &dto.ShowPageInfo{QueryString: q, Page: p, Start: (p - 1) * numInt})
	}

	resultPage.Pagination = &pagination

	return resultPage, true
}
