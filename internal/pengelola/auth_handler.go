package pengelola

import (
	"net/http"
	"os"
	"pgn-server/internal/konfigurasi"
	"pgn-server/internal/model"
	"pgn-server/pkg/utilitas"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	NIP      string `json:"nip" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler menangani otentikasi user
// @Summary User Login
// @Description Login menggunakan NIP dan Password untuk mendapatkan JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginInput true "Credentials"
// @Success 200 {object} utilitas.ResponsAPI
// @Failure 401 {object} utilitas.ResponsAPI
// @Router /auth/login [post]
func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utilitas.FormatRespons(false, "Input tidak valid", nil))
		return
	}

	var user model.User
	if err := konfigurasi.DB.Where("nip = ?", input.NIP).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, utilitas.FormatRespons(false, "NIP atau Password salah", nil))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, utilitas.FormatRespons(false, "NIP atau Password salah", nil))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"nip":     user.NIP,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "pgn-secret-key-2026"
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, "Gagal membuat token", nil))
		return
	}

	c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Login berhasil", gin.H{
		"token": tokenString,
		"user": gin.H{
			"nama": user.Nama,
			"role": user.Role,
		},
	}))
}
