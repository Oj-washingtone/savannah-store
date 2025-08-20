package model

import (
	"github.com/google/uuid"
)

type Cart struct {
	BaseModel
	UserId uuid.UUID `json:"userId"`
}

type CartItem struct {
	BaseModel
	CartId    uuid.UUID `db:"cart_id" json:"cartId"`
	ProductId uuid.UUID `db:"product_id" json:"productId"`
	Quantity  int       `db:"quantity" json:"quantity"`
	Price     int64     `db:"price" json:"price"`
}
