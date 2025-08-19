package model

import (
	"github.com/google/uuid"
)

type Product struct {
	BaseModel
	CategoryID  uuid.UUID `json:"categoryId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Stock       int       `json:"stock"`
}

type ProductCategory struct {
	BaseModel
	Name     string     `json:"name"`
	ParentId *uuid.UUID `db:"parent_id" json:"parentId,omitempty"`
}

type ProductImage struct {
	BaseModel
	ProductID uuid.UUID `json:"productId"`
	URL       string    `json:"url"`
	AltText   string    `json:"altText"`
	IsPrimary bool      `json:"isPrimary"`
}

type DiscountType string

const (
	DiscountPercentage DiscountType = "percentage"
	DiscountFixed      DiscountType = "fixed"
)

type Discount struct {
	BaseModel
	ProductID uuid.UUID    `json:"productId"`
	Type      DiscountType `json:"type"`
	Value     float64      `json:"value"`
}
