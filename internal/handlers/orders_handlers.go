package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/Oj-washingtone/savannah-store/internal/repocitory"
	"github.com/Oj-washingtone/savannah-store/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateOrder godoc
// @Summary Create a new order for the authenticated user
// @Description Creates an order based on the user's cart and persists order items. Calculates total automatically. Requires user authentication.
// @Tags Orders
// @Accept json
// @Produce json
// @Success 201 {object} model.Orders "Order created successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Cart not found or empty"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /orders/create [post]
func CreateOrder(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims := userClaims.(map[string]interface{})

	auth0ID := claims["sub"].(string)

	user, err := repocitory.NewUserRepository().GetByAuth0Id(c, auth0ID)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to load user account", err.Error())
		return
	}

	cart, err := repocitory.NewShoppingCartRepository().GetShoppingCart(c.Request.Context(), user.ID)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			RespondError(c, http.StatusNotFound, "cart not found", "No cart found for the user")
			return
		}
	}

	// cart items

	cartItemsRepo := repocitory.NewCartItemsRepository()

	items, err := cartItemsRepo.GetItems(c.Request.Context(), cart.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			RespondError(c, http.StatusNotFound, "Empty cart", "No items found in the cart")
			return
		} else {
			RespondError(c, http.StatusInternalServerError, "failed to load cart items", err.Error())
			return
		}
	}

	var total int64 = 0
	for _, item := range items {
		total += item.Price * int64(item.Quantity)
	}

	order := &model.Orders{
		UserID: user.ID,
		Total:  total,
	}

	order.ID = uuid.New()

	theOrder, err := repocitory.NewOrdersRepository().Create(c.Request.Context(), order)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to create order", err.Error())
		return
	}

	var orderItems []*model.OrderItems
	for _, item := range items {
		orderItem := &model.OrderItems{
			OrderID:   theOrder.ID,
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}

		orderItem.ID = uuid.New()
		orderItems = append(orderItems, orderItem)
	}

	orderItemsRepo := repocitory.NewOrderItemsRepository()

	err = orderItemsRepo.CreateBulk(c.Request.Context(), orderItems)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to create order items", err.Error())
		return
	}

	err = cartItemsRepo.ClearCart(c.Request.Context(), cart.ID)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Unable to clear cart", err.Error())
		return
	}

	// TODO: send SMS

	// email to admin

	adminEmail := os.Getenv("ADMIN_EMAIL")
	emailBody := service.BuildOrderEmailBody(c.Request.Context(), theOrder, orderItems)

	service.SendEmail(adminEmail, "New Order Created", emailBody)

	RespondSuccess(c, http.StatusCreated, "Order created successfully", theOrder)
}

// GetAllOrders godoc
// @Summary Get all orders
// @Description Fetches all orders with associated order items.
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {array} model.Orders "List of all orders"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /orders [get]
func GetAllOrders(c *gin.Context) {
	orders, err := repocitory.NewOrdersRepository().GetAll(c.Request.Context())

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to fetch orders", err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, "Orders fetched successfully", orders)
}
