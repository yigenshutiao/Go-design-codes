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
//go:linkname helloFromP2 p2.helloFromP2
func helloFromP2()

func (p PP1) HelloFromP2() {
	helloFromP2()
}
