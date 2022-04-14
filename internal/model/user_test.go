package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUser_CheckPassword(t *testing.T) {
	type fields struct {
		Base      Base
		FirstName *string
		LastName  *string
		Username  *string
		Email     *string
		Password  string
		IsAdmin   bool
		Roles     []*Role
	}
	type args struct {
		password string
	}

	FirstName := "test"
	LastName := "test"
	Username := "test"
	Email := "test@email.com"
	Password := "test"
	password, _ := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"passwordMatchedSuccessfully",
			fields{
				Base:      Base{ID: "123"},
				FirstName: &FirstName,
				LastName:  &LastName,
				Username:  &Username,
				Email:     &Email,
				Password:  string(password),
				IsAdmin:   false,
			},
			args{password: "test"},
			true,
		},
		{
			"passwordMatchedFailed",
			fields{
				Base:      Base{ID: "123"},
				FirstName: &FirstName,
				LastName:  &LastName,
				Username:  &Username,
				Email:     &Email,
				Password:  string(password),
				IsAdmin:   false,
			},
			args{password: "test1"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Base:      tt.fields.Base,
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
				Username:  tt.fields.Username,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				IsAdmin:   tt.fields.IsAdmin,
				Roles:     tt.fields.Roles,
			}
			if got := u.CheckPassword(tt.args.password); got != tt.want {
				t.Errorf("User.CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	mockdb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockdb.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockdb,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer mockdb.Close()

	type fields struct {
		Base      Base
		FirstName *string
		LastName  *string
		Username  *string
		Email     *string
		Password  string
		IsAdmin   bool
		Roles     []*Role
	}
	type args struct {
		tx *gorm.DB
	}

	FirstName := "test"
	LastName := "test"
	Username := "test"
	Email := "test@email.com"
	Password := "test"
	// password, _ := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"passwordGeneratedSuccessfully",
			fields{
				Base:      Base{ID: "123"},
				FirstName: &FirstName,
				LastName:  &LastName,
				Username:  &Username,
				Email:     &Email,
				Password:  Password,
				IsAdmin:   false,
			},
			args{tx: gormDB},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Base:      tt.fields.Base,
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
				Username:  tt.fields.Username,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				IsAdmin:   tt.fields.IsAdmin,
			}
			if err := u.BeforeCreate(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("User.BeforeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
