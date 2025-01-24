package repository

import "github.com/ApnanJuanda/superindo/model"

type ProductRepository interface {
	Save(productModel *model.Product) (*model.Product, error)
	GetAll() ([]*model.Product, error)
	GetByIdOrName(id int, name string) ([]*model.Product, error)
	GetByType(productType string) ([]*model.Product, error)
	GetAllAfterSorting(sortingType string) ([]*model.Product, error)
}
