package app

import (
	"net/http"

	"yourapp/internal/model"
	"yourapp/internal/repository"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo repository.OrderRepository
}

func NewOrderHandler(orderRepo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo}
}

// GetPending returns pending orders for the authenticated user
// GET /api/v1/orders/pending
func (h *OrderHandler) GetPending(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	orders, err := h.orderRepo.FindPendingByUserID(userID.(string))
	if err != nil {
		util.InternalServerError(c, "Failed to fetch orders")
		return
	}

	if orders == nil {
		orders = []model.Order{}
	}

	util.SuccessResponse(c, http.StatusOK, "Pending orders retrieved", gin.H{"orders": orders})
}
