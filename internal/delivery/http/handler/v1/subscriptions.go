package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler_dto "github.com/scmbr/subscription-aggregator/internal/delivery/http/handler/dto"
	"github.com/scmbr/subscription-aggregator/internal/service"
	service_dto "github.com/scmbr/subscription-aggregator/internal/service/dto"
	"github.com/scmbr/subscription-aggregator/pkg/logger"
)

func (h *Handler) initSubscriptionsRoutes(api *gin.RouterGroup) {
	subscriptions := api.Group("/subscriptions")
	{
		subscriptions.POST("", h.createSubscription)
		subscriptions.GET("", h.getAllSubscriptions)
		subscriptions.GET("/:id", h.getSubscriptionById)
		subscriptions.PUT("/:id", h.updateSubscriptionById)
		subscriptions.DELETE("/:id", h.deleteSubscriptionById)
		subscriptions.GET("/total", h.getSubscriptionTotalPrice)
	}
}
func (h *Handler) createSubscription(c *gin.Context) {
	var input handler_dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var endDate *time.Time
	if input.EndDate != nil {
		endDate = &input.EndDate.Time
	}

	id, err := h.service.Subscription.CreateSubscription(c.Request.Context(), &service_dto.CreateSubscriptionInput{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   input.StartDate.Time,
		EndDate:     endDate,
	})
	if err != nil {
		logger.Error("error occurred while creating a subscription", err,
			map[string]interface{}{
				"service_name": input.ServiceName,
				"price":        input.Price,
				"user_id":      input.UserID,
				"start_date":   input.StartDate,
				"end_date":     input.EndDate,
			})
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.JSON(http.StatusCreated, handler_dto.CreateSubscriptionResponse{
		Id: id,
	})
}
func (h *Handler) getAllSubscriptions(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		newResponse(c, http.StatusBadRequest, "invalid limit")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		newResponse(c, http.StatusBadRequest, "invalid offset")
		return
	}
	res, err := h.service.Subscription.GetAllSubscriptions(c.Request.Context(), service_dto.GetAllSubscriptionsInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error(
			"error occurred while getting all subscriptions",
			err,
			map[string]interface{}{
				"limit":  limit,
				"offset": offset,
			},
		)
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	subscriptions := make([]handler_dto.GetSubscriptionResponse, 0, len(res.Subscriptions))
	for _, i := range res.Subscriptions {
		var endDate *handler_dto.MonthYear
		if i.EndDate != nil {
			endDate = &handler_dto.MonthYear{Time: *i.EndDate}
		}
		subscriptions = append(subscriptions, handler_dto.GetSubscriptionResponse{
			Id:          i.ID,
			ServiceName: i.ServiceName,
			Price:       i.Price,
			UserID:      i.UserID,
			StartDate:   handler_dto.MonthYear{Time: i.StartDate},
			EndDate:     endDate,
		})
	}
	c.JSON(http.StatusOK, handler_dto.GetAllSubscriptionsResponse{
		Total:         res.Total,
		Subscriptions: subscriptions,
	})
}
func (h *Handler) getSubscriptionById(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}
	res, err := h.service.Subscription.GetSubscriptionById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			newResponse(c, http.StatusNotFound, "subscription not found")
			return
		}
		logger.Error(
			"error occurred while getting subscription by id",
			err,
			map[string]interface{}{
				"subscription_id": id,
			},
		)
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	var endDate *handler_dto.MonthYear
	if res.EndDate != nil {
		endDate = &handler_dto.MonthYear{Time: *res.EndDate}
	}
	subscription := handler_dto.GetSubscriptionResponse{
		Id:          res.ID,
		ServiceName: res.ServiceName,
		Price:       res.Price,
		UserID:      res.UserID,
		StartDate:   handler_dto.MonthYear{Time: res.StartDate},
		EndDate:     endDate,
	}
	c.JSON(http.StatusOK, subscription)
}
func (h *Handler) updateSubscriptionById(c *gin.Context) {
	var input handler_dto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}
	var startDate *time.Time
	if input.StartDate != nil {
		startDate = &input.StartDate.Time
	}

	var endDate *time.Time
	if input.EndDate != nil {
		endDate = &input.EndDate.Time
	}
	err := h.service.Subscription.UpdateSubscriptionById(c.Request.Context(), id, &service_dto.UpdateSubscriptionInput{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	})
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			newResponse(c, http.StatusNotFound, "subscription not found")
			return
		}
		logger.Error(
			"error occurred while updating subscription by id",
			err,
			map[string]interface{}{
				"subscription_id": id,
			},
		)
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) deleteSubscriptionById(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}
	err := h.service.Subscription.DeleteSubscriptionById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			newResponse(c, http.StatusNotFound, "subscription not found")
			return
		}
		logger.Error(
			"error occurred while deleting subscription by id",
			err,
			map[string]interface{}{
				"subscription_id": id,
			},
		)
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) getSubscriptionTotalPrice(c *gin.Context) {

	var input handler_dto.GetTotalPriceRequest

	if err := c.ShouldBindQuery(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.UserID != nil {
		if _, err := uuid.Parse(*input.UserID); err != nil {
			newResponse(c, http.StatusBadRequest, "invalid user_id format")
			return
		}
	}
	var startDate *time.Time
	if input.StartDate != nil {
		startDate = &input.StartDate.Time
	}
	var endDate *time.Time
	if input.EndDate != nil {
		endDate = &input.EndDate.Time
	}
	logger.Debug("handler data:", map[string]interface{}{
		"user_id":      input.UserID,
		"service_name": input.ServiceName,
		"start_date":   startDate,
		"end_date":     endDate,
	})
	total, err := h.service.Subscription.GetSubscriptionsTotalPrice(c.Request.Context(), &service_dto.GetTotalPriceInput{
		UserID:      input.UserID,
		ServiceName: input.ServiceName,
		StartDate:   startDate,
		EndDate:     endDate,
	})
	if err != nil {
		logger.Error("error occurred while getting total price", err, map[string]interface{}{
			"user_id":      input.UserID,
			"service_name": input.ServiceName,
			"start_date":   input.StartDate,
			"end_date":     input.EndDate,
		})
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, handler_dto.GetTotalPriceResponse{
		TotalPrice: total,
	})
}
