package ela

import (
	"fmt"
	"github.com/gogather/com/log"
	"testing"
)

func TestSession(t *testing.T) {
	sess := NewSession("tmp")
	path := sess.getPath("ABCDEFGHIJKLMNOPQRST")
	log.Pinkln(path)
	sess.Set("ABCDEFGHIJKLMNOPQRST", "key", "ABCDEFGHIJKLMNOPQRST")
	sess.Set("ABCDEFGHIJKLMNOPQRST", "number", 123456)
	obj, _ := sess.Get("ABCDEFGHIJKLMNOPQRST", "key")
	number, _ := sess.Get("ABCDEFGHIJKLMNOPQRST", "number")

	log.Pinkln("[session] " + obj.(string))
	log.Pinkln("[session] " + fmt.Sprintf("%d", number.(int)))
}
