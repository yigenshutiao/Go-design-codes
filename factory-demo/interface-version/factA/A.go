package factA

import "fmt"

var AFact = factA{}

type factA struct {
}

func (a factA) DoAction(info string) string {
	return fmt.Sprintf("I am Fact A ins, my info:%+v", info)
}
