package controller

import (
	"github.com/ApnanJuanda/superindo/model"
	"github.com/ApnanJuanda/superindo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByIdAndName(ctx *gin.Context)
	GetByType(ctx *gin.Context)
	GetAllAfterSorting(ctx *gin.Context)
}

type ProductControllerImpl struct {
	ProductService service.ProductService
}

func NewProductControllerImpl(productService service.ProductService) *ProductControllerImpl {
	return &ProductControllerImpl{ProductService: productService}
}

func (c ProductControllerImpl) Create(ctx *gin.Context) {
	productRequest := model.ProductRequest{}
	if err := ctx.ShouldBindJSON(&productRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	productResponse, err := c.ProductService.Save(&productRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, productResponse)
}

func (c ProductControllerImpl) GetAll(ctx *gin.Context) {
	productResponses, err := c.ProductService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, productResponses)
}

func (c ProductControllerImpl) GetByIdAndName(ctx *gin.Context) {
	productIdStr := ctx.Param("productId")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	productRequest := model.GetProductRequest{}
	if err := ctx.ShouldBindJSON(&productRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	productResponses, err := c.ProductService.GetByIdAndName(productId, productRequest.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, productResponses)
}

func (c ProductControllerImpl) GetByType(ctx *gin.Context) {
	productType := ctx.Param("productType")
	productResponses, err := c.ProductService.GetByType(productType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, productResponses)
}

func (c ProductControllerImpl) GetAllAfterSorting(ctx *gin.Context) {
	sortingType := ctx.Param("sortingType")
	productResponses, err := c.ProductService.GetAllAfterSorting(sortingType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, productResponses)
}
