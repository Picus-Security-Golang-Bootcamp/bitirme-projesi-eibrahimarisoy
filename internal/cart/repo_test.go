package cart

import (
	"database/sql"
	"patika-ecommerce/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
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
	id     = uuid.New()
	userId = uuid.New()
	name   = "test"
)

var u = model.User{
	Base: model.Base{ID: userId},
}

var c = model.Cart{
	Base:   model.Base{ID: id},
	Status: "created",
	UserID: userId,
	User:   model.User{},
	Items:  []model.CartItem{},
}

// func TestCartRepository_GetCreatedCart(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &CartRepository{db}

// 	query := `SELECT * FROM "carts" WHERE user_id = $1 AND status = $2 ORDER BY "carts"."id" LIMIT 1`

// 	rows := sqlmock.NewRows([]string{"id", "status", "user_id"}).
// 		AddRow(c.ID, c.Status, userId)

// 	str := fmt.Sprintf("%s", c.UserID)
// 	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(str, model.CartStatusCreated).WillReturnRows(rows)

// 	cart, _ := repo.GetCreatedCart(&u)
// 	fmt.Println(cart)
// 	// assert.Equal(t, cart.ID, id)

// }
