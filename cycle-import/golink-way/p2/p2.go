package p2

import "fmt"

type pp1 interface {
	HelloFromP1()
}

type PP2 struct {
	PP1 pp1
}

func New(pp1 pp1) *PP2 {
	return &PP2{
		PP1: pp1,
	}
}

// HelloFromP1Side p2包通过interface使用p1包中的方法
func (p PP2) HelloFromP1Side() {
	p.PP1.HelloFromP1()
}

func helloFromP2() {
	fmt.Println("Hello, I am p2 func")
}
