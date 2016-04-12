package ela

import (
	"github.com/gogather/com/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRouter(t *testing.T) {
	log.Debug = true
	Convey("Router Session sections", t, func() {
		Router("/:hello/:world/123", 123)
		Router("/:hello/123", 1234)
		Router("/", "index")
		ctrl1, param1 := getController("/param1/param2/123")
		ctrl2, param2 := getController("/param1/123")
		ctrl3, param3 := getController("/")
		log.Greenln(ctrl1)
		log.Greenln(param1)
		log.Greenln(ctrl2)
		log.Greenln(param2)
		log.Greenln(ctrl3)
		log.Greenln(param3)

		getArgs("/param1/param2/123", "/:hello/:world/123")

		So(ctrl1, ShouldEqual, 123)
		So(ctrl2, ShouldEqual, 1234)
		So(ctrl3, ShouldEqual, "index")
	})

}
