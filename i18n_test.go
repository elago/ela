package ela

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestI18n(t *testing.T) {
	Convey("I18n sections", t, func() {
		i18n := NewI18n("etc/locale")

		hello := i18n.Tr("zh_CN", "_", "hello")
		world := i18n.Tr("zh_CN", "_", "world")

		fmt.Println(hello)
		fmt.Println(world)

		So(hello, ShouldEqual, "你好")
		So(world, ShouldEqual, "世界")
	})

}
