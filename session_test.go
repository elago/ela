package ela

import (
	"github.com/gogather/com"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSession(t *testing.T) {

	Convey("Test Session sections", t, func() {
		com.Mkdir("tmp")
		sess := NewSession("tmp")
		path := sess.getPath("ABCDEFGHIJKLMNOPQRST")

		So(path, ShouldEqual, "tmp/A/BC/DEFGHIJKLMNOPQRST")

		sess.Set("ABCDEFGHIJKLMNOPQRST", "key", "ABCDEFGHIJKLMNOPQRST")
		sess.Set("ABCDEFGHIJKLMNOPQRST", "number", 123456)
		obj, _ := sess.Get("ABCDEFGHIJKLMNOPQRST", "key")
		number, _ := sess.Get("ABCDEFGHIJKLMNOPQRST", "number")

		So(obj.(string), ShouldEqual, "ABCDEFGHIJKLMNOPQRST")
		So(number.(int), ShouldEqual, 123456)

	})

}
