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

// Create godoc
// @Summary Создать подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body dto.CreateSubscriptionRequest true "subscription"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Router /subscriptions [post]
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

// GetByID godoc
// @Summary Получить подписку по id
// @Tags subscriptions
// @Produce json
// @Param id path string true "subscription id"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, toResponse(sub))
}

// List godoc
// @Summary Получить список подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} dto.SubscriptionResponse
// @Router /subscriptions [get]
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

// Delete godoc
// @Summary Удалить подписку
// @Tags subscriptions
// @Param id path string true "subscription id"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Total godoc
// @Summary Сумма подписок за период
// @Tags subscriptions
// @Produce json
// @Param from query string true "MM-YYYY"
// @Param to query string true "MM-YYYY"
// @Param user_id query string false "UUID"
// @Param service_name query string false "Service name"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Router /subscriptions/total [get]
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
