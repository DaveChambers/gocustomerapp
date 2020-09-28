package main

import (
	"github.com/DaveChambers/gocustomerapp/customer/delivery"
	"github.com/DaveChambers/gocustomerapp/customer/repository/gormpostgres"
	"github.com/DaveChambers/gocustomerapp/customer/usecase"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	repo := gormpostgres.NewCustomerRepository()
	defer repo.CloseConnection()

	usecase := usecase.NewCustomerUsecase(repo)

	delivery.NewHandler(usecase)

}
