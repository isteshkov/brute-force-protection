package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	MetricsUrn = "/metrics"
	HealthUrn  = "/health"
)

type API interface {
	Run(address ...string) error
}

func buildMetricsApi(s *Service) API {
	router := gin.New()
	router.GET(MetricsUrn, wrapPromHandler)
	router.GET(HealthUrn, s.healthHandler)
	return router
}

func wrapPromHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func (s *Service) healthHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
	return
}
