package factB

import "fmt"

var BFact = factB{}

type factB struct {
}

func (b factB) DoAction(info string) string {
	return fmt.Sprintf("I am Fact B ins, my info:%+v", info)
}
