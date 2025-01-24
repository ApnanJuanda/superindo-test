package service

import (
	"github.com/ApnanJuanda/superindo/model"
)

type ProductService interface {
	Save(request *model.ProductRequest) (*model.ProductResponse, error)
	GetAll() ([]*model.ProductResponse, error)
	GetByIdAndName(productId int, name string) ([]*model.ProductResponse, error)
	GetByType(productType string) ([]*model.ProductResponse, error)
	GetAllAfterSorting(sortingType string) ([]*model.ProductResponse, error)
}
