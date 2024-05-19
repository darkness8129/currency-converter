package httpcontroller

import (
	"darkness8129/currency-converter/app/service"
	"darkness8129/currency-converter/packages/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

type subscriptionController struct {
	logger   *logging.Logger
	services service.Services
}

func newSubscriptionController(opt controllerOptions) {
	logger := opt.Logger.Named("subscriptionController")

	c := subscriptionController{
		logger:   logger,
		services: opt.Services,
	}

	group := opt.RouterGroup.Group("/subscriptions")
	group.POST("/subscribe", panicHandler(logger, c.subscribe))
}

type subscribeForm struct {
	Email string `form:"email" binding:"required"`
}

type subscribeError struct {
	Error string `json:"error"`
}

func (ctrl *subscriptionController) subscribe(c *gin.Context) {
	logger := ctrl.logger.Named("subscribe")

	var formData subscribeForm
	err := c.ShouldBind(&formData)
	if err != nil {
		logger.Info("invalid form", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, subscribeError{err.Error()})
		return
	}
	logger.Debug("parsed form", "formData", formData)

	err = ctrl.services.Subscription.Subscribe(c, formData.Email)
	if err != nil {
		if err == service.ErrSubscriberAlreadyExists {
			logger.Info(err.Error())
			c.AbortWithStatusJSON(http.StatusConflict, subscribeError{err.Error()})
			return
		}

		logger.Error("failed to subscribe", "err", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, subscribeError{"internal server error"})
		return
	}

	logger.Info("successfully subscribed", "email", formData.Email)
	c.JSON(http.StatusOK, nil)
}
