package A

import (
	"fmt"
	"go-design-codes/cycle-import/what-is-cycle-import/B"
)

func Hello(name string) string {
	return fmt.Sprintf("My name is %+v", name)
}

func GetB(b int) int {
	return B.Add(b)
}
