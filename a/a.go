package a

import (
	"github.com/ggq89/mutualdep/b"
	"github.com/ggq89/mutualdep/c"
	"fmt"
)

type A struct {
	Pb *b.B
	Pc *c.C
}

func New(ic int) *A {
	a := &A{
		Pc:c.New(ic),
	}

	a.Pb = b.New(a)

	return a
}

func Printf(v int)  {
	fmt.Printf("%v", v)
}