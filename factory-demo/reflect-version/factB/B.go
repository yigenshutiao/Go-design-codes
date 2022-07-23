package factB

import (
	"context"
	"fmt"

	"go-design-codes/factory-demo/reflect-version/factory"
	"go-design-codes/factory-demo/reflect-version/model"
)

func init() {
	factory.RegisterFactory("B", new(PackageB).GetB)
}

type PackageB struct{}

func (b PackageB) GetB(ctx context.Context, param model.BOpts) (string, error) {
	res := "My ID is" + param.ID + ", my company is " + param.Company
	return res, nil
}

func (b PackageB) PrintAll() {
	p := model.AOpts{
		Name:    "A",
		Address: "Bj",
	}
	fmt.Println(factory.CallFactory("A", context.TODO(), p)[0].(string))
}
