package reopository

import (
	"es-app/model"
	"github.com/sanposhiho/gomockhandler"
	"testing"
)

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

type UserRepositorySuite struct {
	suite.Suite

	c echo.Context
	repository *userRepository
	db *gorm.DB
}

func (s *UserRepositorySuite) SetupSuite() {
	var err error
	s.c = echo.New().NewContext(nil, nil)
	
	// モックDBのセットアップ
	var db *sql.DB
	db, s.mock, err = sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when initializing gorm", err)
	}

	s.repository = &userRepository{
		db: s.db,
	}
}

func (s *UserRepositorySuite) TearDownSuite() {
	s.db.Close()
}

func (s *UserRepositorySuite) TestGetUser() {
	// テストデータのセットアップ
	user := &model.User{
		UserID: "test_user_id",
	}

	// モックのセットアップ
	s.mock.ExpectQuery("SELECT * FROM users WHERE user_id = ?").WithArgs("test_user_id").WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("test_user_id"))

	// テスト対象のメソッドを実行
	result, err := s.repository.GetUser(s.c, "test_user_id")

	// 期待値の検証
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, result)
}