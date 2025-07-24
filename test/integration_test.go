package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-clean-gin/config"
	"go-clean-gin/internal/container"
	"go-clean-gin/internal/entity"
	"go-clean-gin/internal/router"
	"go-clean-gin/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	app *gin.Engine
	db  *gorm.DB
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Setup test database
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Name:     "go_clean_gin_test",
			SSLMode:  "disable",
		},
		JWT: config.JWTConfig{
			Secret:          "test-secret",
			ExpirationHours: 24,
		},
	}

	db, err := database.NewPostgresDB(&cfg.Database)
	suite.Require().NoError(err)

	err = database.RunMigrations(db)
	suite.Require().NoError(err)

	suite.db = db

	// Setup app
	container := container.NewContainer(cfg, db)
	suite.app = router.SetupRouter(container)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *IntegrationTestSuite) TestAuthFlow() {
	// Register user
	registerReq := entity.RegisterRequest{
		Email:     "integration@test.com",
		Username:  "integration",
		Password:  "password123",
		FirstName: "Integration",
		LastName:  "Test",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.app.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Login user
	loginReq := entity.LoginRequest{
		Email:    "integration@test.com",
		Password: "password123",
	}

	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	suite.app.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(suite.T(), err)

	data := loginResp["data"].(map[string]interface{})
	token := data["token"].(string)
	assert.NotEmpty(suite.T(), token)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
