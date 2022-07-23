package factory

import (
	"fmt"
	"reflect"
)

var factoryMap map[string]interface{}

func init() {
	factoryMap = make(map[string]interface{})
}

func RegisterFactory(key string, factory interface{}) {
	factoryMap[key] = factory
}

func CallFactory(key string, args ...interface{}) []interface{} {
	if factor, ok := factoryMap[key]; ok {
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			in[i] = reflect.ValueOf(arg)
		}

		outList := reflect.ValueOf(factor).Call(in)
		result := make([]interface{}, len(outList))
		for i, out := range outList {
			result[i] = out.Interface()
		}
		return result
	} else {
		panic(fmt.Errorf("factory not found:  %s", key))
	}
	return nil
}
