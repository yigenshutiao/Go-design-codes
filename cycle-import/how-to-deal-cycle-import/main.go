package main

import (
	"go-design-codes/cycle-import/how-to-deal-cycle-import/A"
	"go-design-codes/cycle-import/how-to-deal-cycle-import/B"
	C2 "go-design-codes/cycle-import/how-to-deal-cycle-import/C"
)

func main() {
	a := new(A.PackageA)
	b := new(B.PackageB)

	c := new(C2.CombineAB)
	c.A = a
	c.B = b

	c.PrintAll()
}
