package order

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ordering-service/internal/menus"
)

type Order struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	OrderNumber    string     `gorm:"type:varchar;unique"`
	TableID        *uuid.UUID `gorm:"type:uuid"`
	SessionID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	OrderType      string     `gorm:"type:varchar"`
	StaffID        *uuid.UUID `gorm:"type:uuid"`
	CustomerCount  *int       `gorm:"type:int"`
	Subtotal       float64    `gorm:"type:decimal(10,2)"`
	DiscountAmount float64    `gorm:"type:decimal(10,2);default:0"`
	TaxAmount      float64    `gorm:"type:decimal(10,2);default:0"`
	GrandTotal     float64    `gorm:"type:decimal(10,2)"`
	Status         string     `gorm:"type:varchar"`
	SlipImage      *string    `gorm:"type:varchar"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	return
}

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID    uuid.UUID `gorm:"type:uuid"`
	MenuItemID uuid.UUID `gorm:"type:uuid"`
	Quantity   int       `gorm:"type:int"`
	UnitPrice  float64   `gorm:"type:decimal(10,2)"`
	TotalPrice float64   `gorm:"type:decimal(10,2)"`
	Note       *string   `gorm:"type:varchar"`
	Status     string    `gorm:"type:varchar"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	MenuItem menus.MenuItem `gorm:"foreignKey:MenuItemID"`
}

func (item *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	return
}
