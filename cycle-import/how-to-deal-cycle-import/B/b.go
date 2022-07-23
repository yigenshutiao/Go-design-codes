package B

import (
	"fmt"
	"go-design-codes/cycle-import/how-to-deal-cycle-import/model"
)

type PackageB struct {
	A model.PackageAInterface
}

func (PackageB) PrintB() {
	fmt.Println("Hello, I am B")
}

func (b PackageB) PrintAll() {
	b.PrintB()
}
