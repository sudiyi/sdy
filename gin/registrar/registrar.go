package registrar

import (
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var registers map[string]func(r *gin.Engine) = map[string]func(r *gin.Engine){
	"noroute": func(r *gin.Engine) {
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "the interface not found"})
		})
	},

	"monitors": func(r *gin.Engine) {
		r.HEAD("/monitors/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "server running"})
		})
	},

	"debug": func(r *gin.Engine) {
		pprof.Register(r)

		r.GET("/debug/goroutines", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"num": runtime.NumGoroutine()})
		})

		r.GET("/debug/gc_stat", func(c *gin.Context) {
			stat := &debug.GCStats{}
			debug.ReadGCStats(stat)
			pause := stat.Pause
			if len(pause) > 10 {
				pause = pause[:10]
			}
			c.JSON(http.StatusOK, gin.H{
				"last_gc": stat.LastGC.Format("2006-01-02 15:04:05"),
				"num_gc":  stat.NumGC,
				"pause":   pause,
			})
		})
	},
}

func Register(r *gin.Engine, strs ...string) {
	for _, str := range strs {
		if register, ok := registers[str]; ok {
			register(r)
		}
	}
}

func RegisterAll(r *gin.Engine) {
	for _, register := range registers {
		register(r)
	}
}
