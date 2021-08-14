package service

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

const ProfilingUrn = "/profiling"

func buildProfilingApi() API {
	router := gin.New()
	router.GET(ProfilingUrn, wrapPprofHandler)
	return router
}

func wrapPprofHandler(c *gin.Context) {
	http.HandlerFunc(pprof.Index).ServeHTTP(c.Writer, c.Request)
}
