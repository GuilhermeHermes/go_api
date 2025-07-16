package main

import (
	"fmt"
	"net/http"

	"github/GuilhermeHermes/GO_API/configs"
	"github/GuilhermeHermes/GO_API/internal/entity"
	"github/GuilhermeHermes/GO_API/internal/infra/webserver"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Auto migrate
	db.AutoMigrate(&entity.User{}, &entity.Product{})

	// Setup routes
	router := webserver.SetupRoutes(db)

	fmt.Printf("Server starting on port %s\n", cfg.WebServerPort)
	http.ListenAndServe(":"+cfg.WebServerPort, router)
}
