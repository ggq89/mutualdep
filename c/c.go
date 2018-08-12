package c

import (
	"github.com/ggq89/mutualdep/a/f"
)

type C struct {
	Vc int
}

func New(i int) *C {
	return &C{
		Vc: i,
	}
}

func (c *C) Show() {
	f.Printf(c.Vc)
}
