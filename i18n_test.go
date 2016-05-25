package ela

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestI18n(t *testing.T) {
	Convey("I18n sections", t, func() {
		i18n := NewI18n("etc/locale")
		i18n.SetLang("zh_CN")

		hello := i18n.Tr("hello")
		world := i18n.Tr("world")

		fmt.Println(hello)
		fmt.Println(world)

		So(hello, ShouldEqual, "你好")
		So(world, ShouldEqual, "世界")

		i18ne := NewEmptyI18n()

		So(i18ne, ShouldNotBeNil)

		i18ne.Load("etc/locale")

		hello1 := i18n.Tr("_", "hello")
		world1 := i18n.Tr("_", "world")

		fmt.Println(hello1)
		fmt.Println(world1)

		So(hello1, ShouldEqual, "你好")
		So(world1, ShouldEqual, "世界")

		// todo

		Router("/", func(ctx *Context) {
			fmt.Println("ctrl test")
		})

		fmt.Println(middlewares)

		Use(InitI18nModule("etc/locale"))

		fmt.Println(middlewares)

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Body = ioutil.NopCloser(bytes.NewBufferString("This is my request body"))
		// So(err, ShouldBeNil)
		m := elaHandler{}
		m.ServeHTTP(resp, req)

	})

}
