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

	"portarius/middleware"

	inventoryDomain "portarius/internal/inventory/domain"
	inventoryRoutes "portarius/internal/inventory/routes"

	packageDomain "portarius/internal/package/domain"
	packageRoutes "portarius/internal/package/routes"

	reservationDomain "portarius/internal/reservation/domain"
	reservationRoutes "portarius/internal/reservation/routes"

	residentDomain "portarius/internal/resident/domain"
	residentRoutes "portarius/internal/resident/routes"

	userDomain "portarius/internal/user/domain"
	userRoutes "portarius/internal/user/routes"
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
		&reservationDomain.Reservation{},
		&userDomain.User{},
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

	userRoutes.RegisterUserRoutes(apiPrefixGroup, db)

	apiPrefixGroup.Use(middleware.AuthMiddleware())
	{
		inventoryRoutes.RegisterInventoryRoutes(apiPrefixGroup, db)
		residentRoutes.ResidentRegisterRoutes(apiPrefixGroup, db)
		packageRoutes.RegisterPackageRoutes(apiPrefixGroup, db)
		reservationRoutes.RegisterReservationRoutes(apiPrefixGroup, db)
		userRoutes.RegisterUserProtectedRoutes(apiPrefixGroup, db)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
