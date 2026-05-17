package otentikasi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"pgn-server/internal/infrastruktur"
	"pgn-server/internal/otentikasi"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MockRepositoriOtentikasi adalah tiruan lokal untuk menghindari import cycle
type MockRepositoriOtentikasi struct {
	mock.Mock
}

func (m *MockRepositoriOtentikasi) Daftar(pengguna *otentikasi.Pengguna) error {
	args := m.Called(pengguna)
	return args.Error(0)
}

func (m *MockRepositoriOtentikasi) CariBerdasarkanSurel(surel string) (*otentikasi.Pengguna, error) {
	args := m.Called(surel)
	if args.Get(0) != nil {
		return args.Get(0).(*otentikasi.Pengguna), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriOtentikasi) CariBerdasarkanID(id uint) (*otentikasi.Pengguna, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*otentikasi.Pengguna), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriOtentikasi) PerbaruiSandi(pengguna *otentikasi.Pengguna) error {
	args := m.Called(pengguna)
	return args.Error(0)
}

func (m *MockRepositoriOtentikasi) PerbaruiTenggatSesi(id uint, tenggatSesi *gorm.DeletedAt) error {
	args := m.Called(id, tenggatSesi)
	return args.Error(0)
}

func TestIntegrasiAlurHTTPDanJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set kunci rahasia JWT untuk pengetesan
	jwtSecret := "kunci-rahasia-integrasi-test-12345"
	os.Setenv("JWT_SECRET", jwtSecret)
	defer os.Unsetenv("JWT_SECRET")

	// 1. Persiapkan Mock Repository
	mockRepo := new(MockRepositoriOtentikasi)
	layanan := otentikasi.KonstruksiLayananBaru(mockRepo)
	handler := otentikasi.KonstruksiPenangananBaru(layanan)

	// 2. Persiapkan Gin Engine dan Daftarkan Router
	router := gin.New()
	router.Use(gin.Recovery())

	// Grup otentikasi publik
	authGrup := router.Group("/api/v1/otentikasi")
	{
		authGrup.POST("/daftar", handler.TanganiRegistrasi)
		authGrup.POST("/masuk", handler.TanganiLogin)
		authGrup.POST("/lupa-sandi", handler.TanganiLupaSandi)
		authGrup.POST("/keluar", handler.TanganiLogout)
		// Proteksi JWT ditambahkan khusus untuk endpoint profil
		authGrup.GET("/profil", infrastruktur.PenjagaSesiJWT(), handler.TanganiProfil)
	}

	// Password hashing untuk dummy user login
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	dummyUser := &otentikasi.Pengguna{
		ID:                    42,
		SurelKredensial:       "tester.qc@pgn-quality.co.id",
		KataSandiTerenskripsi: string(hashedPassword),
		PeranOtorisasi:        "LEADER",
		DibuatPada:            time.Now(),
		DiubahPada:            time.Now(),
	}

	t.Run("1. Registrasi Akun - Validasi Payload & Retro-kompatibilitas NIP", func(t *testing.T) {
		mockRepo.On("CariBerdasarkanSurel", "2211019@pgn-quality.co.id").Return(nil, nil).Once()
		mockRepo.On("Daftar", mock.AnythingOfType("*otentikasi.Pengguna")).Return(nil).Once()

		reqBody := map[string]string{
			"nip":        "2211019",
			"kata_sandi": "admin123",
			"peran":      "LEADER",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/otentikasi/daftar", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.True(t, resp["sukses"].(bool))
		assert.Equal(t, "Registrasi akun berhasil", resp["pesan"])
		assert.NotNil(t, resp["data"])

		dataMap := resp["data"].(map[string]interface{})
		assert.Equal(t, "2211019@pgn-quality.co.id", dataMap["surelKredensial"])
		assert.Equal(t, "LEADER", dataMap["peranOtorisasi"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("2. Login Sukses & JWT Token Generation", func(t *testing.T) {
		mockRepo.On("CariBerdasarkanSurel", "tester.qc@pgn-quality.co.id").Return(dummyUser, nil).Once()

		reqBody := map[string]string{
			"surel": "tester.qc@pgn-quality.co.id",
			"sandi": "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/otentikasi/masuk", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.True(t, resp["sukses"].(bool))
		assert.Equal(t, "Autentikasi berhasil, token diterbitkan", resp["pesan"])

		dataMap := resp["data"].(map[string]interface{})
		tokenStr := dataMap["token"].(string)
		assert.NotEmpty(t, tokenStr)

		// Verifikasi secara internal signature token
		parsedToken, errParse := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		assert.NoError(t, errParse)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, float64(42), claims["id"])
		assert.Equal(t, "tester.qc@pgn-quality.co.id", claims["surel"])
		assert.Equal(t, "LEADER", claims["peran"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("3. Profil Endpoint dengan JWT Valid (Retro-kompatibilitas Double Envelope)", func(t *testing.T) {
		mockRepo.On("CariBerdasarkanID", uint(42)).Return(dummyUser, nil).Once()

		// Generate token valid
		claims := jwt.MapClaims{
			"id":    dummyUser.ID,
			"surel": dummyUser.SurelKredensial,
			"peran": dummyUser.PeranOtorisasi,
			"exp":   time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, _ := token.SignedString([]byte(jwtSecret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/otentikasi/profil", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Validasi double envelope untuk KMP Client compatibility
		assert.True(t, resp["sukses"].(bool))
		assert.True(t, resp["success"].(bool)) // Legacy
		assert.Equal(t, "Profil berhasil diambil", resp["pesan"])
		assert.Equal(t, "Profil berhasil diambil", resp["message"]) // Legacy

		dataMap := resp["data"].(map[string]interface{})
		assert.Equal(t, float64(42), dataMap["id"])
		assert.Equal(t, "tester.qc@pgn-quality.co.id", dataMap["surel"])
		assert.Equal(t, "tester.qc@pgn-quality.co.id", dataMap["email"]) // Legacy
		assert.Equal(t, "LEADER", dataMap["peran"])
		assert.Equal(t, "LEADER", dataMap["role"]) // Legacy
		assert.Equal(t, "Leader QC", dataMap["nama"])
		assert.Equal(t, "Leader QC", dataMap["fullName"]) // Legacy
		assert.Equal(t, "tester.qc", dataMap["nip"])      // NIP parsed from email prefix

		mockRepo.AssertExpectations(t)
	})

	t.Run("4. Profil Endpoint Tanpa Token JWT - Menolak Akses", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/otentikasi/profil", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.False(t, resp["sukses"].(bool))
		assert.False(t, resp["success"].(bool))
		assert.Equal(t, "Akses ditolak: Tajuk otorisasi tidak ditemukan", resp["pesan"])
	})

	t.Run("5. Profil Endpoint dengan Token JWT Kadaluwarsa/Tidak Sah", func(t *testing.T) {
		// Token ditandatangani dengan kunci rahasia yang salah
		badClaims := jwt.MapClaims{
			"id":  uint(42),
			"exp": time.Now().Add(-time.Hour).Unix(), // Kadaluwarsa
		}
		badToken := jwt.NewWithClaims(jwt.SigningMethodHS256, badClaims)
		badTokenStr, _ := badToken.SignedString([]byte("kunci-rahasia-salah"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/otentikasi/profil", nil)
		req.Header.Set("Authorization", "Bearer "+badTokenStr)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.False(t, resp["sukses"].(bool))
		assert.Equal(t, "Akses ditolak: Sesi tidak sah atau telah kedaluwarsa", resp["pesan"])
	})
}
