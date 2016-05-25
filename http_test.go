package ela

import (
	"github.com/gogather/com/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type MidStruct1 struct {
	name string
	age  int
}

func action(m MidStruct1, t *testing.T) {
	log.Println(m.name)
	if m.name != "lee" {
		t.Error("maybe injection failed")
	}
}

func TestInjection(t *testing.T) {
	log.Debug = true
	Convey("Injection sections", t, func() {
		e := Web()
		e.Use(MidStruct1{name: "lee", age: 24})
		e.Use(&testing.T{})
		injectFuc(action)
	})

}
