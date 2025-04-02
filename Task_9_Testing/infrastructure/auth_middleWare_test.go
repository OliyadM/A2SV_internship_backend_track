package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID, username, role string) (string, error) {
	return "", nil 
}

func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

type AuthMiddlewareTestSuite struct {
	suite.Suite
	mockJWT    *MockJWTService
	router     *gin.Engine
	authMw     gin.HandlerFunc
	adminMw    gin.HandlerFunc
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockJWT = &MockJWTService{}
	s.authMw = AuthMiddleware(s.mockJWT)
	s.adminMw = AdminMiddleware()
	s.router = gin.New()
	
	s.router.GET("/protected", s.authMw, func(c *gin.Context) {
		c.String(http.StatusOK, "Protected OK")
	})
	
	s.router.GET("/admin", s.authMw, s.adminMw, func(c *gin.Context) {
		c.String(http.StatusOK, "Admin OK")
	})
}

func (s *AuthMiddlewareTestSuite) TearDownTest() {
	s.mockJWT.AssertExpectations(s.T())
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware() {
	s.Run("Success", func() {
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"sub":  "1",
				"name": "testuser",
				"role": "user",
			},
			Valid: true,
		}
		s.mockJWT.On("ValidateToken", "valid-token").Return(token, nil).Once()

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code, "Should return 200 OK")
		s.Equal("Protected OK", w.Body.String(), "Response body should match")
	})

	s.Run("NoToken", func() {
		req, _ := http.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized")
		s.Contains(w.Body.String(), `"error":"Authorization header required"`, "Error message should match")
	})

	s.Run("InvalidToken", func() {
		s.mockJWT.On("ValidateToken", "invalid-token").Return((*jwt.Token)(nil), jwt.ErrSignatureInvalid).Once()

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized")
		s.Contains(w.Body.String(), `"error":"Invalid token"`, "Error message should match")
	})

	s.Run("InvalidClaims", func() {
		token := &jwt.Token{
			Claims: nil, 
			Valid:  true,
		}
		s.mockJWT.On("ValidateToken", "bad-claims-token").Return(token, nil).Once()

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer bad-claims-token")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized")
		s.Contains(w.Body.String(), `"error":"Invalid token claims"`, "Error message should match")
	})
}

func (s *AuthMiddlewareTestSuite) TestAdminMiddleware() {
	s.Run("SuccessAdmin", func() {
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"sub":  "1",
				"name": "adminuser",
				"role": "admin",
			},
			Valid: true,
		}
		s.mockJWT.On("ValidateToken", "admin-token").Return(token, nil).Once()

		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code, "Should return 200 OK")
		s.Equal("Admin OK", w.Body.String(), "Response body should match")
	})

	s.Run("NonAdmin", func() {
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"sub":  "1",
				"name": "testuser",
				"role": "user",
			},
			Valid: true,
		}
		s.mockJWT.On("ValidateToken", "user-token").Return(token, nil).Once()

		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer user-token")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusForbidden, w.Code, "Should return 403 Forbidden")
		s.Contains(w.Body.String(), `"error":"Unauthorized access"`, "Error message should match")
	})

	s.Run("NoRole", func() {
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"sub":  "1",
				"name": "testuser",
				
			},
			Valid: true,
		}
		s.mockJWT.On("ValidateToken", "no-role-token").Return(token, nil).Once()

		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer no-role-token")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusForbidden, w.Code, "Should return 403 Forbidden")
		s.Contains(w.Body.String(), `"error":"Unauthorized access"`, "Error message should match")
	})
}

func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}