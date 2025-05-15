package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

	_ "portarius/docs"

	_ "portarius/internal/inventory/handler"

	middleware "portarius/internal/middleware/auth"

	"portarius/internal/infra"

	reminderListeners "portarius/internal/reminder/listeners"
	"portarius/internal/reminder/scheduler"

	reminderRepository "portarius/internal/reminder/repository"

	residentRepository "portarius/internal/resident/repository"

	packageRepository "portarius/internal/package/repository"

	reservationRepository "portarius/internal/reservation/repository"

	inventoryRoutes "portarius/internal/inventory/routes"

	packageRoutes "portarius/internal/package/routes"

	reservationRoutes "portarius/internal/reservation/routes"

	residentRoutes "portarius/internal/resident/routes"

	userRoutes "portarius/internal/user/routes"

	reminderRoutes "portarius/internal/reminder/routes"

	whatsappDomain "portarius/internal/whatsapp/domain"
	"portarius/internal/whatsapp/handler"
)

// @title Portarius API
// @version 1.0
// @description API for managing a package delivery system in a residential building
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.portarius.com/support
// @contact.email support@portarius.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
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
	packageRepo := packageRepository.NewPackageRepository(db)
	reservationRepo := reservationRepository.NewReservationRepository(db)

	whatsappService := whatsappDomain.NewWhatsAppService()

	whatsappHandler := handler.NewWhatsAppHandler(whatsappService)

	reminderListeners.RegisterReminderListeners(reminderRepo, residentRepo, packageRepo, reservationRepo, whatsappHandler)

	reminderScheduler := scheduler.NewReminderScheduler(reminderRepo)

	reminderScheduler.Run()

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
		reminderRoutes.RegisterReminderProtectedRoutes(apiPrefixGroup, db)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + port)
}
