package extractutil

import (
	"regexp"
	"strings"
)

// 提取标题(title)
func ExtractTitle(content string) string {
	reg := regexp.MustCompile("<title>?\\s*([^>]*?)\\s*<?/title>")
	allsubmatch := reg.FindAllStringSubmatch(content, -1)
	ans := ""
	for _, v2 := range allsubmatch {
		for k, v := range v2 {
			if k > 0 {
				ans = ans + v
			}
		}
	}
	return ans
}

// 提取正文(body)
func ExtractBody(content string) string {
	// 去掉head标签
	reg := regexp.MustCompile("<head>([\\s\\S]*?)</head>")
	content = reg.ReplaceAllString(content, "$1")

	// 去掉script中的所有内容，包括script标签
	reg = regexp.MustCompile("<script>?[\\s\\S]*?</script>")
	content = reg.ReplaceAllString(content, "")

	// 去掉style中的所有内容，包括style标签
	reg = regexp.MustCompile("<style>?[\\s\\S]*?</style>")
	content = reg.ReplaceAllString(content, "")

	// 将td换成空格，li 换成 \t,  tr,br,p 换成 \r\n
	reg = regexp.MustCompile("<td[^>]*>")
	content = reg.ReplaceAllString(content, " ")
	rep := strings.NewReplacer("<li>", "\t", "<tr>", "\r\n", "<br>", "\r\n", "<p>", "\r\n")
	content = rep.Replace(content)

	// 去掉所有的成对的尖括号<>
	reg = regexp.MustCompile("<[^>]*>")
	content = reg.ReplaceAllString(content, "")

	// 将&nbsp;等转义字符替换成相应的符号
	rep = strings.NewReplacer("&lt;", "<", "&gt;", ">", "&amp;", "&", "&nbsp;", " ", "&quot;", "\"", "&apos;", "'")
	content = rep.Replace(content)
	reg = regexp.MustCompile("&#.{2,6};")
	content = reg.ReplaceAllString(content, " ")

	// 去掉多余的空行等
	reg = regexp.MustCompile(" +")
	content = reg.ReplaceAllString(content, " ")
	reg = regexp.MustCompile("(\\s*\\t)+")
	content = reg.ReplaceAllString(content, "\t")
	reg = regexp.MustCompile("(\\s*\\r\\n)+")
	content = reg.ReplaceAllString(content, "\r\n")

	return content
}
