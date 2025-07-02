package database

import (
	"testing"

	"github/GuilhermeHermes/GO_API/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProductTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&entity.Product{})
	require.NoError(t, err)

	return db
}

func createTestProduct(t *testing.T) *entity.Product {
	product, err := entity.NewProduct("testproduct", "descricao x", 100.0)
	require.NoError(t, err)
	return product
}

func TestNewProductRepository(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := NewProductRepository(db)

	require.NotNil(t, productRepo)
	require.Equal(t, db, productRepo.DB)
}

func TestProduct_Create(t *testing.T) {
	t.Run("should create product successfully", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		product := createTestProduct(t)
		err := productRepo.Create(product)
		require.NoError(t, err)
		assert.NotEmpty(t, product.ID)
		assert.NotEmpty(t, product.CreatedAt)
		assert.NotEmpty(t, product.UpdatedAt)

		// Verify product was actually saved to database
		var count int64
		err = db.Model(&entity.Product{}).Where("id = ?", product.ID).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		foundProduct, err := productRepo.FindByID(product.ID.String())
		require.NoError(t, err)
		require.NotNil(t, foundProduct)
		assert.Equal(t, product.ID, foundProduct.ID)
		assert.Equal(t, product.Name, foundProduct.Name)
		assert.Equal(t, product.Description, foundProduct.Description)
		assert.Equal(t, product.Price, foundProduct.Price)
	})
	t.Run("should return error when product is nil", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		err := productRepo.Create(nil)
		require.Error(t, err)
		assert.Equal(t, entity.ErrIdIsRequired, err)

	})
	t.Run("should return error when product name is empty", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		product := createTestProduct(t)
		product.Name = ""
		err := productRepo.Create(product)
		require.Error(t, err)
		assert.Equal(t, entity.ErrNameIsRequired, err)
	})
}

func TestProduct_FindByID(t *testing.T) {
	t.Run("should find product by ID successfully", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		product := createTestProduct(t)
		err := productRepo.Create(product)
		require.NoError(t, err)

		foundProduct, err := productRepo.FindByID(product.ID.String())
		require.NoError(t, err)
		require.NotNil(t, foundProduct)
		assert.Equal(t, product.ID, foundProduct.ID)
		assert.Equal(t, product.Name, foundProduct.Name)
		assert.Equal(t, product.Description, foundProduct.Description)
		assert.Equal(t, product.Price, foundProduct.Price)
	})
	t.Run("should return error when ID is empty", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		_, err := productRepo.FindByID("")
		assert.Error(t, err)
		assert.Equal(t, "id cannot be empty", err.Error())
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)

		foundProduct, err := productRepo.FindByID("non-existent-id")
		assert.Error(t, err)
		assert.Nil(t, foundProduct)
	})
}

func TestProduct_FindAll(t *testing.T) {

	t.Run("should find all products with pagination and sorting", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		products := []*entity.Product{
			createTestProduct(t),
			createTestProduct(t),
			createTestProduct(t),
		}
		for _, product := range products {
			err := productRepo.Create(product)
			require.NoError(t, err)
		}
		foundProducts, err := productRepo.FindAll(1, 10, "asc")
		require.NoError(t, err)
		require.Len(t, foundProducts, 3)
		assert.Equal(t, products[0].Name, foundProducts[0].Name)
		assert.Equal(t, products[1].Name, foundProducts[1].Name)
		assert.Equal(t, products[2].Name, foundProducts[2].Name)
	})

	t.Run("should return empty slice when no products found", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)

		foundProducts, err := productRepo.FindAll(1, 10, "asc")
		require.NoError(t, err)
		assert.Empty(t, foundProducts)
	})

	t.Run("should return error for invalid sort order", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)

		_, err := productRepo.FindAll(1, 10, "invalid")
		assert.Error(t, err)
		assert.Equal(t, "sort must be 'asc' or 'desc'", err.Error())
	})
	t.Run("should return error for invalid pagination", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)

		_, err := productRepo.FindAll(0, 10, "asc")
		assert.Error(t, err)
		assert.Equal(t, "page and limit must be greater than 0", err.Error())
	})
	t.Run("should return error for invalid limit", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)

		_, err := productRepo.FindAll(1, 0, "asc")
		assert.Error(t, err)
		assert.Equal(t, "page and limit must be greater than 0", err.Error())
	})
	t.Run("should return error for invalid page and limit", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)

		_, err := productRepo.FindAll(0, 0, "asc")
		assert.Error(t, err)
		assert.Equal(t, "page and limit must be greater than 0", err.Error())
	})
	t.Run("should return products with default pagination and sorting", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		products := []*entity.Product{
			createTestProduct(t),
			createTestProduct(t),
			createTestProduct(t),
		}
		for _, product := range products {
			err := productRepo.Create(product)
			require.NoError(t, err)
		}

		foundProducts, err := productRepo.FindAll(1, 10, "")
		require.NoError(t, err)
		require.Len(t, foundProducts, 3)
		assert.Equal(t, products[0].Name, foundProducts[0].Name)
		assert.Equal(t, products[1].Name, foundProducts[1].Name)
		assert.Equal(t, products[2].Name, foundProducts[2].Name)
	})
}

func TestProduct_Update(t *testing.T) {
	t.Run("should update product successfully", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		product := createTestProduct(t)
		err := productRepo.Create(product)
		require.NoError(t, err)

		product.Name = "updated name"
		err = productRepo.Update(product)
		require.NoError(t, err)

		foundProduct, err := productRepo.FindByID(product.ID.String())
		require.NoError(t, err)
		assert.Equal(t, "updated name", foundProduct.Name)
	})
}

func TestProduct_Delete(t *testing.T) {
	t.Run("should delete product successfully", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)
		product := createTestProduct(t)
		err := productRepo.Create(product)
		require.NoError(t, err)

		err = productRepo.Delete(product.ID.String())
		require.NoError(t, err)

		foundProduct, err := productRepo.FindByID(product.ID.String())
		assert.Error(t, err)
		assert.Nil(t, foundProduct)
	})
	t.Run("should return error when deleting non-existent product", func(t *testing.T) {
		db := setupProductTestDB(t)
		productRepo := NewProductRepository(db)

		err := productRepo.Delete("non-existent-id")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
