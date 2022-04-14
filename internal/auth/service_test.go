package auth

// type TestSuiteEnv struct {
// 	suite.Suite
// 	database *gorm.DB
// 	cfg      *config.Config
// }

// // Tests are run before they start
// func (suite *TestSuiteEnv) SetupSuite() {
// 	cfg, err := config.LoadConfig("./pkg/config/config-local")
// 	if err != nil {
// 		log.Fatalf("Failed to load configuration: %v", err)
// 	}
// 	// Connect to database

// 	mockdb, _, err := sqlmock.New()
// 	if err != nil {
// 		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer mockdb.Close()

// 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn: mockdb,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	suite.cfg = cfg
// 	suite.database = gormDB
// }

// // Running after each test
// func (suite *TestSuiteEnv) TearDownTest() {
// 	// suite.database.ClearTable()
// }

// // Running after all tests are completed
// func (suite *TestSuiteEnv) TearDownSuite() {
// 	// suite.database.Close()
// }

// // This gets run automatically by `go test` so we call `suite.Run` inside it
// func TestSuite(t *testing.T) {
// 	// This is what actually runs our suite
// 	suite.Run(t, new(TestSuiteEnv))
// }

// func TestAuthService_Register(t *testing.T) {
// 	mockdb, _, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer mockdb.Close()

// 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn: mockdb,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	defer mockdb.Close()
// 	cfg, _ := config.LoadConfig("./pkg/config/config-local")
// 	userRepo := user.NewUserRepository(gormDB)
// 	roleRepo := role.NewRoleRepository(gormDB)

// 	type fields struct {
// 		cfg      *config.Config
// 		userRepo *user.UserRepository
// 		roleRepo *role.RoleRepository
// 	}
// 	type args struct {
// 		user *model.User
// 	}
// 	FirstName := "test"
// 	LastName := "test"
// 	Username := "test"
// 	Email := "test@email.com"
// 	Password := "test"

// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    api.TokenResponse
// 		wantErr bool
// 	}{
// 		{
// 			name: "registerUser_Success",
// 			fields: fields{
// 				cfg:      cfg,
// 				userRepo: userRepo,
// 				roleRepo: roleRepo,
// 			},
// 			args: args{
// 				user: &model.User{
// 					FirstName: &FirstName,
// 					LastName:  &LastName,
// 					Username:  &Username,
// 					Email:     &Email,
// 					Password:  Password,
// 					IsAdmin:   false,
// 					Roles:     []*model.Role{},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			a := &AuthService{
// 				cfg:      tt.fields.cfg,
// 				userRepo: tt.fields.userRepo,
// 				roleRepo: tt.fields.roleRepo,
// 			}
// 			got, err := a.Register(tt.args.user)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("AuthService.Register() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestAuthService_Login(t *testing.T) {
// 	mockdb, _, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer mockdb.Close()

// 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn: mockdb,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	defer mockdb.Close()
// 	cfg, _ := config.LoadConfig("./pkg/config/config-local")
// 	userRepo := user.NewUserRepository(gormDB)
// 	roleRepo := role.NewRoleRepository(gormDB)

// 	type fields struct {
// 		cfg      *config.Config
// 		userRepo *user.UserRepository
// 		roleRepo *role.RoleRepository
// 	}
// 	type args struct {
// 		user *model.User
// 	}

// 	FirstName := "test"
// 	LastName := "test"
// 	Username := "test"
// 	Email := "emre.arisoy@gmail.com"
// 	Password := "test"

// 	newUser := model.User{
// 		FirstName: &FirstName,
// 		LastName:  &LastName,
// 		Username:  &Username,
// 		Email:     &Email,
// 		Password:  Password,
// 		IsAdmin:   false,
// 		Roles:     []*model.Role{},
// 	}

// 	email := "emre.arisoy@gmail.com"
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    api.TokenResponse
// 		wantErr bool
// 	}{
// 		{
// 			name: "loginUser_Success",
// 			fields: fields{
// 				cfg:      cfg,
// 				userRepo: userRepo,
// 				roleRepo: roleRepo,
// 			},
// 			args: args{
// 				user: &model.User{
// 					Email:    &email,
// 					Password: "test",
// 				},
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			a := &AuthService{
// 				cfg:      tt.fields.cfg,
// 				userRepo: tt.fields.userRepo,
// 				roleRepo: tt.fields.roleRepo,
// 			}
// 			a.Register(&newUser)

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
