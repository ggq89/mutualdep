package c

type C struct {
	Vc int
}

func New(i int) *C {
	return &C{
		Vc:i,
	}
}