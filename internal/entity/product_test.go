package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	name        = "Test Product"
	description = "This is a test product"
	price       = 10.99
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct(name, description, price)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, name, product.Name)
	assert.Equal(t, description, product.Description)
	assert.Equal(t, price, product.Price)
	assert.NotEmpty(t, product.ID)       // ID should be generated
	assert.NotZero(t, product.CreatedAt) // CreatedAt should be set
	assert.NotZero(t, product.UpdatedAt) // UpdatedAt should be set
}

func TestProductWhenNameIsRequired(t *testing.T) {
	product, err := NewProduct("", description, price)
	assert.Nil(t, product)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	product, err := NewProduct(name, description, -1)
	assert.Nil(t, product)
	assert.NotNil(t, err)
	assert.Equal(t, ErrPriceMustBePositive, err)
}
