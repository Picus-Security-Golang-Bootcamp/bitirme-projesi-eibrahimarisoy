package category

import (
	"database/sql"
	"patika-ecommerce/internal/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
	name        = "test"
	id          = uuid.New()
	description = "test"
)

var c = model.Category{
	Base:        model.Base{ID: id},
	Name:        &name,
	Description: description,
}

func TestCategoryRepository_GetCategoryByID(t *testing.T) {
	db, mock := NewMock()
	repo := &CategoryRepository{db}

	// query := "SELECT id, first_name, last_name, username, email, is_admin FROM users WHERE id = \\?"

	query := `SELECT * FROM "categories" WHERE id = $1 ORDER BY "categories"."id" LIMIT 1`

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(c.ID, c.Name, c.Description)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(c.ID).WillReturnRows(rows)

	category, _ := repo.GetCategoryByID(c.ID)

	assert.Equal(t, category.ID, id)
	assert.Equal(t, category.Name, name)

}

func TestCategoryRepository_GetCategories(t *testing.T) {
	db, mock := NewMock()
	repo := &CategoryRepository{db}

	query := `SELECT * FROM "categories"`

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(c.ID, c.Name, c.Description)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	_, err := repo.GetCategories()

	assert.Equal(t, err, nil)

}
