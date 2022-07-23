package main

import (
	"go-design-codes/cycle-import/how-to-deal-cycle-import/A"
	"go-design-codes/cycle-import/how-to-deal-cycle-import/B"
)

func main() {
	a := new(A.PackageA)
	b := new(B.PackageB)
	a.PrintAll()
	b.PrintAll()
}
