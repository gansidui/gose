package participleutil

import (
	"fmt"
	"testing"
)

func TestParticiple(t *testing.T) {
	
	LoadDic("./mydic.txt")

	ss := Participle("学历史学好")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()

	ss = Participle("搜噶，我爱豆豆猪")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()
}
