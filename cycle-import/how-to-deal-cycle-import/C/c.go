package C

import (
	"go-design-codes/cycle-import/how-to-deal-cycle-import/A"
	"go-design-codes/cycle-import/how-to-deal-cycle-import/B"
)

type CombineAB struct {
	A *A.PackageA
	B *B.PackageB
}

func (c *CombineAB) PrintAll() {
	c.A.PrintA()
	c.B.PrintB()
}
