package app

import "github.com/gin-gonic/gin"

func InitRouter(init *Initialization) *gin.Engine {
	router := gin.New()
	router.POST("/product", init.ProductController.Create)
	router.GET("/product/search/:productId", init.ProductController.GetByIdAndName)
	router.GET("/product/type/:productType", init.ProductController.GetByType)
	router.GET("/product/sorting/:sortingType", init.ProductController.GetAllAfterSorting)
	router.GET("/product", init.ProductController.GetAll)
	return router
}
