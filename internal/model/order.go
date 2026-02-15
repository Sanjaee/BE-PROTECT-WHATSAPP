package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Status    string         `gorm:"type:varchar(50);default:'pending';index" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

func (Order) TableName() string {
	return "orders"
}

const (
	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)
