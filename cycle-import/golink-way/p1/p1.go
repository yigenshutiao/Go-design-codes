package p1

import (
	"fmt"
	_ "unsafe"
)

type PP1 struct{}

func (p *PP1) HelloFromP1() {
	fmt.Println("Hello, I am p1 func")
}

// p1 包中使用p2的私有方法，原文：https://golang.org/cmd/compile/#hdr-Compiler_Directives
/*
	This special directive does not apply to the Go code that follows it. Instead,
	the //go:linkname directive instructs the compiler to use “importpath.name” as
	the object file symbol name for the variable or function declared as “localname”
	in the source code. If the “importpath.name” argument is omitted, the directive uses
	the symbol's default object file symbol name and only has the effect of making the symbol
	accessible to other packages. Because this directive can subvert the type system and package modularity,
	it is only enabled in files that have imported "unsafe".
*/
//go:linkname helloFromP2 p2.helloFromP2
func helloFromP2()

func (p PP1) HelloFromP2() {
	helloFromP2()
}
