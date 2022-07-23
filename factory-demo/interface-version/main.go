package main

import (
	"fmt"
	"go-design-codes/factory-demo/interface-version/factory"
)

func main() {
	a, err := factory.GetFactoryIns("FactA")
	if err != nil {
		panic("invalid fact type")
	}
	fmt.Println(a.DoAction("have a nice day"))
}
