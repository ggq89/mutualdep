package main

import "github.com/ggq89/mutualdep/a"

func main() {
	a := a.New(3)
	a.Pb.DisplayC()
}
