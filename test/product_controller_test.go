package test

import (
	"encoding/json"
	"github.com/ApnanJuanda/superindo/app"
	"github.com/ApnanJuanda/superindo/controller"
	"github.com/ApnanJuanda/superindo/helper"
	"github.com/ApnanJuanda/superindo/model"
	"github.com/ApnanJuanda/superindo/repository"
	"github.com/ApnanJuanda/superindo/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setupTestDB() (*gorm.DB, *redis.Client) {

	dialect := mysql.Open("root@tcp(localhost:3306)/superindo")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	helper.PanicIfError(err)

	var redisDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	return db, redisDB
}

func setupRouter(db *gorm.DB, redisDB *redis.Client) *gin.Engine {
	validate := validator.New()
	productRepository := repository.NewProductRepositoryImpl(db, redisDB)
	productService := service.NewProductServiceImpl(productRepository, validate)
	productController := controller.NewProductControllerImpl(productService)
	initialization := app.NewInitialization(productRepository, productService, productController)
	router := app.InitRouter(initialization)
	return router
}

func truncateProduct(db *gorm.DB) {
	db.Exec("TRUNCATE products")
}

func TestAddProduct(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)
	router := setupRouter(db, redisDB)

	requestBody := strings.NewReader(`{"name": "Kentang", "price": 15000, "productType": "Sayuran", "expiredDate": "26-01-2025"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/product", requestBody)
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, "Kentang", responseBody["name"])
	assert.Equal(t, "Sayuran", responseBody["productType"])
	assert.Equal(t, float64(15000), responseBody["price"])
	assert.Equal(t, "26-01-2025", responseBody["expiredDate"])
}

func TestGetAllProducts(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)

	productRepository := repository.NewProductRepositoryImpl(db, redisDB)
	parsedDate1, err := time.Parse("02-01-2006", "25-01-2025")
	parsedDate2, err := time.Parse("02-01-2006", "26-01-2025")
	helper.PanicIfError(err)
	product1, _ := productRepository.Save(&model.Product{
		Name:        "Apel",
		ProductType: "Buah",
		Price:       25000,
		ExpiredDate: parsedDate1,
	})
	product2, _ := productRepository.Save(&model.Product{
		Name:        "Wortel",
		ProductType: "Sayuran",
		Price:       10000,
		ExpiredDate: parsedDate2,
	})

	router := setupRouter(db, redisDB)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/product", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody []map[string]interface{}
	json.Unmarshal(body, &responseBody)

	productResponse1 := responseBody[1]
	productResponse2 := responseBody[0]

	assert.Equal(t, product2.Id, int(productResponse1["id"].(float64)))
	assert.Equal(t, product2.Name, productResponse1["name"])
	assert.Equal(t, product1.Id, int(productResponse2["id"].(float64)))
	assert.Equal(t, product1.Name, productResponse2["name"])
}

func TestGetProductByType(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)

	productRepository := repository.NewProductRepositoryImpl(db, redisDB)
	parsedDate1, err := time.Parse("02-01-2006", "25-01-2025")
	parsedDate2, err := time.Parse("02-01-2006", "26-01-2025")
	helper.PanicIfError(err)
	product1, _ := productRepository.Save(&model.Product{
		Name:        "Apel",
		ProductType: "Buah",
		Price:       25000,
		ExpiredDate: parsedDate1,
	})
	_, _ = productRepository.Save(&model.Product{
		Name:        "Wortel",
		ProductType: "Sayuran",
		Price:       10000,
		ExpiredDate: parsedDate2,
	})

	router := setupRouter(db, redisDB)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/product/type/Buah", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody []map[string]interface{}
	json.Unmarshal(body, &responseBody)

	productResponse1 := responseBody[0]

	assert.Equal(t, product1.Id, int(productResponse1["id"].(float64)))
	assert.Equal(t, product1.Name, productResponse1["name"])

}

func TestGetProductByIdAndName(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)

	productRepository := repository.NewProductRepositoryImpl(db, redisDB)
	parsedDate1, err := time.Parse("02-01-2006", "25-01-2025")
	helper.PanicIfError(err)
	product1, _ := productRepository.Save(&model.Product{
		Name:        "Bayam",
		ProductType: "Sayuran",
		Price:       5000,
		ExpiredDate: parsedDate1,
	})

	router := setupRouter(db, redisDB)
	requestBody := strings.NewReader(`{"name": "Bayam"}`)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/product/search/1", requestBody)
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody []map[string]interface{}
	json.Unmarshal(body, &responseBody)

	productResponse1 := responseBody[0]

	assert.Equal(t, product1.Id, int(productResponse1["id"].(float64)))
	assert.Equal(t, product1.Name, productResponse1["name"])
}

func TestGetProductSortingByName(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)

	productRepository := repository.NewProductRepositoryImpl(db, redisDB)
	parsedDate1, err := time.Parse("02-01-2006", "25-01-2025")
	parsedDate2, err := time.Parse("02-01-2006", "26-01-2025")
	helper.PanicIfError(err)
	product1, _ := productRepository.Save(&model.Product{
		Name:        "Apel",
		ProductType: "Buah",
		Price:       25000,
		ExpiredDate: parsedDate1,
	})
	product2, _ := productRepository.Save(&model.Product{
		Name:        "Wortel",
		ProductType: "Sayuran",
		Price:       10000,
		ExpiredDate: parsedDate2,
	})

	router := setupRouter(db, redisDB)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/product/sorting/name_desc", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody []map[string]interface{}
	json.Unmarshal(body, &responseBody)

	productResponse1 := responseBody[0]
	productResponse2 := responseBody[1]

	assert.Equal(t, product2.Id, int(productResponse1["id"].(float64)))
	assert.Equal(t, product2.Name, productResponse1["name"])
	assert.Equal(t, product1.Id, int(productResponse2["id"].(float64)))
	assert.Equal(t, product1.Name, productResponse2["name"])
}

func TestGetProductSortingByPrice(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)

	productRepository := repository.NewProductRepositoryImpl(db, redisDB)
	parsedDate1, err := time.Parse("02-01-2006", "25-01-2025")
	parsedDate2, err := time.Parse("02-01-2006", "26-01-2025")
	helper.PanicIfError(err)
	product1, _ := productRepository.Save(&model.Product{
		Name:        "Apel",
		ProductType: "Buah",
		Price:       25000,
		ExpiredDate: parsedDate1,
	})
	product2, _ := productRepository.Save(&model.Product{
		Name:        "Wortel",
		ProductType: "Sayuran",
		Price:       10000,
		ExpiredDate: parsedDate2,
	})

	router := setupRouter(db, redisDB)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/product/sorting/price_desc", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody []map[string]interface{}
	json.Unmarshal(body, &responseBody)

	productResponse1 := responseBody[0]
	productResponse2 := responseBody[1]

	assert.Equal(t, product1.Id, int(productResponse1["id"].(float64)))
	assert.Equal(t, product1.Name, productResponse1["name"])
	assert.Equal(t, product2.Id, int(productResponse2["id"].(float64)))
	assert.Equal(t, product2.Name, productResponse2["name"])
}

func TestGetProductSortingByDate(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)

	productRepository := repository.NewProductRepositoryImpl(db, redisDB)
	parsedDate1, err := time.Parse("02-01-2006", "25-01-2025")
	parsedDate2, err := time.Parse("02-01-2006", "26-01-2025")
	helper.PanicIfError(err)
	product1, _ := productRepository.Save(&model.Product{
		Name:        "Apel",
		ProductType: "Buah",
		Price:       25000,
		ExpiredDate: parsedDate1,
	})
	product2, _ := productRepository.Save(&model.Product{
		Name:        "Wortel",
		ProductType: "Sayuran",
		Price:       10000,
		ExpiredDate: parsedDate2,
	})

	router := setupRouter(db, redisDB)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/product/sorting/date_desc", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody []map[string]interface{}
	json.Unmarshal(body, &responseBody)

	productResponse1 := responseBody[0]
	productResponse2 := responseBody[1]

	assert.Equal(t, product2.Id, int(productResponse1["id"].(float64)))
	assert.Equal(t, product2.Name, productResponse1["name"])
	assert.Equal(t, product1.Id, int(productResponse2["id"].(float64)))
	assert.Equal(t, product1.Name, productResponse2["name"])
}

func TestGetProductByTypeFailed(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)

	productRepository := repository.NewProductRepositoryImpl(db, redisDB)
	parsedDate1, err := time.Parse("02-01-2006", "25-01-2025")
	parsedDate2, err := time.Parse("02-01-2006", "26-01-2025")
	helper.PanicIfError(err)
	_, _ = productRepository.Save(&model.Product{
		Name:        "Apel",
		ProductType: "Buah",
		Price:       25000,
		ExpiredDate: parsedDate1,
	})
	_, _ = productRepository.Save(&model.Product{
		Name:        "Wortel",
		ProductType: "Sayuran",
		Price:       10000,
		ExpiredDate: parsedDate2,
	})

	router := setupRouter(db, redisDB)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/product/type/Protein", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, "Maaf, Produk tidak ditemukan", responseBody["message"])
}

func TestAddProductFailed(t *testing.T) {
	db, redisDB := setupTestDB()
	truncateProduct(db)
	router := setupRouter(db, redisDB)

	requestBody := strings.NewReader(`{"name": "Kentang", "price": 15000, "productType": "Sayuran", "expiredDate": "26-01-2020"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/product", requestBody)
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, "Data produk tidak valid", responseBody["message"])
}
