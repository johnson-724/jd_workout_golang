package tests

import (
	"encoding/json"
	"jd_workout_golang/internal/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine

func init() {
	// 註冊 router
	r = gin.Default()
	router.RegisterUser(r.Group("/api"))
}

func TestUserWithoutToken(t *testing.T) {
	// 註冊 router
	// r := gin.Default()
	// router.RegisterUser(r.Group("/api"))

	// 建立 response & request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/", nil)

	r.ServeHTTP(w, req)

	// assert response
	response := struct {
		Error string `json:"error"`
	}{}

	assert.Equal(t, 403, w.Code)
	// assert.Equal(t, "JWT token is empty", w.Body.String())
	json.Unmarshal((w.Body.Bytes()), &response)
	assert.Equal(t, "JWT token is empty", response.Error)
}

func TestUserWithInvalidToken(t *testing.T) {
	testCases := []struct {
		token string
		status int
		error string
	}{
		{"test123", 403, "invalid token"},
		{"Bearer test1234", 403, "invalid token"},
		{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjN9.mX1ysOl8jt8Rdg2uYNX8B0dhnKsqfyy2UTT28_1pwZQ123", 403, "invalid token"},
	}

	// 建立 response & request
	
	for _, testCase := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user/", nil)
		req.Header.Set("Authorization", testCase.token)

		r.ServeHTTP(w, req)

		// assert response
		response := struct {
			Error string `json:"error"`
		}{}

		assert.Equal(t, testCase.status, w.Code)
		json.Unmarshal((w.Body.Bytes()), &response)
		assert.Equal(t, testCase.error, response.Error)
	}
}
