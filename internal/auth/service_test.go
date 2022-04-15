// package auth

// import (
// 	"errors"
// 	"patika-ecommerce/internal/api"
// 	"patika-ecommerce/internal/model"
// 	"patika-ecommerce/internal/role"
// 	user "patika-ecommerce/internal/user"
// 	"patika-ecommerce/pkg/config"
// 	"reflect"
// 	"testing"

// 	"github.com/google/uuid"
// )

// func TestAuthService_Login(t *testing.T) {
// 	// mockdb, _, err := sqlmock.New()
// 	// if err != nil {
// 	// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	// }
// 	// defer mockdb.Close()

// 	// gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 	// 	Conn: mockdb,
// 	// }), &gorm.Config{})
// 	// if err != nil {
// 	// 	t.Fatal(err.Error())
// 	// }
// 	// defer mockdb.Close()
// 	cfg, _ := config.LoadConfig("./pkg/config/config-local")

// 	type fields struct {
// 		cfg      *config.Config
// 		userRepo *user.UserRepository
// 		roleRepo *role.RoleRepository
// 	}
// 	type args struct {
// 		user *model.User
// 	}

// 	firstname, lastname, Email, password := "test", "test", "emre@arisoy.com", "123456"
// 	mockRepo := &mockUserRepository{
// 		items: []model.User{
// 			{
// 				Base:      model.Base{ID: uuid.New()},
// 				FirstName: &firstname,
// 				LastName:  &lastname,
// 				Email:     &Email,
// 				Password:  password,
// 			},
// 		},
// 	}

// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    api.TokenResponse
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Login_Success",
// 			fields: fields{
// 				cfg:      cfg,
// 				userRepo: mockRepo,
// 				roleRepo: &mockRoleRepository{
// 					items: []model.Role{
// 						{
// 							ID:   uuid.New().String(),
// 							Name: "admin",
// 						},
// 					},
// 				},
// 			},
// 			args: args{
// 				user: &model.User{
// 					ID:        uuid.New().String(),
// 					FirstName: "John",
// 					LastName:  "Doe",
// 					Email:     "eibrahimarisoy@arisoy.com",
// 					Password:  "password",
// 				},
// 			},
// 			want:    api.TokenResponse{},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			a := &AuthService{
// 				cfg:      tt.fields.cfg,
// 				userRepo: tt.fields.userRepo,
// 				roleRepo: tt.fields.roleRepo,
// 			}
// 			got, err := a.Login(tt.args.user)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("AuthService.Login() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// /////////////////////////////////////////////////////////////////////////
// var errCRUD = errors.New("Mock: Error crud operation")
// var userNotFound = errors.New("User not found")

// type mockUserRepository struct {
// 	items []model.User
// }

// // InsertUser insert user to database
// func (u *mockUserRepository) InsertUser(user *model.User) (*model.User, error) {

// 	result := u.db.Create(user)

// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return user, nil

// }

// func (u *mockUserRepository) GetAll() (*[]model.User, error) {
// 	var users []model.User

// 	result := u.db.Find(&users)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &users, nil
// }

// // GetUser
// func (u *mockUserRepository) GetUser(id string) (*model.User, error) {

// 	if len(id) == 0 {
// 		return nil, errCRUD
// 	}
// 	uu, _ := uuid.Parse(id)
// 	for _, item := range u.items {
// 		if item.ID == uu {
// 			return &item, nil
// 		}
// 	}
// 	return nil, userNotFound
// }

// // DeleteUser
// func (u *mockUserRepository) DeleteUser(id string) error {
// 	result := u.db.Debug().Delete(&model.User{}, "id = ?", id)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// // UpdateUser
// func (u *mockUserRepository) UpdateUser(user *model.User) (*model.User, error) {
// 	result := u.db.Save(user)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return user, nil
// }

// // GetUserByEmail
// func (u *mockUserRepository) GetUserByEmail(email string) (*model.User, error) {
// 	var user model.User

// 	result := u.db.First(&user, "email = ?", email)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &user, nil
// }

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
