package repository

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"es-app/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

type UserRepositorySuite struct {
	suite.Suite
	c          echo.Context
	repository *userRepository
	db         *gorm.DB
	mock       sqlmock.Sqlmock
}

func (s *UserRepositorySuite) SetupSuite() {
	var err error
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	s.c = e.NewContext(req, rec)

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
	sqlDB, err := s.db.DB()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when getting sql.DB from gorm.DB", err)
	}
	sqlDB.Close()
}

func (s *UserRepositorySuite) TestGetUser() {
	const timeFormat = "2006-01-02 15:04:05.999999-07"
	defaultTimeStr := "2023-01-01 00:00:00.000000-00"
	defaultTime, err := time.Parse(timeFormat, defaultTimeStr)
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when parsing default time", err)
	}

	// テストデータのセットアップ
	user := &model.User{
		UserID:         "87643a78-4041-70f8-c35c-2968abba0bd9",
		Username:       "test_user",
		Email:          "test@test.com",
		WorkExperience: "test_work_experience",
		Skills:         "test_skills",
		SelfPR:         "test_self_pr",
		FutureGoals:    "test_future_goals",
		CreatedAt:      defaultTime,
		UpdatedAt:      defaultTime,
	}

	// モックのセットアップ
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE user_id = \\$1 ORDER BY \"users\".\"user_id\" LIMIT \\$2").
		WithArgs("87643a78-4041-70f8-c35c-2968abba0bd9", 1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "email", "work_experience", "skills", "self_pr", "future_goals", "created_at", "updated_at"}).
			AddRow(user.UserID, user.Username, user.Email, user.WorkExperience, user.Skills, user.SelfPR, user.FutureGoals, user.CreatedAt, user.UpdatedAt))

	// テスト対象のメソッドを実行
	result, err := s.repository.GetUser(s.c, "87643a78-4041-70f8-c35c-2968abba0bd9")
	s.NoError(err)
	s.Equal(user.UserID, result.UserID)
	s.Equal(user.Username, result.Username)
	s.Equal(user.Email, result.Email)
	s.Equal(user.WorkExperience, result.WorkExperience)
	s.Equal(user.Skills, result.Skills)
	s.Equal(user.SelfPR, result.SelfPR)
	s.Equal(user.FutureGoals, result.FutureGoals)
	s.Equal(user.CreatedAt, result.CreatedAt)
	s.Equal(user.UpdatedAt, result.UpdatedAt)
}
