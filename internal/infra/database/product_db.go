package database

import (
	"errors"
	"strings"

	"github/GuilhermeHermes/GO_API/internal/entity"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Create(product *entity.Product) error {
	if product == nil {
		return entity.ErrIdIsRequired
	}
	if strings.TrimSpace(product.Name) == "" {
		return entity.ErrNameIsRequired
	}
	if product.Price <= 0 {
		return entity.ErrPriceIsRequired
	}

	return r.DB.Create(product).Error
}

func (p *ProductRepository) FindByID(id string) (*entity.Product, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id cannot be empty")
	}

	var product entity.Product
	if err := p.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *ProductRepository) FindAll(page int, limit int, sort string) ([]*entity.Product, error) {
	var products []*entity.Product

	if sort != "" && sort != "asc" && sort != "desc" {
		return nil, errors.New("sort must be 'asc' or 'desc'")
	}

	if page <= 0 || limit <= 0 {
		return nil, errors.New("page and limit must be greater than 0")
	}

	if sort == "" {
		sort = "asc"
	}

	query := p.DB.Order("created_at " + sort)

	offset := (page - 1) * limit
	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductRepository) Update(product *entity.Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	_, err := p.FindByID(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(product).Error
}

func (p *ProductRepository) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id cannot be empty")
	}
	return p.DB.Delete(&entity.Product{}, "id = ?", id).Error
}
