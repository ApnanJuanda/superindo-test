package service

import (
	"errors"
	"github.com/ApnanJuanda/superindo/helper"
	"github.com/ApnanJuanda/superindo/model"
	"github.com/ApnanJuanda/superindo/repository"
	"github.com/go-playground/validator/v10"
	"time"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
	Validator         *validator.Validate
}

func NewProductServiceImpl(productRepository repository.ProductRepository, validator *validator.Validate) *ProductServiceImpl {
	return &ProductServiceImpl{ProductRepository: productRepository, Validator: validator}
}

func (s *ProductServiceImpl) Save(request *model.ProductRequest) (*model.ProductResponse, error) {
	var productModel = new(model.Product)
	var productResponse *model.ProductResponse

	err := s.Validator.Struct(request)
	helper.PanicIfError(err)

	expiredDateParse, err := time.Parse("02-01-2006", request.ExpiredDate)
	if expiredDateParse.Before(time.Now()) {
		return nil, errors.New("Data produk tidak valid")
	}

	productModel.FromCreateRequest(request)

	_, err = s.ProductRepository.Save(productModel)
	if err != nil {
		return nil, errors.New("Terjadi kesalahan saat menyimpan data")
	}
	productResponse = productModel.ToProductResponse()

	return productResponse, nil
}

func (s *ProductServiceImpl) GetAll() ([]*model.ProductResponse, error) {
	var productResponses []*model.ProductResponse
	productModels, err := s.ProductRepository.GetAll()
	if err != nil {
		return nil, err
	}
	for _, product := range productModels {
		var productResponse *model.ProductResponse
		productResponse = product.ToProductResponse()
		productResponses = append(productResponses, productResponse)
	}

	return productResponses, nil
}

func (s *ProductServiceImpl) GetByIdAndName(productId int, name string) ([]*model.ProductResponse, error) {
	var productResponses []*model.ProductResponse
	productModels, err := s.ProductRepository.GetByIdOrName(productId, name)
	if err != nil {
		return nil, err
	}
	for _, product := range productModels {
		var productResponse *model.ProductResponse
		productResponse = product.ToProductResponse()
		productResponses = append(productResponses, productResponse)
	}
	return productResponses, nil
}

func (s *ProductServiceImpl) GetByType(productType string) ([]*model.ProductResponse, error) {
	var productResponses []*model.ProductResponse
	productModels, err := s.ProductRepository.GetByType(productType)
	if err != nil {
		return nil, err
	}
	for _, product := range productModels {
		var productResponse *model.ProductResponse
		productResponse = product.ToProductResponse()
		productResponses = append(productResponses, productResponse)
	}
	return productResponses, nil
}

func (s *ProductServiceImpl) GetAllAfterSorting(sortingType string) ([]*model.ProductResponse, error) {
	var productResponses []*model.ProductResponse
	productModels, err := s.ProductRepository.GetAllAfterSorting(sortingType)
	if err != nil {
		return nil, err
	}
	for _, product := range productModels {
		var productResponse *model.ProductResponse
		productResponse = product.ToProductResponse()
		productResponses = append(productResponses, productResponse)
	}

	return productResponses, nil
}
