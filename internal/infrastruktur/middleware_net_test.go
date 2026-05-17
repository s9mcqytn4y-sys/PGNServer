package infrastruktur

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareCORS(t *testing.T) {
	// Setup environment
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://example.com")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(MiddlewareCORS())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test 1: Origin not in allowed list
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://malicious.com")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code)

	// Test 2: Origin allowed
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("Origin", "http://localhost:3000")
	resp2 := httptest.NewRecorder()
	r.ServeHTTP(resp2, req2)
	assert.Equal(t, http.StatusOK, resp2.Code)
	assert.Equal(t, "http://localhost:3000", resp2.Header().Get("Access-Control-Allow-Origin"))

	// Test 3: Preflight request (OPTIONS) allowed
	req3, _ := http.NewRequest("OPTIONS", "/test", nil)
	req3.Header.Set("Origin", "http://example.com")
	req3.Header.Set("Access-Control-Request-Method", "GET")
	resp3 := httptest.NewRecorder()
	r.ServeHTTP(resp3, req3)
	assert.Equal(t, http.StatusNoContent, resp3.Code)
	assert.Equal(t, "http://example.com", resp3.Header().Get("Access-Control-Allow-Origin"))
}

func TestMiddlewareSecureHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(MiddlewareSecureHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, "DENY", resp.Header().Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", resp.Header().Get("X-Content-Type-Options"))
}

func TestMiddlewareIPWhitelist(t *testing.T) {
	os.Setenv("IP_WHITELIST", "127.0.0.1,192.168.1.1")
	defer os.Unsetenv("IP_WHITELIST")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(MiddlewareIPWhitelist())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test 1: Allowed IP
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test 2: Blocked IP
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("X-Forwarded-For", "10.0.0.1")
	resp2 := httptest.NewRecorder()
	r.ServeHTTP(resp2, req2)
	assert.Equal(t, http.StatusForbidden, resp2.Code)
}
