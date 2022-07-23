package B

import "go-design-codes/cycle-import/what-is-cycle-import/A"

func Add(a int) int {
	return a + 10
}

func GetA(s string) string {
	return A.Hello(s)
}
