package order

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateOrderWithItems(order *Order, items []OrderItem) error
	GetOrderByID(orderID uuid.UUID) (*Order, error)
	GetOrderByIDAndTableIDAndSessionID(orderID, tableID, sessionID uuid.UUID) (*Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderWithItems(order *Order, items []OrderItem) error {
	// ເລີ່ມ Transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		return nil // ບັນທຶກສຳເລັດ (Commit)
	})
}

func (r *orderRepository) GetOrderByID(orderID uuid.UUID) (*Order, error) {
	var order Order
	err := r.db.Preload("OrderItems").First(&order, "id = ?", orderID).Error
	return &order, err
}

func (r *orderRepository) GetOrderByIDAndTableIDAndSessionID(orderID, tableID, sessionID uuid.UUID) (*Order, error) {
	var order Order
	// 🔥 ຄົ້ນຫາດ້ວຍ id ແລະ table_id ພ້ອມກັນ
	err := r.db.Preload("OrderItems").
		Where("id = ? AND table_id = ? AND session_id = ?", orderID, tableID, sessionID).
		First(&order).Error
	return &order, err
}
