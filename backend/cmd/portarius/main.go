package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"portarius/middleware"

	"portarius/internal/infra"

	reminderListeners "portarius/internal/reminder/listeners"

	reminderRepository "portarius/internal/reminder/repository"

	residentRepository "portarius/internal/resident/repository"

	inventoryRoutes "portarius/internal/inventory/routes"

	packageRoutes "portarius/internal/package/routes"

	reservationRoutes "portarius/internal/reservation/routes"

	residentRoutes "portarius/internal/resident/routes"

	userRoutes "portarius/internal/user/routes"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := infra.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	infra.RunMigrations(db)

	reminderRepo := reminderRepository.NewReminderRepository(db)
	residentRepo := residentRepository.NewResidentRepository(db)
	reminderListeners.RegisterReminderListeners(reminderRepo, residentRepo)

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
