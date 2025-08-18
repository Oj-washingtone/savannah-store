package model

import "github.com/google/uuid"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type Orders struct {
	BaseModel
	UserID uuid.UUID   `json:"userId"`
	Status OrderStatus `json:"status"`
	Total  int64       `json:"total"` // capture in cents why ?
	Paid   bool        `json:"paid"`
}

type OrderItems struct {
	BaseModel
	OrderID   uuid.UUID `json:"orderId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
	Price     int64     `json:"price"`
}
