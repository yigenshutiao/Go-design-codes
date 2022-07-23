package main

import "go-design-codes/cycle-import/what-is-cycle-import/A"

/*
output:

package go-design-codes/cycle-import/what-is-cycle-import
	imports go-design-codes/cycle-import/what-is-cycle-import/A
	imports go-design-codes/cycle-import/what-is-cycle-import/B
	imports go-design-codes/cycle-import/what-is-cycle-import/A: import cycle not allowed
*/
func main() {
	A.GetB(13)
}
