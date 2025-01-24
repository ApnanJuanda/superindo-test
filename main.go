package main

import (
	"fmt"
	"github.com/ApnanJuanda/superindo/app"
	"github.com/ApnanJuanda/superindo/config"
	"github.com/ApnanJuanda/superindo/helper"
	"github.com/ApnanJuanda/superindo/model"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func main() {
	//SeederData()

	err := godotenv.Load("environment/.env")
	helper.PanicIfError(err)

	port := os.Getenv("PORT")
	init := InitializedServer()

	app := app.InitRouter(init)
	fmt.Println("My application is running")
	app.Run(":" + port)
}

func SeederData() {
	// execute first migrate -database "mysql://root@tcp(localhost:3306)/superindo" -path migrations up

	db := config.NewDB()

	var products = []model.Product{
		{Name: "Nanas Sunpride", Price: 20000, ExpiredDate: time.Date(2025, time.January, 30, 0, 0, 0, 0, time.UTC), ProductType: "Buah"},
		{Name: "Bayam", Price: 5000, ExpiredDate: time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC), ProductType: "Sayuran"},
		{Name: "Pisang Sunpride", Price: 15000, ExpiredDate: time.Date(2025, time.February, 02, 0, 0, 0, 0, time.UTC), ProductType: "Buah"},
		{Name: "Ikan Tuna", Price: 40000, ExpiredDate: time.Date(2025, time.January, 28, 0, 0, 0, 0, time.UTC), ProductType: "Protein"},
		{Name: "Chitato", Price: 20000, ExpiredDate: time.Date(2025, time.July, 10, 0, 0, 0, 0, time.UTC), ProductType: "Snack"},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			fmt.Println("Error seeding data: ", err)
		}
	}
	fmt.Println("Success seeding data")
}
