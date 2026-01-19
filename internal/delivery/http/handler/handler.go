package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/scmbr/subscription-aggregator/internal/delivery/http/handler/v1"
	"github.com/scmbr/subscription-aggregator/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery(), gin.Logger())
	h.initAPI(router)
	return router
}
func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.service)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
