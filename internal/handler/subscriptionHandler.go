package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"subscriptionManagement/internal/handler/dto"
	"subscriptionManagement/internal/model"
	"subscriptionManagement/internal/service"
	"subscriptionManagement/internal/utils"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(s service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: s}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req dto.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start, err := utils.ParseMonthYear(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}

	var end *time.Time
	if req.EndDate != "" {
		t, err := utils.ParseMonthYear(req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
			return
		}
		end = &t
	}

	sub := model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start,
		EndDate:     end,
	}

	created, err := h.service.Create(c.Request.Context(), sub)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toResponse(created))
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, toResponse(sub))
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	subs, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]dto.SubscriptionResponse, 0, len(subs))
	for _, s := range subs {
		tmp := s
		resp = append(resp, toResponse(&tmp))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) Total(c *gin.Context) {

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := utils.ParseMonthYear(fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
		return
	}

	to, err := utils.ParseMonthYear(toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to"})
		return
	}

	filter := model.SubscriptionFilter{
		From: from,
		To:   to,
	}

	if user := c.Query("user_id"); user != "" {
		filter.UserID = &user
	}

	if name := c.Query("service_name"); name != "" {
		filter.ServiceName = &name
	}

	total, err := h.service.CalculateTotal(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}

func toResponse(sub *model.Subscription) dto.SubscriptionResponse {
	resp := dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   utils.FormatMonthYear(sub.StartDate),
	}

	if sub.EndDate != nil {
		resp.EndDate = utils.FormatMonthYear(*sub.EndDate)
	}

	return resp
}
