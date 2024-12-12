package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"authentication_medods/cmd/api/config"
	"authentication_medods/cmd/api/controller"
	"authentication_medods/cmd/api/service"
	"authentication_medods/cmd/api/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := &config.Config{}
	cfg.Init()

	dns := fmt.Sprintf(
		"host=postgres user=%s password=%s  dbname=%s port=5432 sslmode=disable",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Name,
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var (
		str    = storage.NewStorage(db)
		svc    = service.NewService(str)
		ctrl   = controller.NewAPIController(svc)
		router = gin.Default()
	)

	ctrl.SetupRouter(router)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router.Handler(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
