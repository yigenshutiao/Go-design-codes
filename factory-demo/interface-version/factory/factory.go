package factory

import (
	"errors"
	"go-design-codes/factory-demo/interface-version/factA"
	"go-design-codes/factory-demo/interface-version/factB"
)

type Factory interface {
	DoAction(info string) string
}

var FactMap = map[string]Factory{}

func init() {
	FactMap["FactA"] = factA.AFact
	FactMap["FactB"] = factB.BFact
}

func GetFactoryIns(factType string) (Factory, error) {
	if _, ok := FactMap[factType]; !ok {
		return nil, errors.New("invalid fact type")
	}
	return FactMap[factType], nil
}
