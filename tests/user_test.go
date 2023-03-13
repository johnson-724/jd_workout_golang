package tests

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"jd_workout_golang/internal/router"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	// 註冊 router
	r := gin.Default()
	router.RegisterUser(r.Group("/api"))

	// 建立 response & request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user", nil)

	r.ServeHTTP(w, req)

	// assert response
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "user", w.Body.String())
	
	if w.Code != 200 {
		t.Errorf("Response code is %v", w.Code)
	}
}
