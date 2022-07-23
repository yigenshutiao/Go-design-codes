package factA

import (
	"context"
	"fmt"

	"go-design-codes/factory-demo/reflect-version/factory"
	"go-design-codes/factory-demo/reflect-version/model"
)

func init() {
	factory.RegisterFactory("A", new(PackageA).GetA)
}

type PackageA struct{}

func (a PackageA) GetA(ctx context.Context, param model.AOpts) (string, error) {
	res := "my name is" + param.Name + ", My address is" + param.Address
	return res, nil
}

func (a PackageA) PrintAll() {
	p := model.BOpts{
		ID:      "10086",
		Company: "fk",
	}

	fmt.Println(factory.CallFactory("B", context.TODO(), p)[0].(string))
}
