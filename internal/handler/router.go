package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *SubscriptionHandler) {

	api := r.Group("/subscriptions")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/total", h.Total)
		api.GET("/:id", h.GetByID)
		api.DELETE("/:id", h.Delete)
	}
}
