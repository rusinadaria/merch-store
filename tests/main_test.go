package tests

import (
	"github.com/stretchr/testify/suite"
	"merch-store/internal/handlers"
	"merch-store/internal/services"
	"merch-store/internal/repository"
	"testing"
	"database/sql"
	"os"
)

type APITestSuite struct {
	suite.Suite

	db *sql.DB
	handler *handlers.Handler
	services *services.Service
	repos *repository.Repository
	token    string
	token_two_user string
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	connStr := "user=postgres password=root dbname=test_shop sslmode=disable"
	// connStr := os.Getenv("user=postgres password=root dbname=test_shop sslmode=disable")
	db, err := sql.Open("postgres", connStr)
	s.db = db
	if err != nil {
		s.FailNow("Failed connect database", err)
	}

	s.initDeps()
	// s.createTestUsers()
	// s.createTestTwoUsers()	
}

func (s *APITestSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			s.FailNow("Failed to close database connection", err)
		}
	}
}

func (s *APITestSuite) initDeps() {
	s.repos = repository.NewRepository(s.db)
	s.services = services.NewService(s.repos)
	s.handler = handlers.NewHandler(s.services)
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}