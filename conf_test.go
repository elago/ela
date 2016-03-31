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
			conf.parseItems()

			log.Greenln("warning stack")
			log.Pinkln("==============")
			log.Blueln(conf.GetWarnings())

			log.Greenln("mistake value")
			log.Pinkln("==============")
			log.Blueln(conf.GetString("mysql", "host"))

			log.Greenln("about bool")
			log.Pinkln("==============")
			conf.SetBool("test", "dev", true)
			dev, _ := conf.GetBool("_", "dev")
			log.Blueln("%v", dev)

			content = conf.serialize()
			log.Greenln("serialize value")
			log.Pinkln("==============")
			log.Blueln(content)

			log.Greenln("hex")
			log.Pinkln("==============")
			hex, _ := conf.GetInt("_", "hex")
			log.Blueln(hex)

			log.Greenln("get empty")
			log.Pinkln("==============")
			_, err = conf.GetInt("_", "heloo")
			if err != nil {
				log.Blueln(err)
			}

			log.Greenln("save config file")
			log.Pinkln("==============")
			err := conf.Save("tmp/test.ini")
			if err == nil {
				log.Blueln("done!")
			}

		}

		val1, _ := conf.Get("_", "port")
		So(val1, ShouldEqual, 80)

		val2, _ := conf.GetString("_", "appname")
		So(val2, ShouldEqual, "my application")

		val3, _ := conf.GetString("mysql", "password")
		So(val3, ShouldEqual, "liju#n")

		val4, _ := conf.GetString("mysql", "host")
		So(val4, ShouldEqual, `"192.168.1.11" = GHJ`)

		val5, _ := conf.GetBool("_", "dev")
		So(val5, ShouldEqual, true)

		val6, _ := conf.GetFloat("_", "pi")
		So(val6, ShouldEqual, 3.14)

		val7, _ := conf.GetInt("_", "hex")
		So(val7, ShouldEqual, 0x24)

		val8 := conf.GetIntDefault("section", "key", 100)
		So(val8, ShouldEqual, 100)
	})
}
