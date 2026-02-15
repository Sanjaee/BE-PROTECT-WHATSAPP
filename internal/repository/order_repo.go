package repository

import (
	"yourapp/internal/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	FindPendingByUserID(userID string) ([]model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindPendingByUserID(userID string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Where("user_id = ? AND status = ?", userID, model.OrderStatusPending).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}
