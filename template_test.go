package ela

import (
	"bytes"
	"fmt"
	"github.com/gogather/com/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTemplate(t *testing.T) {
	log.Debug = true
	header := `{{define "header.tpl"}}this is a header template{{end}}`
	index := `{{define "index.tpl"}}this is index file. {{template "header.tpl"}}{{end}}`
	Convey("Template sections", t, func() {
		templates["header.tpl"] = header
		templates["index.tpl"] = index

		getTemplateNames()
		log.Greenln(templatesName)

		So(templatesName, ShouldContain, "header.tpl")
		So(templatesName, ShouldContain, "index.tpl")

		b := bytes.NewBuffer(make([]byte, 0))
		t, _ := parseFiles(templatesName...)
		t.ExecuteTemplate(b, "index.tpl", nil)
		result := fmt.Sprintf("%s", b)

		So(result, ShouldEqual, "this is index file. this is a header template")

	})

}
