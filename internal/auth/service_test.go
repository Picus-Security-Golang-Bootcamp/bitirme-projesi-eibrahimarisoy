package auth

import (
	"errors"
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	user "patika-ecommerce/internal/user"
	"patika-ecommerce/pkg/config"
	jwtHelper "patika-ecommerce/pkg/jwt"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	cfg, _ := config.LoadConfig("../../pkg/config/config-local")

	type fields struct {
		cfg      *config.Config
		userRepo user.UserRepositoryForMock
	}
	type args struct {
		user *model.User
	}

	firstname, lastname, email, username, password := "test", "test", "emre@arisoy.com", "emre", "123456Aa"
	wrongEmail := "wrong@email.com"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockRepo := &mockUserRepository{
		items: []model.User{
			{
				Base:      model.Base{ID: uuid.New()},
				FirstName: &firstname,
				LastName:  &lastname,
				Email:     &email,
				Password:  string(hashed),
			},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.TokenResponse
		wantErr bool
	}{
		{
			name: "userLogin_Success",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname,
					LastName:  &lastname,
					Username:  &username,
					Email:     &email,
					Password:  password,
					IsAdmin:   false,
				},
			},
			want:    api.TokenResponse{},
			wantErr: false,
		},
		{
			name: "userLogin_Failed_wrongPassword",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname,
					LastName:  &lastname,
					Username:  &username,
					Email:     &email,
					Password:  "123456",
					IsAdmin:   false,
				},
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
		{
			name: "userLogin_Failed_wrongEmail",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname,
					LastName:  &lastname,
					Username:  &username,
					Email:     &wrongEmail,
					Password:  password,
					IsAdmin:   false,
				},
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				cfg:      tt.fields.cfg,
				userRepo: tt.fields.userRepo,
			}
			_, err := a.Login(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	cfg, _ := config.LoadConfig("../../pkg/config/config-local")

	firstname, lastname, email, username, password := "test", "test", "emre@arisoy.com", "emre", "123456Aa"
	firstname2, lastname2, email2, username2, password2 := "test2", "test2", "emre2@arisoy.com", "emre2", "123456Aa"

	mockRepo := &mockUserRepository{
		items: []model.User{
			{
				Base:      model.Base{ID: uuid.New()},
				FirstName: &firstname,
				LastName:  &lastname,
				Username:  &username,
				Email:     &email,
				Password:  password,
			},
		},
	}

	type fields struct {
		cfg      *config.Config
		userRepo user.UserRepositoryForMock
	}
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.TokenResponse
		wantErr bool
	}{
		{
			name: "userRegister_Successful",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname2,
					LastName:  &lastname2,
					Username:  &username2,
					Email:     &email2,
					Password:  password2,
					IsAdmin:   false,
				},
			},
			want:    api.TokenResponse{},
			wantErr: false,
		},
		{
			name: "userRegister_Failed_duplicateEmail",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname2,
					LastName:  &lastname2,
					Username:  &username2,
					Email:     &email,
					Password:  "123456",
					IsAdmin:   true,
				},
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
		{
			name: "userRegister_Failed_duplicateUsername",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname2,
					LastName:  &lastname2,
					Username:  &username,
					Email:     &email2,
					Password:  "123456",
					IsAdmin:   true,
				},
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				cfg:      tt.fields.cfg,
				userRepo: tt.fields.userRepo,
			}
			_, err := a.Register(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	cfg, _ := config.LoadConfig("../../pkg/config/config-local")

	firstname, lastname, email, username, password := "test", "test", "emre@arisoy.com", "emre", "123456Aa"

	mockRepo := &mockUserRepository{
		items: []model.User{
			{
				Base:      model.Base{ID: uuid.New()},
				FirstName: &firstname,
				LastName:  &lastname,
				Username:  &username,
				Email:     &email,
				Password:  password,
			},
		},
	}
	invalidRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	expiredRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDcmVhdGVkQXQiOjE2NTAwMjUyMTYsIkV4cGlyZXNBdCI6MTY1MDAyODgxNiwiVXNlcklkIjoiNWQ0ODBmZTItYzY2Yy00MzkwLWJhMTctNmJmZTMxMTc5ZDY5In0.HyfE06jGkTNPyBXFMQwyUa1lGs83xsSv5N7my4kWtYY"
	anonymousUserToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDcmVhdGVkQXQiOjE2NTAwMjUyMTYsIkV4cGlyZXNBdCI6MTY1MDAyODgxNiwiVXNlcklkIjoiNWQ0ODBmZTItYzY2Yy00MzkwLWJhMTctNmJmZTMxMTc5ZDkifQ.OWIWseL90E9h3ainPmcqj6cgUvWZDoUQZ4xEzUbxTPk"

	type fields struct {
		cfg      *config.Config
		userRepo user.UserRepositoryForMock
	}
	type args struct {
		refreshToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.TokenResponse
		wantErr bool
	}{
		{
			name: "refreshToken_Successful",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				refreshToken: jwtHelper.GetAuthToken(&mockRepo.items[0], cfg).RefreshToken,
			},
			want:    api.TokenResponse{},
			wantErr: false,
		},
		{
			name: "refreshToken_Failed_InvalidRefreshToken",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				refreshToken: invalidRefreshToken,
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
		{
			name: "refreshToken_Failed_ExpiredRefreshToken",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				refreshToken: expiredRefreshToken,
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
		{
			name: "refreshToken_Failed_AnonymousUser",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				refreshToken: anonymousUserToken,
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				cfg:      tt.fields.cfg,
				userRepo: tt.fields.userRepo,
			}
			_, err := a.RefreshToken(tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////
var errCRUD = errors.New("Mock: Error crud operation")
var userNotFound = errors.New("User not found")

type mockUserRepository struct {
	items []model.User
}

// InsertUser insert user to database
func (u *mockUserRepository) InsertUser(user *model.User) (*model.User, error) {
	for _, item := range u.items {
		fmt.Println(*item.Username, *user.Username)
		fmt.Println(*item.Email, *user.Email)
		if *item.Username == *user.Username || *item.Email == *user.Email {
			fmt.Println("dasdasdasdsa")
			return nil, errCRUD
		}
	}
	u.items = append(u.items, *user)
	return user, nil
}

func (u *mockUserRepository) GetAll() (*[]model.User, error) {
	var users []model.User

	// result := u.db.Find(&users)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }
	return &users, nil
}

// GetUser
func (u *mockUserRepository) GetUser(id string) (*model.User, error) {

	if len(id) == 0 {
		return nil, errCRUD
	}
	uu, _ := uuid.Parse(id)

	for _, item := range u.items {
		if item.ID == uu {
			return &item, nil
		}
	}
	return nil, userNotFound
}

// DeleteUser
func (u *mockUserRepository) DeleteUser(id string) error {
	// result := u.db.Debug().Delete(&model.User{}, "id = ?", id)
	// if result.Error != nil {
	// 	return result.Error
	// }
	return nil
}

// UpdateUser
func (u *mockUserRepository) UpdateUser(user *model.User) (*model.User, error) {
	// 	result := u.db.Save(user)
	// 	if result.Error != nil {
	// 		return nil, result.Error
	// 	}
	return user, nil
}

// GetUserByEmail
func (u *mockUserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	if len(email) == 0 {
		return nil, errCRUD
	}
	fmt.Println("email: ", email)
	for _, item := range u.items {
		if *item.Email == email {
			fmt.Println("username: ", item.Password)
			return &item, nil
		}
	}
	return &user, nil
}

// // type TestSuiteEnv struct {
// // 	suite.Suite
// // 	database *gorm.DB
// // 	cfg      *config.Config
// // }

// // // Tests are run before they start
// // func (suite *TestSuiteEnv) SetupSuite() {
// // 	cfg, err := config.LoadConfig("./pkg/config/config-local")
// // 	if err != nil {
// // 		log.Fatalf("Failed to load configuration: %v", err)
// // 	}
// // 	// Connect to database

// // 	mockdb, _, err := sqlmock.New()
// // 	if err != nil {
// // 		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// // 	}
// // 	defer mockdb.Close()

// // 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// // 		Conn: mockdb,
// // 	}), &gorm.Config{})
// // 	if err != nil {
// // 		log.Fatal(err.Error())
// // 	}
// // 	suite.cfg = cfg
// // 	suite.database = gormDB
// // }

// // // Running after each test
// // func (suite *TestSuiteEnv) TearDownTest() {
// // 	// suite.database.ClearTable()
// // }

// // // Running after all tests are completed
// // func (suite *TestSuiteEnv) TearDownSuite() {
// // 	// suite.database.Close()
// // }

// // // This gets run automatically by `go test` so we call `suite.Run` inside it
// // func TestSuite(t *testing.T) {
// // 	// This is what actually runs our suite
// // 	suite.Run(t, new(TestSuiteEnv))
// // }

// // func TestAuthService_Register(t *testing.T) {
// // 	mockdb, _, err := sqlmock.New()
// // 	if err != nil {
// // 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// // 	}
// // 	defer mockdb.Close()

// // 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// // 		Conn: mockdb,
// // 	}), &gorm.Config{})
// // 	if err != nil {
// // 		t.Fatal(err.Error())
// // 	}
// // 	defer mockdb.Close()
// // 	cfg, _ := config.LoadConfig("./pkg/config/config-local")
// // 	userRepo := user.NewUserRepository(gormDB)
// // 	roleRepo := role.NewRoleRepository(gormDB)

// // 	type fields struct {
// // 		cfg      *config.Config
// // 		userRepo *user.UserRepository
// // 		roleRepo *role.RoleRepository
// // 	}
// // 	type args struct {
// // 		user *model.User
// // 	}
// // 	FirstName := "test"
// // 	LastName := "test"
// // 	Username := "test"
// // 	Email := "test@email.com"
// // 	Password := "test"

// // 	tests := []struct {
// // 		name    string
// // 		fields  fields
// // 		args    args
// // 		want    api.TokenResponse
// // 		wantErr bool
// // 	}{
// // 		{
// // 			name: "registerUser_Success",
// // 			fields: fields{
// // 				cfg:      cfg,
// // 				userRepo: userRepo,
// // 				roleRepo: roleRepo,
// // 			},
// // 			args: args{
// // 				user: &model.User{
// // 					FirstName: &FirstName,
// // 					LastName:  &LastName,
// // 					Username:  &Username,
// // 					Email:     &Email,
// // 					Password:  Password,
// // 					IsAdmin:   false,
// // 					Roles:     []*model.Role{},
// // 				},
// // 			},
// // 			wantErr: true,
// // 		},
// // 	}
// // 	for _, tt := range tests {
// // 		t.Run(tt.name, func(t *testing.T) {
// // 			a := &AuthService{
// // 				cfg:      tt.fields.cfg,
// // 				userRepo: tt.fields.userRepo,
// // 				roleRepo: tt.fields.roleRepo,
// // 			}
// // 			got, err := a.Register(tt.args.user)
// // 			if (err != nil) != tt.wantErr {
// // 				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
// // 				return
// // 			}
// // 			if !reflect.DeepEqual(got, tt.want) {
// // 				t.Errorf("AuthService.Register() = %v, want %v", got, tt.want)
// // 			}
// // 		})
// // 	}
// // }

// // func TestAuthService_Login(t *testing.T) {
// // 	mockdb, _, err := sqlmock.New()
// // 	if err != nil {
// // 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// // 	}
// // 	defer mockdb.Close()

// // 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// // 		Conn: mockdb,
// // 	}), &gorm.Config{})
// // 	if err != nil {
// // 		t.Fatal(err.Error())
// // 	}
// // 	defer mockdb.Close()
// // 	cfg, _ := config.LoadConfig("./pkg/config/config-local")
// // 	userRepo := user.NewUserRepository(gormDB)
// // 	roleRepo := role.NewRoleRepository(gormDB)

// // 	type fields struct {
// // 		cfg      *config.Config
// // 		userRepo *user.UserRepository
// // 		roleRepo *role.RoleRepository
// // 	}
// // 	type args struct {
// // 		user *model.User
// // 	}

// // 	FirstName := "test"
// // 	LastName := "test"
// // 	Username := "test"
// // 	Email := "emre.arisoy@gmail.com"
// // 	Password := "test"

// // 	newUser := model.User{
// // 		FirstName: &FirstName,
// // 		LastName:  &LastName,
// // 		Username:  &Username,
// // 		Email:     &Email,
// // 		Password:  Password,
// // 		IsAdmin:   false,
// // 		Roles:     []*model.Role{},
// // 	}

// // 	email := "emre.arisoy@gmail.com"
// // 	tests := []struct {
// // 		name    string
// // 		fields  fields
// // 		args    args
// // 		want    api.TokenResponse
// // 		wantErr bool
// // 	}{
// // 		{
// // 			name: "loginUser_Success",
// // 			fields: fields{
// // 				cfg:      cfg,
// // 				userRepo: userRepo,
// // 				roleRepo: roleRepo,
// // 			},
// // 			args: args{
// // 				user: &model.User{
// // 					Email:    &email,
// // 					Password: "test",
// // 				},
// // 			},
// // 			wantErr: true,
// // 		},
// // 	}
// // 	for _, tt := range tests {
// // 		t.Run(tt.name, func(t *testing.T) {
// // 			a := &AuthService{
// // 				cfg:      tt.fields.cfg,
// // 				userRepo: tt.fields.userRepo,
// // 				roleRepo: tt.fields.roleRepo,
// // 			}
// // 			a.Register(&newUser)

// // 			got, err := a.Login(tt.args.user)
// // 			if (err != nil) != tt.wantErr {
// // 				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
// // 				return
// // 			}
// // 			if !reflect.DeepEqual(got, tt.want) {
// // 				t.Errorf("AuthService.Login() = %v, want %v", got, tt.want)
// // 			}
// // 		})
// // 	}
// // }
