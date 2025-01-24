package config

import (
	"fmt"
	"github.com/ApnanJuanda/superindo/helper"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func NewDB() *gorm.DB {
	err := godotenv.Load("environment/.env")
	helper.PanicIfError(err)

	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYUSERNAME"),
		os.Getenv("MYPASSWORD"),
		os.Getenv("MYHOST"),
		os.Getenv("MYPORT"),
		os.Getenv("MYDATABASE"))

	dialect := mysql.Open(mysqlInfo)
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	helper.PanicIfError(err)

	return db
}

func NewRedisDB() *redis.Client {
	var redisDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	return redisDB
}
