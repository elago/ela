package ela

import (
	"github.com/gogather/com/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfig(t *testing.T) {
	Convey("Test Config sections", t, func() {
		conf := NewConfig("etc/test.ini")
		content, err := conf.readConfigFile()
		if err != nil {
			log.Redln(err)
		} else {
			log.Greenln("raw content")
			log.Pinkln("==============")

			log.Blueln(content)
			content = conf.filterComment()

			log.Greenln("filter comment")
			log.Pinkln("==============")

			log.Blueln(content)
			arraylize := conf.arraylize()

			log.Greenln("arraylize")
			log.Pinkln("==============")
			count := len(arraylize)
			log.Bluef("[%d]\n", count)
			for i := 0; i < count; i++ {
				log.Bluef("[%d]\t%s\n", i, arraylize[i])
			}

			log.Greenln("parse items")
			log.Pinkln("==============")
		}
	})
}
