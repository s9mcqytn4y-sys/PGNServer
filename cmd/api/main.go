package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"pgn-server/internal/konfigurasi"
	"pgn-server/internal/pengelola"
	"pgn-server/pkg/utilitas"
	_ "pgn-server/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

// @title PGNServer API
// @version 1.0
// @description API Server untuk PGN (Produk Gagal & NG) Intelligence System.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	konfigurasi.HubungkanDatabase()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(pengelola.LogSlogMiddleware())
	r.Use(pengelola.TimeoutMiddleware(60 * time.Second))
	r.Use(pengelola.LimitRequestMiddleware(rate.Limit(200), 400))

	r.Static("/penyimpanan", "./penyimpanan")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/kesehatan", func(c *gin.Context) {
			c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Server PGN berjalan dengan baik", gin.H{
				"waktu_server": time.Now().Format(time.RFC3339),
				"status":       "stabil",
			}))
		})

		// 1. Auth (Login)
		v1.POST("/auth/login", pengelola.LoginHandler)

		// 2. QC Checksheet & QC Tools (Terproteksi)
		terproteksi := v1.Group("/")
		terproteksi.Use(pengelola.AutentikasiMiddleware())
		{
			// QC (Checksheet)
			qc := terproteksi.Group("/qc")
			{
				qc.POST("/checksheet", pengelola.TambahChecksheet)
				qc.GET("/form-checksheet/:produk_id", pengelola.GetFormChecksheet)
			}

			// QC Tools (Analytics)
			tools := terproteksi.Group("/qc-tools")
			{
				tools.GET("/pareto", pengelola.GetParetoData)
				tools.GET("/control-chart", pengelola.GetControlChartData)
			}

			// Media
			media := terproteksi.Group("/media")
			{
				media.POST("/upload", pengelola.UploadHandler)
			}
		}
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server PGNServer dimulai", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("Gagal menjalankan server", "error", err)
		os.Exit(1)
	}
}
