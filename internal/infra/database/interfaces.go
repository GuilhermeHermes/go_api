package database

import "github/GuilhermeHermes/GO_API/internal/entity"

type UserDB interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
	Exists(email string) (bool, error)
}

type ProductDB interface {
	Create(product *entity.Product) error
	FindByID(id string) (*entity.Product, error)
	FindAll(page int, limit int, sort string) ([]*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}
