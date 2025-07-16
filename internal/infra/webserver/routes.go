package webserver

import (
	"github/GuilhermeHermes/GO_API/configs"
	"github/GuilhermeHermes/GO_API/internal/infra/database"
	"github/GuilhermeHermes/GO_API/internal/infra/webserver/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *chi.Mux {

	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Repositories
	productRepo := database.NewProductRepository(db)
	userRepo := database.NewUserRepository(db)

	// Handlers
	productHandler := handlers.NewProductHandler(productRepo)
	userHandler := handlers.NewUserHandler(userRepo, cfg.TokenAuth, cfg.JwtExpiration)

	// Product routes
	r.Route("/products", func(r chi.Router) {
		r.Post("/", productHandler.CreateProduct)       // POST /products
		r.Get("/", productHandler.GetAllProducts)       // GET /products?page=1&limit=10&sort=asc
		r.Get("/{id}", productHandler.GetProduct)       // GET /products/{id}
		r.Put("/{id}", productHandler.UpdateProduct)    // PUT /products/{id}
		r.Delete("/{id}", productHandler.DeleteProduct) // DELETE /products/{id}
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)                 // POST /users
		r.Get("/email/{email}", userHandler.GetUserByEmail) // GET /users/email/{email}
		r.Get("/{id}", userHandler.GetUserByID)             // GET /users/{id}
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
		r.Post("/generate-jwt", userHandler.GetJwt) // POST /users/generate-jwt
	})

	return r
}
