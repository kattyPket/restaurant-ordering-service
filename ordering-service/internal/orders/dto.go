// dto.go
package order

import "github.com/google/uuid"

type PlaceOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderItemRequest struct {
	MenuItemID uuid.UUID `json:"menu_item_id" validate:"required"`
	Quantity   int       `json:"quantity" validate:"required,min=1"`
	Note       string    `json:"note"` // Optional
}

type OrderResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	OrderNumber string    `json:"order_number"`
	Status      string    `json:"status"`
	GrandTotal  float64   `json:"grand_total"`
}
