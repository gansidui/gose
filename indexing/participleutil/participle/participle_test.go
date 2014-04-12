package participle

import (
	"fmt"
	"testing"
)

func TestFind(t *testing.T) {
	p := NewParticiple()
	p.Insert("ab")
	p.Insert("你好啊")
	p.Insert("")

	if p.Num() != 3 {
		t.Error("Error")
	}

	if p.ForwardFind("cc") || p.BackwardFind("cc") {
		t.Error("Error")
	}

	if !p.ForwardFind("ab") || !p.BackwardFind("ab") {
		t.Error("Error")
	}

	if !p.ForwardFind("你好啊") || !p.BackwardFind("你好啊") {
		t.Error("Error")
	}

	if !p.ForwardFind("") || !p.BackwardFind("") {
		t.Error("Error")
	}
}

func TestForwardMaxMatch(t *testing.T) {
	p := NewParticiple()
	p.Insert("我A")
	p.Insert("B是")
	p.Insert("我AB是D")

	ss := p.ForwardMaxMatch("我A我AB是DB是D我AB")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()

	ss = p.ForwardMaxMatch("D我AB是")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()

	p.Insert("学")
	p.Insert("历")
	p.Insert("史")
	p.Insert("学")
	p.Insert("好")
	p.Insert("学历")
	p.Insert("历史")
	p.Insert("史学")
	p.Insert("学好")

	ss = p.ForwardMaxMatch("学历史学好")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()
}

func TestBackwardMaxMatch(t *testing.T) {
	p := NewParticiple()
	p.Insert("学")
	p.Insert("历")
	p.Insert("史")
	p.Insert("学")
	p.Insert("好")
	p.Insert("学历")
	p.Insert("历史")
	p.Insert("史学")
	p.Insert("学好")

	ss := p.BackwardMaxMatch("学历史学好")
	for _, v := range ss {
		fmt.Printf("%s/", v)
	}
	fmt.Println()
}
