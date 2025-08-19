package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/Oj-washingtone/savannah-store/internal/repocitory"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createProductBody struct {
	Name        string
	CategoryId  uuid.UUID
	Description string
	Price       int64
	Stock       int
}

// CreateProduct godoc
// @Summary Add a new product
// @Tags products
// @Accept json
// @Produce json
// @Param body body createProductBody true "Product body"
// @Success 201 {object} model.Product
// @Failure 400 {object} map[string]string
// @Router /products/create [post]
func CreateProduct(c *gin.Context) {
	var body createProductBody

	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	product := &model.Product{
		Name:        strings.ToLower(body.Name),
		CategoryID:  body.CategoryId,
		Description: body.Description,
		Price:       body.Price,
		Stock:       body.Stock,
	}

	product.ID = uuid.New()

	productRepository := repocitory.NewProductRepository()

	err := productRepository.Create(c, product)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failled to create product", err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, "Product Added successfully", product)

}

// GetProductById godoc
// @Summary Get product by ID
// @Description get product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product Id"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
func GetProductById(c *gin.Context) {
	idParam := c.Param("id")

	productId, err := uuid.Parse(idParam)

	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid product id", err.Error())
		return
	}

	product, err := repocitory.NewProductRepository().GetById(c, productId)

	if err != nil {
		RespondError(c, http.StatusNotFound, "Failed to fetch product", err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, "success", product)
}

// ListProducts godoc
// @Summary List product
// @Description Returns paginated list of product
// @Tags products
// @Produce json
// @Param limit query int false "Number of products to return" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} model.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [get]
func ListProducts(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil || limit <= 0 {
		RespondError(c, http.StatusBadRequest, "Invalid limit", err.Error())
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if err != nil || offset < 0 {
		RespondError(c, http.StatusBadRequest, "Invalid offset", err.Error())
		return
	}

	products, err := repocitory.NewProductRepository().List(c, limit, offset)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failled to fetch categories", err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, "Products found", products)

}

func UpdateProduct(c *gin.Context) {

}

func DeleteProduct(c *gin.Context) {

}
