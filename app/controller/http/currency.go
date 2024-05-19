package httpcontroller

import (
	"darkness8129/currency-converter/app/service"
	"darkness8129/currency-converter/packages/errs"
	"darkness8129/currency-converter/packages/logging"
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

type getRateRequestQuery struct {
}

func (ctrl *currencyController) getRate(c *gin.Context) {
	logger := ctrl.logger.Named("getRate")

	var query getRateRequestQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		logger.Info("invalid request query", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	logger.Debug("parsed request query", "query", query)

	err = ctrl.services.Currency.GetRate(c, service.GetRateOpt{})
	if err != nil {
		if errs.IsCustom(err) {
			logger.Info(err.Error())
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, nil)
			return
		}

		logger.Error("failed to get rate", "err", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	logger.Info("successfully got rate")
	c.JSON(http.StatusOK, nil)
}
