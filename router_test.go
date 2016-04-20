package ela

import (
	"github.com/gogather/com/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRouter(t *testing.T) {
	log.Debug = true
	Convey("Router sections", t, func() {
		Router("/:hello/:world/123", 123, 456)
		Router("/:hello/123", 1234)
		Router("/", "index")
		node1:= getController("/param1/param2/123").(uriMode)
		ctrl1:=node1.fun
		param1:=node1.argsMap
		node2:=getController("/param1/123").(uriMode)
		ctrl2:=node2.fun
		param2 :=node2.argsMap
		node3:=getController("/").(uriMode)
		ctrl3:=node3.fun
		param3 := node3.argsMap
		log.Greenln(ctrl1)
		log.Greenln(param1)
		log.Greenln(ctrl2)
		log.Greenln(param2)
		log.Greenln(ctrl3)
		log.Greenln(param3)

		getArgs("/param1/param2/123", "/:hello/:world/123")

		So(ctrl1, ShouldContain, 123)
		So(ctrl1, ShouldContain, 456)
		So(ctrl2, ShouldContain, 1234)
		So(ctrl3, ShouldContain, "index")
		So(param1["hello"], ShouldEqual, "param1")
		So(param1["world"], ShouldEqual, "param2")
	})

}
