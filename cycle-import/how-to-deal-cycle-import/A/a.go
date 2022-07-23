package A

import (
	"fmt"
	"go-design-codes/cycle-import/how-to-deal-cycle-import/model"
)

type PackageA struct {
	B model.PackageBInterface
}

func (PackageA) PrintA() {
	fmt.Println("Hello, I am A")
}

func (a PackageA) PrintAll() {
	a.PrintA()
}
