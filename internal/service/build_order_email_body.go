package service

import (
	"context"
	"fmt"

	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/Oj-washingtone/savannah-store/internal/repocitory"
)

func BuildOrderEmailBody(ctx context.Context, order *model.Orders, items []*model.OrderItems) string {
	body := "Your order has been created successfully!\n\n"
	body += "Order ID: " + order.ID.String() + "\n"
	body += "Total: Ksh." + fmt.Sprintf("%d", order.Total) + "\n\n"

	paymentStatus := "Not Paid"
	if order.Paid {
		paymentStatus = "Paid"
	}

	body += "Payment Status: " + paymentStatus + "\n"
	body += "Items:\n"

	productRepo := repocitory.NewProductRepository()

	for _, item := range items {
		product, err := productRepo.GetById(ctx, item.ProductID)
		if err != nil {

			body += fmt.Sprintf("- Product ID: %s, Quantity: %d, Price: Ksh.%d\n", item.ProductID, item.Quantity, item.Price)
			continue
		}

		body += fmt.Sprintf("- %s: %s\n  Quantity: %d\n  Unit Price: Ksh.%d\n",
			product.Name,
			product.Description,
			item.Quantity,
			product.Price,
		)
	}

	return body
}
