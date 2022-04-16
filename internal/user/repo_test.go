package auth

import (
	"database/sql"
	"fmt"
	"patika-ecommerce/internal/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository AuthRepository
	person     *model.User
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(s.T(), err) //sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(s.T(), err)

	// s.DB.LogMode(true)

	s.repository = NewUserRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	// require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

// func (s *Suite) TestUserRepository_InsertUser() {

// 	FirstName := "test"
// 	LastName := "test"
// 	Username := "test"
// 	Email := "test@email.com"
// 	Password := "test"
// 	id := uuid.New()
// 	user := &model.User{
// 		Base:      model.Base{ID: id},
// 		FirstName: &FirstName,
// 		LastName:  &LastName,
// 		Username:  &Username,
// 		Email:     &Email,
// 		Password:  Password,
// 		IsAdmin:   false,
// 	}
// 	fmt.Println("id", id)
// 	s.mock.ExpectBegin()
// 	sqlInsert := `INSERT INTO "users" ("created_at","updated_at","deleted_at","first_name","last_name","username","email","password","is_admin","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "users"."id"`

// 	s.mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
// 		WithArgs(time.Now(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.FirstName, user.LastName, user.Username, user.Email, user.Password, user.IsAdmin, id).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

// 	_, err := s.repository.InsertUser(user)

// 	require.NoError(s.T(), err)

// }

func (s *Suite) TestUserRepository_GetUser() {

	var (
		id       = uuid.New()
		username = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(id.String()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
			AddRow(id.String(), username))

	user, err := s.repository.GetUser(id.String())

	require.NoError(s.T(), err)
	assert.Equal(s.T(), user.ID, id)
	assert.Equal(s.T(), user.Username, username)

	// s.mock.ExpectCommit()

}

func (s *Suite) TestUserRepository_GetAll() {

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(uuid.New().String()).
			AddRow(uuid.New().String()))

	users, err := s.repository.GetAll()
	fmt.Println(users)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(*users), 2)

}

// func (s *Suite) TestUserRepository_DeleteUser() {

// 	id := uuid.New()
// 	// FirstName := "test"
// 	// LastName := "test"
// 	// Username := "test"
// 	// Email := "test@email.com"
// 	// Password := "test"
// 	// user := &model.User{
// 	// 	Base:      model.Base{ID: id},
// 	// 	FirstName: &FirstName,
// 	// 	LastName:  &LastName,
// 	// 	Username:  &Username,
// 	// 	Email:     &Email,
// 	// 	Password:  Password,
// 	// 	IsAdmin:   false,
// 	// 	Roles:     []*model.Role{},
// 	// }
// 	// query := `DELETE FROM "users" WHERE id = $1`

// 	s.mock.ExpectBegin()

// 	// rows := s.mock.
// 	// 	NewRows([]string{"id", "first_name", "last_name", "username", "email", "password"}).
// 	// 	AddRow(user.ID, user.FirstName, user.LastName, user.Username, user.Email, user.Password)

// 	s.mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM "users" WHERE id = $1`)).
// 		WithArgs(id)
// 	// s.mock.ExpectCommit()

// 	// s.mock.ExpectExec(query).

// 	err := s.repository.DeleteUser(id.String())
// 	// fmt.Println(err)
// 	require.NoError(s.T(), err)
// 	assert.Equal(s.T(), nil, nil)
// 	s.mock.ExpectCommit()

// }
