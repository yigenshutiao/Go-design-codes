package main

import (
	"go-design-codes/cycle-import/golink-way/p1"
	"go-design-codes/cycle-import/golink-way/p2"
)

func main() {
	pp1 := p1.PP1{}
	pp2 := p2.New(&pp1)
	pp2.HelloFromP1Side()
}
