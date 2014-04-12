package trie

type trieNode struct {
	child map[rune]*trieNode
	flag  bool
}

func newTrieNode() *trieNode {
	t := new(trieNode)
	t.flag = false
	t.child = make(map[rune]*trieNode, 0)
	return t
}

type Trie struct {
	root *trieNode
	num  int
}

func NewTrie() *Trie {
	return &Trie{root: newTrieNode(), num: 0}
}

func (this *Trie) Insert(src string) {
	curNode := this.root
	for _, v := range src {
		if curNode.child[v] == nil {
			newNode := newTrieNode()
			curNode.child[v] = newNode
		}
		curNode = curNode.child[v]
	}
	curNode.flag = true
	this.num++
}

// 若不存在src,则flag为false，且preWordLastIndex保存该路径上离失配地点最近的一个词的最后一个rune的末位置
// 若存在src,则flag为true，且应该忽视preWordLastIndex
func (this *Trie) Find(src string) (flag bool, preWordLastIndex int) {
	curNode := this.root
	ff := false
	for k, v := range src {
		if ff {
			preWordLastIndex = k
			ff = false
		}
		if curNode.child[v] == nil {
			return false, preWordLastIndex
		}
		curNode = curNode.child[v]
		if curNode.flag {
			ff = true
		}
	}
	return curNode.flag, preWordLastIndex
}

func (this *Trie) Num() int {
	return this.num
}

// 正向最大匹配分词，按照trie中的词典将src分词，分词结果以[]string形式返回
func (this *Trie) Participle(src string) (target []string) {
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
