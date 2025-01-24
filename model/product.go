package model

import (
	"github.com/ApnanJuanda/superindo/helper"
	"time"
)

type Product struct {
	Id          int       `gorm:"primary_key;autoIncrement;column:id"`
	Name        string    `gorm:"column:name"`
	Price       int       `gorm:"column:price"`
	ProductType string    `gorm:"column:product_type;type:enum('Buah', 'Sayuran', 'Protein', 'Snack')"`
	ExpiredDate time.Time `gorm:"column:expired_date"`
}

func (m *Product) TableName() string {
	return "products"
}

func (m *Product) FromCreateRequest(request *ProductRequest) {
	layout := "02-01-2006"
	parsedTime, err := time.Parse(layout, request.ExpiredDate)
	helper.PanicIfError(err)

	m.Name = request.Name
	m.Price = request.Price
	m.ExpiredDate = parsedTime
	m.ProductType = request.ProductType
}

func (m *Product) ToProductResponse() *ProductResponse {

	layout := "02-01-2006"
	expiredDateString := m.ExpiredDate.Format(layout)

	return &ProductResponse{
		Id:          m.Id,
		Name:        m.Name,
		Price:       m.Price,
		ProductType: m.ProductType,
		ExpiredDate: expiredDateString,
	}
}
