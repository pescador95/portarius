package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"portarius/internal/inventory"
	pkg "portarius/internal/package"
	"portarius/internal/reservation"
	"portarius/internal/resident/domain"
	residentHandler "portarius/internal/resident/handler"
	"portarius/middleware"

	"portarius/internal/user"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(
		&inventory.Inventory{},
		&pkg.Package{},
		&domain.Resident{},
		&reservation.Reservation{},
		&user.User{},
	)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	apiPrefixGroup := r.Group("/api")

	user.RegisterRoutes(apiPrefixGroup, db)

	apiPrefixGroup.Use(middleware.AuthMiddleware())
	{
		inventory.RegisterRoutes(apiPrefixGroup, db)
		pkg.RegisterRoutes(apiPrefixGroup, db)
		residentHandler.RegisterRoutes(apiPrefixGroup, db)
		reservation.RegisterRoutes(apiPrefixGroup, db)
		user.RegisterProtectedRoutes(apiPrefixGroup, db)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
