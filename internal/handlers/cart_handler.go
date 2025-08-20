package handlers

import (
	"errors"
	"net/http"

	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/Oj-washingtone/savannah-store/internal/repocitory"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ItemBody struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

//TODO: Add security annotation

// Add product to cart
// AddProductToCart godoc
// @Summary Add product to cart
// @Description Add a product to the shopping cart
// @Tags Shopping Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body ItemBody true "Item to add"
// @Success 201 {object} model.CartItem
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/create [post]
func AddToCart(c *gin.Context) {

	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims := userClaims.(map[string]interface{})

	// Auth0 user_id (sub claim)
	auth0ID := claims["sub"].(string)

	user, err := repocitory.NewUserRepository().GetByAuth0Id(c, auth0ID)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to get user", err.Error())
		return
	}

	cartRepo := repocitory.NewShoppingCartRepository()

	cart, err := cartRepo.GetShoppingCart(c.Request.Context(), user.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			cart, err = cartRepo.CreateCart(c.Request.Context(), user.ID)

			if err != nil {
				RespondError(c, http.StatusInternalServerError, "failed to create cart", err.Error())
				return
			}
		} else {
			RespondError(c, http.StatusInternalServerError, "failed to get cart", err.Error())
			return
		}
	}

	var body ItemBody

	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	productId, err := uuid.Parse(body.ProductID)

	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid product ID", err.Error())
		return
	}

	product, err := repocitory.NewProductRepository().GetById(c.Request.Context(), productId)

	if err != nil {
		RespondError(c, http.StatusNotFound, "product not found", err.Error())
		return
	}

	if body.Quantity <= 0 {
		RespondError(c, http.StatusBadRequest, "invalid quantity", "Quantity must be greater than zero")
		return
	}

	cartItemRepo := repocitory.NewCartItemsRepository()

	// item already exist in cart

	exists, err = cartItemRepo.Exists(c.Request.Context(), cart.ID, body.ProductID)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to check item existence", err.Error())
		return
	}

	if exists {
		RespondError(c, http.StatusConflict, "item already exists in cart", "Item already exists in cart, updated quantity")
		return
	}

	cartItem := &model.CartItem{
		ProductId: productId,
		Quantity:  body.Quantity,
		CartId:    cart.ID,
		Price:     product.Price,
	}

	cartItem.ID = uuid.New()

	err = cartItemRepo.AddItem(c.Request.Context(), cartItem)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to add item to cart", err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, "Item added to cart successfully", cartItem)
}

// Remove item from cart
// @Summary Remove item from cart
// @Description Remove a product from the shopping cart by its item ID
// @Tags Shopping Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart Item ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/remove/{id} [delete]
func RemoveFromCart(c *gin.Context) {
	idParam := c.Param("id")

	itemId, err := uuid.Parse(idParam)

	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid item id", err.Error())
		return
	}

	err = repocitory.NewCartItemsRepository().RemoveItem(c.Request.Context(), itemId)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to remove item from cart", err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, "Item removed from cart successfully", nil)
}

// GetCartItems godoc
// @Summary Get items in the cart
// @Description Retrieve all items from the authenticated user's shopping cart
// @Tags Shopping Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "cart_id and list of items"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart [get]
func GetCartItems(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "unauthorized", "User not authenticated")
		return
	}

	claims := userClaims.(map[string]interface{})

	auth0ID := claims["sub"].(string)

	user, err := repocitory.NewUserRepository().GetByAuth0Id(c, auth0ID)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to get user", err.Error())
		return
	}

	cartRepo := repocitory.NewShoppingCartRepository()

	cart, err := cartRepo.GetShoppingCart(c.Request.Context(), user.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusOK, gin.H{"message": "Cart is empty"})
			return
		}
		RespondError(c, http.StatusInternalServerError, "failed to get cart", err.Error())
		return
	}

	items, err := repocitory.NewCartItemsRepository().GetItems(c.Request.Context(), cart.ID)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to get cart items", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart_id": cart.ID.String(), "items": items})
}

// UpdateQuantity godoc
// @Summary Update quantity of an item in the cart
// @Description Update the quantity of a specific cart item by its ID
// @Tags Shopping Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID"
// @Param item body ItemBody true "Updated quantity"
// @Success 200 {object} map[string]string "Item quantity updated successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/update/quantity/{id} [patch]
func UpdateQuantity(c *gin.Context) {
	idParam := c.Param("id")

	itemId, err := uuid.Parse(idParam)

	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid item id", err.Error())
		return
	}

	var body ItemBody

	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if body.Quantity <= 0 {
		RespondError(c, http.StatusBadRequest, "Invalid quantity", "Quantity must be greater than zero")
		return
	}

	err = repocitory.NewCartItemsRepository().UpdateQuantity(c.Request.Context(), itemId, body.Quantity)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to update item quantity", err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, "Item quantity updated successfully", nil)
}
