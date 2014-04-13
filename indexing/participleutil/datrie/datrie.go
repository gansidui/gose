// double array trie
// base(s) + c --> t

package datrie

type trieNode struct {
	c    rune
	flag bool
}

func newTrieNode() *trieNode {
	return &trieNode{
		c:    0,
		flag: false,
	}
}

type pair struct {
	s *trieNode
	c rune
}

type DATrie struct {
	root *trieNode
	darr map[pair]*trieNode
	num  int
}

func NewDATrie() *DATrie {
	return &DATrie{
		root: newTrieNode(),
		darr: make(map[pair]*trieNode),
		num:  0,
	}
}

func (this *DATrie) Insert(src string) {
	curNode := this.root
	for _, v := range src {
		p := pair{s: curNode, c: v}
		if this.darr[p] == nil {
			newNode := newTrieNode()
			this.darr[p] = newNode
		}
		curNode = this.darr[p]
	}
	curNode.flag = true
	this.num++
}

// 若不存在src,则flag为false，且preWordLastIndex保存该路径上离失配地点最近的一个词的最后一个rune的末位置
// 若存在src,则flag为true，且应该忽视preWordLastIndex
func (this *DATrie) Find(src string) (flag bool, preWordLastIndex int) {
	curNode := this.root
	ff := false
	for k, v := range src {
		if ff {
			preWordLastIndex = k
			ff = false
		}
		p := pair{s: curNode, c: v}
		if this.darr[p] == nil {
			return false, preWordLastIndex
		}
		curNode = this.darr[p]
		if curNode.flag {
			ff = true
		}
	}
	return curNode.flag, preWordLastIndex
}

func (this *DATrie) Num() int {
	return this.num
}

// 正向最大匹配分词，按照词典将src分词，分词结果以[]string形式返回
func (this *DATrie) Participle(src string) (target []string) {
	if len(src) == 0 {
		return
	}

	flag, preWordLastIndex, length := false, 0, len(src)
	left, right := 0, length

	for left < right {
		flag, preWordLastIndex = this.Find(src[left:right])
		preWordLastIndex += left
		if flag {
			target = append(target, src[left:right])
			left = right
			right = length
		} else {
			if preWordLastIndex == left {
				left++ // 多个字节的rune一定会到这里多次 :)
			} else {
				right = preWordLastIndex
			}
		}
	}
	return
}
