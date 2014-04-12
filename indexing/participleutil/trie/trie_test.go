package trie

import (
	"fmt"
	"testing"
)

func TestFind(t *testing.T) {

	tr := NewTrie()
	tr.Insert("ab")
	tr.Insert("cd")
	tr.Insert("abcd")
	tr.Insert("abcde")

	if tr.Num() != 4 {
		t.Error("Error")
	}

	flag, preWordLastflag := tr.Find("ab")
	if !flag {
		t.Error("Error")
	}

	flag, preWordLastflag = tr.Find("cde")
	if flag || preWordLastflag != 2 {
		t.Error("Error")
	}

	flag, preWordLastflag = tr.Find("abcdf")
	if flag || preWordLastflag != 4 {
		t.Error("Error")
	}

	flag, preWordLastflag = tr.Find("abcg")
	if flag || preWordLastflag != 2 {
		t.Error("Error")
	}
}

func TestParticiple(t *testing.T) {

	tr := NewTrie()
	tr.Insert("我A")
	tr.Insert("B是")
	tr.Insert("我AB是D")
	tr.Insert("我")

	ss := tr.Participle("我A我AB是DB是D我AB") // 我A/我AB是D/B是/D/我A/B, 其中D和B不是词，会被丢掉
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()

	ss = tr.Participle("我A我AB是DB是擦擦擦我AB") // 我A/我AB是D/B是/擦/擦/擦/我A/B,其中擦和B不是词，会被丢掉
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()

	ss = tr.Participle("D我AB是D")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()

	ss = tr.Participle("D我AB是")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()

	ss = tr.Participle("我A")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()

	ss = tr.Participle("我")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()
}
