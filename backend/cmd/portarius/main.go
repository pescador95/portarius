package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	inventoryDomain "portarius/internal/inventory/domain"
	inventoryHandler "portarius/internal/inventory/handler"

	packageDomain "portarius/internal/package/domain"
	packageHandler "portarius/internal/package/handler"

	"portarius/internal/reservation"

	residentDomain "portarius/internal/resident/domain"
	residentHandler "portarius/internal/resident/handler"

	"portarius/middleware"

	userHandler "portarius/internal/user"
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
		&inventoryDomain.Inventory{},
		&packageDomain.Package{},
		&residentDomain.Resident{},
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

	userHandler.RegisterUserRoutes(apiPrefixGroup, db)

	apiPrefixGroup.Use(middleware.AuthMiddleware())
	{
		inventoryHandler.RegisterInventoryRoutes(apiPrefixGroup, db)
		residentHandler.ResidentRegisterRoutes(apiPrefixGroup, db)
		packageHandler.RegisterPackageRoutes(apiPrefixGroup, db)
		reservation.RegisterReservationRoutes(apiPrefixGroup, db)
		userHandler.RegisterUserProtectedRoutes(apiPrefixGroup, db)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
