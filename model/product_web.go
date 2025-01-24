package model

type ProductRequest struct {
	Name        string `validate:"required,max=200,min=1" json:"name"`
	Price       int    `validate:"required" json:"price"`
	ProductType string `validate:"required,max=200,min=1" json:"productType"`
	ExpiredDate string `validate:"required" json:"expiredDate"`
}

type GetProductRequest struct {
	Name        string `validate:"max=200,min=1" json:"name"`
	ProductType string `validate:"max=200,min=1" json:"productType"`
}

type ProductResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	ProductType string `json:"productType"`
	ExpiredDate string `json:"expiredDate"`
}
