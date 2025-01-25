package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ApnanJuanda/superindo/helper"
	"github.com/ApnanJuanda/superindo/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"slices"
	"strings"
	"time"
)

type ProductRepositoryImpl struct {
	DB      *gorm.DB
	RedisDB *redis.Client
}

var ctx = context.Background()

func NewProductRepositoryImpl(DB *gorm.DB, RedisDB *redis.Client) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{
		DB:      DB,
		RedisDB: RedisDB,
	}
}

func (r *ProductRepositoryImpl) Save(productModel *model.Product) (*model.Product, error) {
	err := r.DB.Create(productModel).Error
	if err != nil {
		return nil, err
	}

	productData, err := json.Marshal(productModel)
	if err != nil {
		return nil, err
	}

	productKey := fmt.Sprintf("product:%d", productModel.Id)
	if productModel.Id > 0 {
		err = r.RedisDB.Set(ctx, productKey, productData, time.Hour).Err()
		if err != nil {
			return nil, err
		}
	}

	if productModel.Id > 0 {
		return productModel, nil
	}
	return nil, errors.New("Terjadi kesalahan saat menyimpan data produk")
}

func (r *ProductRepositoryImpl) GetAll() ([]*model.Product, error) {
	var products []*model.Product
	products, err := r.getAllProductsFromRedis()
	if products == nil {
		err = r.DB.Find(&products).Error
		helper.PanicIfError(err)

		if len(products) > 0 {
			r.saveToRedis(products)
		}
	}
	if len(products) > 0 {
		return products, nil
	}
	return nil, errors.New("Maaf, Produk tidak ditemukan")
}

func (r *ProductRepositoryImpl) getAllProductsFromRedis() ([]*model.Product, error) {
	var products []*model.Product
	var cursor uint64

	for {
		keys, nextCursor, err := r.RedisDB.Scan(ctx, cursor, "product:*", 0).Result()
		if err != nil {
			return nil, err
		}
		cursor = nextCursor

		for _, key := range keys {
			productData, err := r.RedisDB.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				return nil, err
			}

			var product model.Product
			err = json.Unmarshal([]byte(productData), &product)
			if err != nil {
				return nil, err
			}

			products = append(products, &product)
		}
		if cursor == 0 {
			break
		}
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetByIdOrName(id int, name string) ([]*model.Product, error) {
	var products []*model.Product
	if id > 0 || name != "" {
		products, _ = r.GetByIdOrNameFromRedis(id, name)
		if products == nil {
			if id > 0 {
				err := r.DB.Find(&products, "id = ?", id).Error
				helper.PanicIfError(err)
			}

			if len(products) <= 0 && name != "" {
				err := r.DB.Find(&products, "name = ?", name).Error
				helper.PanicIfError(err)
			}
		}
		if len(products) > 0 {
			r.saveToRedis(products)
		}
	}

	if len(products) > 0 {
		return products, nil
	}
	return nil, errors.New("Maaf, Produk tidak ditemukan")
}

func (r *ProductRepositoryImpl) GetByIdOrNameFromRedis(id int, name string) ([]*model.Product, error) {
	var products []*model.Product
	if id > 0 {
		var productKey = fmt.Sprintf("product:%d", id)
		productData, err := r.RedisDB.Get(ctx, productKey).Result()
		if err != nil {
			return nil, err
		}

		var product model.Product
		err = json.Unmarshal([]byte(productData), &product)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if len(products) <= 0 && name != "" {
		var productsFilter []*model.Product

		products, _ = r.getAllProductsFromRedis()
		if len(products) > 0 {
			for _, value := range products {
				if value.Name == name {
					productsFilter = append(productsFilter, value)
				}
			}
		}
		if len(productsFilter) > 0 {
			return productsFilter, nil
		} else {
			return nil, errors.New("Produk tidak ditemukan")
		}
	}

	return products, nil
}

func (r *ProductRepositoryImpl) GetByType(productType string) ([]*model.Product, error) {
	var products []*model.Product
	products, err := r.GetByTypeFromRedis(productType)
	if products == nil {
		err = r.DB.Where("product_type = ?", productType).Find(&products).Error
		helper.PanicIfError(err)

		if len(products) > 0 {
			r.saveToRedis(products)
		}
	}
	if len(products) > 0 {
		return products, nil
	}
	return nil, errors.New("Maaf, Produk tidak ditemukan")
}

func (r *ProductRepositoryImpl) GetByTypeFromRedis(productType string) ([]*model.Product, error) {
	var products []*model.Product
	keys, err := r.RedisDB.Keys(ctx, "product:*").Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		productData, err := r.RedisDB.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var product model.Product
		err = json.Unmarshal([]byte(productData), &product)
		if err != nil {
			return nil, err
		}
		if product.ProductType == productType {
			products = append(products, &product)
		}
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetAllAfterSorting(sortingType string) ([]*model.Product, error) {
	var products []*model.Product
	params := strings.Split(sortingType, "_")
	var sortingTypeValue string
	switch params[0] {
	case "date":
		params[0] = "expired_date"
	case "price":
		params[0] = "price"
	case "name":
		params[0] = "name"
	default:
		params[0] = ""
	}
	if params[0] != "" {
		sortingTypeValue = strings.Join(params, " ")
	}

	products, err := r.getAllProductsFromRedis()
	if products == nil {
		err = r.DB.Order(sortingTypeValue).Find(&products).Error
		helper.PanicIfError(err)
	} else {
		r.sortingDataFromRedis(params[0], params[1], products)
	}

	if len(products) > 0 {
		return products, nil
	}
	return nil, errors.New("Maaf, Produk tidak ditemukan")
}

func (r *ProductRepositoryImpl) sortingDataFromRedis(sorting string, sortingType string, products []*model.Product) {
	if sorting == "expired_date" {
		if sortingType == "asc" {
			slices.SortFunc(products, func(a, b *model.Product) int {
				if a.ExpiredDate.Before(b.ExpiredDate) {
					return -1
				} else if a.ExpiredDate.After(b.ExpiredDate) {
					return 1
				}
				return 0
			})
		} else if sortingType == "desc" {
			slices.SortFunc(products, func(a, b *model.Product) int {
				if a.ExpiredDate.After(b.ExpiredDate) {
					return -1
				} else if a.ExpiredDate.Before(b.ExpiredDate) {
					return 1
				}
				return 0
			})
		}
	} else {
		if sorting == "name" {
			if sortingType == "asc" {
				slices.SortFunc(products, func(a, b *model.Product) int {
					if a.Name < b.Name {
						return -1
					} else if a.Name > b.Name {
						return 1
					}
					return 0
				})
			} else if sortingType == "desc" {
				slices.SortFunc(products, func(a, b *model.Product) int {
					if a.Name < b.Name {
						return 1
					} else if a.Name > b.Name {
						return -1
					}
					return 0
				})
			}
		} else if sorting == "price" {
			if sortingType == "asc" {
				slices.SortFunc(products, func(a, b *model.Product) int {
					if a.Price < b.Price {
						return -1
					} else if a.Price > b.Price {
						return 1
					}
					return 0
				})
			} else if sortingType == "desc" {
				slices.SortFunc(products, func(a, b *model.Product) int {
					if a.Price < b.Price {
						return 1
					} else if a.Price > b.Price {
						return -1
					}
					return 0
				})
			}
		}
	}
}

func (r *ProductRepositoryImpl) saveToRedis(products []*model.Product) {
	for _, product := range products {
		productData, err := json.Marshal(product)
		helper.PanicIfError(err)
		if product.Id > 0 {
			productKey := fmt.Sprintf("product:%d", product.Id)
			err = r.RedisDB.Set(ctx, productKey, productData, time.Hour).Err()
			helper.PanicIfError(err)
		}
	}
}
