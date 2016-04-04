package ela

import (
	"github.com/gogather/com/log"
	// "net/http"
	"time"
)

var (
	logTimeFormat = "2006-01-02 15:04:05"
	tag           = "ela"
)

func requestLog(ctx Context) {
	log.Printf("[%s] %s: starting %s %s %s\n", tag, time.Now().Format(logTimeFormat), ctx.GetRequest().Method, ctx.GetRequest().RequestURI, ctx.GetRequest().RemoteAddr)
}

func responseLog(ctx Context) {
	content := "[%s] %s: complete %s %s %d\n"
	switch ctx.GetResponseWriter().Status() {
	case 301, 302:
		log.Bluef(content, tag, time.Now().Format(logTimeFormat), ctx.GetRequest().Method, ctx.GetRequest().RequestURI, ctx.GetResponseWriter().Status())
	case 304:
		log.Greenf(content, tag, time.Now().Format(logTimeFormat), ctx.GetRequest().Method, ctx.GetRequest().RequestURI, ctx.GetResponseWriter().Status())
	case 401, 403:
		log.Yellowf(content, tag, time.Now().Format(logTimeFormat), ctx.GetRequest().Method, ctx.GetRequest().RequestURI, ctx.GetResponseWriter().Status())
	case 404:
		log.Redf(content, tag, time.Now().Format(logTimeFormat), ctx.GetRequest().Method, ctx.GetRequest().RequestURI, ctx.GetResponseWriter().Status())
	case 500:
		log.Pinkf(content, tag, time.Now().Format(logTimeFormat), ctx.GetRequest().Method, ctx.GetRequest().RequestURI, ctx.GetResponseWriter().Status())
	default:
		log.Printf(content, tag, time.Now().Format(logTimeFormat), ctx.GetRequest().Method, ctx.GetRequest().RequestURI, ctx.GetResponseWriter().Status())
	}

}
