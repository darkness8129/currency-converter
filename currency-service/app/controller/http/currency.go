package httpcontroller

import (
	"darkness8129/currency-converter/currency-service/app/service"
	"darkness8129/currency-converter/currency-service/packages/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

type currencyController struct {
	logger   *logging.Logger
	services service.Services
}

func newCurrencyController(opt controllerOptions) {
	logger := opt.Logger.Named("currencyController")

	c := currencyController{
		logger:   logger,
		services: opt.Services,
	}

	group := opt.RouterGroup.Group("/currency")
	group.GET("/rate", panicHandler(logger, c.getRate))
}

type getRateResponseBody struct {
	Rate float64 `json:"rate"`
}

type getRateError struct {
	Error string `json:"error"`
}

func (ctrl *currencyController) getRate(c *gin.Context) {
	logger := ctrl.logger.Named("getRate")

	rate, err := ctrl.services.Currency.GetRate(c)
	if err != nil {
		logger.Error("failed to get rate", "err", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, getRateError{"internal server error"})
		return
	}

	logger.Info("successfully got rate", "rate", rate)
	c.JSON(http.StatusOK, getRateResponseBody{rate})
}
