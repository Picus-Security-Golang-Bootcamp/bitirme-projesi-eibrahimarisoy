package product

import (
	"database/sql"
	"patika-ecommerce/internal/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-openapi/strfmt"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMock() (DB *gorm.DB, mock sqlmock.Sqlmock) {
	var (
		db *sql.DB
	)

	db, mock, _ = sqlmock.New()

	DB, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	return DB, mock
}

var (
	id          = uuid.New()
	name        = "test"
	description = "test"
	price       = float64(100)
	stock       = int64(10)
	sku         = "test"
)

var c = model.Product{
	Base:         model.Base{ID: id},
	Name:         &name,
	Description:  description,
	Price:        price,
	Stock:        &stock,
	SKU:          &sku,
	Categories:   []model.Category{},
	CategoriesID: []strfmt.UUID{},
}

func TestCategoryRepository_GetProductWithoutCategories(t *testing.T) {
	db, mock := NewMock()
	repo := &ProductRepository{db}

	// query := "SELECT id, first_name, last_name, username, email, is_admin FROM users WHERE id = \\?"

	query := `SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT 1`

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "sku"}).
		AddRow(c.ID, c.Name, c.Description, c.Price, c.Stock, c.SKU)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(c.ID).WillReturnRows(rows)

	category, _ := repo.GetProductWithoutCategories(c.ID)

	assert.Equal(t, category.ID, id)
	assert.Equal(t, category.Name, name)

}
