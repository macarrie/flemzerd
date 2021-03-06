package stats

import (
	"expvar"
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"

	. "github.com/macarrie/flemzerd/objects"
)

var Stats StatsFields

func init() {
	Stats = StatsFields{}
}

func Get() StatsFields {
	//Update runtime fields before updating
	Stats.Runtime.GoRoutines = runtime.NumGoroutine()
	Stats.Runtime.GoMaxProcs = runtime.GOMAXPROCS(0)
	Stats.Runtime.NumCPU = runtime.NumCPU()

	return Stats
}

func Handler() gin.HandlerFunc {
	expvar.Publish("stats", expvar.Func(func() interface{} {
		return Get()
	}))

	return func(c *gin.Context) {
		w := c.Writer
		c.Header("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte("{\n"))
		first := true
		expvar.Do(func(kv expvar.KeyValue) {
			if !first {
				w.Write([]byte(",\n"))
			}
			first = false
			fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
		})
		w.Write([]byte("\n}\n"))
		c.AbortWithStatus(200)
	}
}
