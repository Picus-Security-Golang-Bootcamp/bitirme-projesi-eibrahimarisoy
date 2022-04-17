package auth

import (
	"database/sql"
	"fmt"
	"patika-ecommerce/internal/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	firstName = "test"
	lastName  = "test"
	username  = "test"
	email     = "test@email.com"
	password  = "test"
	id        = uuid.New()
)

var u = model.User{
	Base:      model.Base{ID: id},
	FirstName: &firstName,
	LastName:  &lastName,
	Username:  &username,
	Email:     &email,
	Password:  password,
	IsAdmin:   false,
}

func TestUserRepository_GetUser(t *testing.T) {
	db, mock := NewMock()
	repo := &UserRepository{db}

	// query := "SELECT id, first_name, last_name, username, email, is_admin FROM users WHERE id = \\?"

	query := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_email", "username", "email", "is_admin"}).
		AddRow(u.ID, u.FirstName, u.LastName, u.Username, u.Email, u.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnRows(rows)

	user, _ := repo.GetUser(u.ID.String())

	assert.Equal(t, user.ID, id)
	assert.Equal(t, *user.Username, username)

}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	db, mock := NewMock()
	repo := &UserRepository{db}

	query := `SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT 1`

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_email", "username", "email", "is_admin"}).
		AddRow(u.ID, u.FirstName, u.LastName, u.Username, u.Email, u.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(*u.Email).WillReturnRows(rows)

	user, _ := repo.GetUserByEmail(*u.Email)
	fmt.Println(user)
	assert.Equal(t, user.ID, id)
	assert.Equal(t, *user.Username, username)
}

// func TestUserRepository_InserUser(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &UserRepository{db}

// 	query := `INSERT INTO "users" ("created_at", "updated_at", "deleted_at", "first_name", "last_name", "username", "email", "password", "is_admin") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING "id"`

// 	// rows := sqlmock.NewRows([]string{"id", "first_name", "last_email", "username", "email", "is_admin"}).
// 	// 	AddRow(u.ID, u.FirstName, u.LastName, u.Username, u.Email, u.IsAdmin)
// 	mock.ExpectBegin()
// 	prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
// 	prep.ExpectExec().WithArgs(
// 		time.Now(),
// 		time.Now(),
// 		nil,
// 		*u.FirstName,
// 		*u.LastName,
// 		*u.Username,
// 		*u.Email,
// 		u.Password,
// 		u.IsAdmin,
// 	).WillReturnResult(sqlmock.NewResult(0, 1))
// 	mock.ExpectCommit()
// 	_, err := repo.InsertUser(&u)
// 	assert.NoError(t, err)

// }
