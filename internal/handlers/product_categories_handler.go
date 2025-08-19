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

type CreateCategoryBody struct {
	Name     string     `json:"name" binding:"required"`
	ParentId *uuid.UUID `json:"parentId"`
}

// CreateCategory godoc
// @Summary Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param body body CreateCategoryBody true "Category body"
// @Success 201 {object} model.ProductCategory
// @Failure 400 {object} map[string]string
// @Router /products/categories/create [post]
func AddProductCategory(c *gin.Context) {
	var body CreateCategoryBody

	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	category := &model.ProductCategory{
		Name:     strings.ToLower(body.Name),
		ParentId: body.ParentId,
	}
	category.ID = uuid.New()

	err := repocitory.NewCategoryRepository().Create(c, category)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failled to create category", err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, "Category added successfully", category)
}

// ListCategories godoc
// @Summary List product categories
// @Description Returns paginated list of product categories
// @Tags categories
// @Produce json
// @Param limit query int false "Number of categories to return" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} model.ProductCategory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/categories [get]
func ListCategories(c *gin.Context) {
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

	categoryRepo := repocitory.NewCategoryRepository()
	categories, err := categoryRepo.List(c.Request.Context(), limit, offset)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failled to fetch categories", err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, "Success", categories)
}

type updateCategoryBody struct {
	Name     *string    `json:"name,omitempty"`
	ParentId *uuid.UUID `json:"parentId,omitempty"`
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body updateCategoryBody true "Category data"
// @Success 200 {object} model.ProductCategory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/categories/{id} [patch]
func UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")

	categoryID, err := uuid.Parse(idParam)

	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invaluid category id", err.Error())

		return
	}

	var body updateCategoryBody
	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request body", err.Error())

		return
	}

	categoryRepo := repocitory.NewCategoryRepository()

	category, err := categoryRepo.GetById(c.Request.Context(), categoryID)
	if err != nil {
		RespondError(c, http.StatusNotFound, "Category not found", err.Error())
		return
	}

	if body.Name != nil {
		category.Name = *body.Name
	}

	if body.ParentId != nil {
		category.ParentId = body.ParentId
	}

	if err := categoryRepo.Update(c.Request.Context(), category); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to update category", err.Error())

		return
	}

	RespondSuccess(c, http.StatusOK, "category updated successfully", category)
}
