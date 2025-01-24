//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ApnanJuanda/superindo/app"
	"github.com/ApnanJuanda/superindo/config"
	"github.com/ApnanJuanda/superindo/controller"
	"github.com/ApnanJuanda/superindo/repository"
	"github.com/ApnanJuanda/superindo/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func ProvideValidatorOptions() []validator.Option {
	return []validator.Option{}
}

var productSet = wire.NewSet(
	repository.NewProductRepositoryImpl,
	wire.Bind(new(repository.ProductRepository), new(*repository.ProductRepositoryImpl)),
	service.NewProductServiceImpl,
	wire.Bind(new(service.ProductService), new(*service.ProductServiceImpl)),
	controller.NewProductControllerImpl,
	wire.Bind(new(controller.ProductController), new(*controller.ProductControllerImpl)),
)

func InitializedServer() *app.Initialization {
	wire.Build(config.NewDB,
		config.NewRedisDB,
		ProvideValidatorOptions,
		validator.New,
		productSet,
		app.NewInitialization)
	return nil
}
