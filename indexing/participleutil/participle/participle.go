// 采用了double array trie
// 可以随时替换成 trie

package participle

import (
	"github.com/gansidui/gose/indexing/participleutil/datrie"
)

type Participle struct {
	forwardTrie  *datrie.DATrie
	backwardTrie *datrie.DATrie
}

// 正向最大匹配构造一个datrie即可，逆向最大匹配刚好相反，在插入前先反转，查询的时候也先反转再查询
func NewParticiple() *Participle {
	p := new(Participle)
	p.forwardTrie = datrie.NewDATrie()
	p.backwardTrie = datrie.NewDATrie()
	return p
}

func ReverseString(src string) string {
	b := []rune(src)
	length := len(b)
	for i := 0; i < length; i++ {
		if i < length-i-1 {
			b[i], b[length-i-1] = b[length-i-1], b[i]
		}
	}
	return string(b)
}

func ReverseStringArray(src []string) {
	length := len(src)
	for i := 0; i < length; i++ {
		if i < length-i-1 {
			src[i], src[length-i-1] = src[length-i-1], src[i]
		}
	}
}

func (this *Participle) Insert(src string) {
	this.forwardTrie.Insert(src)
	// 反转后再插入到backwardTrie中
	this.backwardTrie.Insert(ReverseString(src))
}

func (this *Participle) ForwardFind(src string) bool {
	flag, _ := this.forwardTrie.Find(src)
	return flag
}

func (this *Participle) BackwardFind(src string) bool {
	flag, _ := this.backwardTrie.Find(ReverseString(src))
	return flag
}

func (this *Participle) Num() int {
	return this.forwardTrie.Num()
}

// 正向最大匹配分词，按照dic中的词典将src分词，分词结果以[]string形式返回
func (this *Participle) ForwardMaxMatch(src string) (target []string) {
	return this.forwardTrie.Participle(src)
}

// 逆向最大匹配分词，和正向恰好相反
func (this *Participle) BackwardMaxMatch(src string) (target []string) {
	ss := this.backwardTrie.Participle(ReverseString(src))
	ReverseStringArray(ss)
	for k, v := range ss {
		ss[k] = ReverseString(v)
	}
	return ss
}

// 双向最大匹配，就是进行正向 + 逆向最大匹配
// 如果 正反向分词结果一样，说明没有歧义，就是分词成功
// 如果 正反向结果不一样，说明有歧义，就要处理
// 处理策略：一个词减一分，一个单字减一分，选取分高的一种方法, 若正向和逆向相同，则选取逆向
// 遵循最少切分法原则，目的是选取词少的，单字也少的
func (this *Participle) BidirectionalMatch(src string) (target []string) {
	fs := this.ForwardMaxMatch(src)
	fb := this.BackwardMaxMatch(src)

	fsLen := len(fs)
	fbLen := len(fb)
	same := true
	if fsLen != fbLen {
		same = false
	} else {
		for i := 0; i < fsLen; i++ {
			if fs[i] != fb[i] {
				same = false
				break
			}
		}
	}

	if same {
		return fs
	}

	fsScore := -fsLen
	fbScore := -fbLen

	for i := 0; i < fsLen; i++ {
		if len(fs[i]) == 1 {
			fsScore--
		}
	}
	for i := 0; i < fbLen; i++ {
		if len(fb[i]) == 1 {
			fbScore--
		}
	}

	if fsScore > fbScore {
		return fs
	}
	return fb
}
