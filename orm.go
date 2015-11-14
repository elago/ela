package elaeagnus

import (
	"fmt"
	"reflect"
)

type Pojo struct {
}

func (this *Pojo) Get() {

}

func RegisterModel(model interface{}) {
	value := reflect.ValueOf(model)
	typ := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		fmt.Printf("method[%d]%s\n", i, typ.Method(i).Name)
	}
}
