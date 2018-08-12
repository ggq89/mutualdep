package b

import (
	"github.com/ggq89/mutualdep/c"
)

type B struct {
	Pa a
}

type a interface {
	GetC() *c.C
}

func New(a a) *B {
	return &B{
		Pa: a,
	}
}

func (b *B) DisplayC() {
	b.Pa.GetC().Show()
}
