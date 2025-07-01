package entity

import (
	"errors"
	"github/GuilhermeHermes/GO_API/pkg/entity"
	"time"
)

var (
	ErrIdIsRequired        = errors.New("product ID is required")
	ErrInvalidID           = errors.New("invalid product ID")
	ErrNameIsRequired      = errors.New("product name is required")
	ErrPriceIsRequired     = errors.New("product price is required")
	ErrPriceMustBePositive = errors.New("product price must be positive")
)

type Product struct {
	ID          entity.ID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIdIsRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidID
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price <= 0 {
		return ErrPriceMustBePositive
	}
	if p.Price < 0 {
		return ErrPriceIsRequired
	}
	return nil
}

func NewProduct(name, description string, price float64) (*Product, error) {
	product := &Product{
		ID:          entity.NewID(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := product.Validate(); err != nil {
		return nil, err
	}
	return product, nil
}
