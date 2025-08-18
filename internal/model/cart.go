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
	CartId    uuid.UUID `json:"cartId"`
	ProductId uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
}
