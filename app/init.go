package app

import (
	"github.com/ApnanJuanda/superindo/controller"
	"github.com/ApnanJuanda/superindo/repository"
	"github.com/ApnanJuanda/superindo/service"
)

type Initialization struct {
	ProductRepository repository.ProductRepository
	ProductService    service.ProductService
	ProductController controller.ProductController
}

func NewInitialization(productRepository repository.ProductRepository,
	productService service.ProductService,
	productController controller.ProductController) *Initialization {
	return &Initialization{
		ProductRepository: productRepository,
		ProductService:    productService,
		ProductController: productController,
	}
}
