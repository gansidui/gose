package participleutil

import (
	"bufio"
	"github.com/gansidui/gose/indexing/participleutil/participle"
	"os"
)

var p *participle.Participle

func init() {
	p = participle.NewParticiple()
}

// 加载词库
func LoadDic(dicPath string) {
	file, err := os.Open(dicPath)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	re := bufio.NewReader(file)

	for {
		line, _, err := re.ReadLine()
		if err != nil {
			break
		}
		p.Insert(string(line))
	}
}

// 分词
func Participle(src string) []string {
	return p.BidirectionalMatch(src)
}
