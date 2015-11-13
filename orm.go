package elaeagnus

import (
	"fmt"
	"reflect"
)

type Pojo struct {
	t string
}

func (this *Pojo) TT() {

}

func init() {
	u := &Pojo{}
	value := reflect.ValueOf(u)
	typ := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		fmt.Printf("method[%d]%s\n", i, typ.Method(i).Name)
	}
}
