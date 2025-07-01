package database

import "github/GuilhermeHermes/GO_API/internal/entity"

type ProductRepository interface {
}

type UserDB interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type ProductDB interface {
	Create(product *entity.Product) error
	FindByID(id string) (*entity.Product, error)
	FindAll() ([]*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}
