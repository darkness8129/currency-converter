package httpcontroller

import (
	"darkness8129/currency-converter/app/service"
	"darkness8129/currency-converter/packages/logging"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Options struct {
	Router   *gin.Engine
	Logger   *logging.Logger
	Services service.Services
}

type controllerOptions struct {
	RouterGroup *gin.RouterGroup
	Logger      *logging.Logger
	Services    service.Services
}

func New(opt Options) {
	opt.Router.Use(gin.Logger(), gin.Recovery(), corsMiddleware)

	controllerOpt := controllerOptions{
		RouterGroup: opt.Router.Group("/api/v1"),
		Logger:      opt.Logger.Named("httpController"),
		Services:    opt.Services,
	}

	newCurrencyController(controllerOpt)
}

func corsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}

// panicHandler provides unified panic handling for all http controllers
func panicHandler(logger *logging.Logger, handler func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.Named("panicHandler")

		defer func() {
			err := recover()
			if err != nil {
				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", err))
				if err != nil {
					logger.Error("failed to abort with error", "err", err)
				}
			}
		}()

		handler(c)
	}
}
