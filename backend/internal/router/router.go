package router

import (
	"github.com/ai-fitness-planner/backend/internal/config"
	"github.com/ai-fitness-planner/backend/internal/handler"
	"github.com/ai-fitness-planner/backend/internal/middleware"
	"github.com/ai-fitness-planner/backend/internal/pkg/jwt"
	"github.com/ai-fitness-planner/backend/internal/pkg/session"
	"github.com/ai-fitness-planner/backend/internal/repository"
	"github.com/ai-fitness-planner/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Dependencies holds all dependencies needed for router setup
type Dependencies struct {
	DB             *gorm.DB
	RedisClient    *redis.Client
	JWTManager     jwt.JWTManager
	SessionManager session.SessionManager
	RateLimiter    *middleware.RateLimiter

	// Services
	AuthService       service.AuthService
	UserService       service.UserService
	AIAPIService      service.AIAPIService
	TrainingService   service.TrainingService
	NutritionService  service.NutritionService
	StatisticsService service.StatisticsService

	// Repositories
	AssessmentRepo repository.AssessmentRepository
}

// SetupRouter configures and returns the Gin router with all routes and middleware
func SetupRouter(deps *Dependencies) *gin.Engine {
	// Set Gin mode based on configuration
	if config.GlobalConfig.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware stack (order matters!)
	// 1. Recovery - catch panics first
	router.Use(middleware.RecoveryMiddleware(nil))

	// 2. Logging - log all requests
	router.Use(middleware.LoggingMiddleware(nil))

	// 3. CORS - handle cross-origin requests
	corsConfig := middleware.DefaultCORSConfig()
	if config.GlobalConfig.App.Mode == "release" {
		// In production, specify allowed origins
		// corsConfig = middleware.ProductionCORSConfig([]string{"https://yourdomain.com"})
	}
	router.Use(middleware.CORSMiddleware(corsConfig))

	// 4. Security - input sanitization and security headers
	router.Use(middleware.SecurityMiddleware(nil))

	// Health check endpoint (no authentication required)
	healthHandler := handler.NewHealthHandler()
	router.GET("/health", healthHandler.HealthCheck)

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		setupPublicRoutes(v1, deps)

		// Protected routes (authentication required)
		setupProtectedRoutes(v1, deps)
	}

	return router
}

// setupPublicRoutes configures public API routes (no authentication)
func setupPublicRoutes(rg *gin.RouterGroup, deps *Dependencies) {
	authHandler := handler.NewAuthHandler(deps.AuthService)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}
}

// setupProtectedRoutes configures protected API routes (authentication required)
func setupProtectedRoutes(rg *gin.RouterGroup, deps *Dependencies) {
	// Create protected group with authentication and rate limiting
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware(deps.JWTManager, deps.SessionManager))
	protected.Use(deps.RateLimiter.RateLimitMiddleware())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(deps.AuthService)
	userHandler := handler.NewUserHandler(deps.UserService)
	aiAPIHandler := handler.NewAIAPIHandler(deps.AIAPIService)
	assessmentHandler := handler.NewAssessmentHandler(deps.AssessmentRepo)
	trainingHandler := handler.NewTrainingHandler(deps.TrainingService)
	nutritionHandler := handler.NewNutritionHandler(deps.NutritionService)
	statisticsHandler := handler.NewStatisticsHandler(deps.StatisticsService)

	// Auth routes (logout requires authentication)
	{
		protected.POST("/auth/logout", authHandler.Logout)
	}

	// User routes
	user := protected.Group("/user")
	{
		user.GET("/profile", userHandler.GetProfile)
		user.PUT("/profile", userHandler.UpdateProfile)
		user.POST("/body-data", userHandler.AddBodyData)
		user.GET("/body-data", userHandler.GetBodyDataHistory)
		user.POST("/fitness-goals", userHandler.SetFitnessGoals)
		user.GET("/fitness-goals", userHandler.GetFitnessGoals)
		user.PUT("/fitness-goals", userHandler.UpdateFitnessGoals)
	}

	// AI API management routes
	aiAPIs := protected.Group("/ai-apis")
	{
		aiAPIs.POST("", aiAPIHandler.AddAPI)
		aiAPIs.GET("", aiAPIHandler.ListAPIs)
		aiAPIs.GET("/:id", aiAPIHandler.GetAPI)
		aiAPIs.PUT("/:id", aiAPIHandler.UpdateAPI)
		aiAPIs.DELETE("/:id", aiAPIHandler.DeleteAPI)
		aiAPIs.POST("/:id/test", aiAPIHandler.TestAPI)
		aiAPIs.POST("/:id/set-default", aiAPIHandler.SetDefault)
	}

	// Assessment routes
	assessments := protected.Group("/assessments")
	{
		assessments.POST("", assessmentHandler.CreateAssessment)
		assessments.GET("/latest", assessmentHandler.GetLatestAssessment)
	}

	// Training plan routes (with stricter rate limiting for generation)
	trainingPlans := protected.Group("/training-plans")
	{
		// AI generation endpoint with stricter rate limit
		generation := trainingPlans.Group("")
		generation.Use(deps.RateLimiter.AIGenerationRateLimitMiddleware())
		generation.POST("/generate", trainingHandler.GeneratePlan)

		// Regular endpoints
		trainingPlans.GET("/tasks/:taskId", trainingHandler.GetPlanStatus)
		trainingPlans.GET("", trainingHandler.ListPlans)
		trainingPlans.GET("/:id", trainingHandler.GetPlanDetail)
		trainingPlans.GET("/today", trainingHandler.GetTodayTraining)
	}

	// Training record routes
	trainingRecords := protected.Group("/training-records")
	{
		trainingRecords.POST("", trainingHandler.RecordTraining)
		trainingRecords.GET("", trainingHandler.ListTrainingRecords)
	}

	// Nutrition plan routes (with stricter rate limiting for generation)
	nutritionPlans := protected.Group("/nutrition-plans")
	{
		// AI generation endpoint with stricter rate limit
		generation := nutritionPlans.Group("")
		generation.Use(deps.RateLimiter.AIGenerationRateLimitMiddleware())
		generation.POST("/generate", nutritionHandler.GeneratePlan)
		nutritionPlans.GET("/tasks/:taskId", nutritionHandler.GetPlanStatus)

		// Regular endpoints
		nutritionPlans.GET("", nutritionHandler.ListPlans)
		nutritionPlans.GET("/:id", nutritionHandler.GetPlanDetail)
		nutritionPlans.GET("/today", nutritionHandler.GetTodayMeals)
	}

	// Nutrition record routes
	nutritionRecords := protected.Group("/nutrition-records")
	{
		nutritionRecords.POST("", nutritionHandler.RecordMeal)
		nutritionRecords.GET("", nutritionHandler.ListNutritionRecords)
		nutritionRecords.GET("/daily-summary", nutritionHandler.GetDailySummary)
	}

	// Statistics routes
	stats := protected.Group("/stats")
	{
		stats.GET("/training", statisticsHandler.GetTrainingStatistics)
		stats.GET("/progress", statisticsHandler.GetProgressReport)
		stats.GET("/trends", statisticsHandler.GetTrends)
	}
}
