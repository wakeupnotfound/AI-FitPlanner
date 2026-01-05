package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ai-fitness-planner/backend/docs"
	"github.com/ai-fitness-planner/backend/internal/config"
	"github.com/ai-fitness-planner/backend/internal/middleware"
	"github.com/ai-fitness-planner/backend/internal/pkg/crypto"
	"github.com/ai-fitness-planner/backend/internal/pkg/database"
	"github.com/ai-fitness-planner/backend/internal/pkg/jwt"
	"github.com/ai-fitness-planner/backend/internal/pkg/logger"
	"github.com/ai-fitness-planner/backend/internal/pkg/redis"
	"github.com/ai-fitness-planner/backend/internal/pkg/session"
	"github.com/ai-fitness-planner/backend/internal/repository"
	"github.com/ai-fitness-planner/backend/internal/router"
	"github.com/ai-fitness-planner/backend/internal/service"
	customvalidator "github.com/ai-fitness-planner/backend/internal/validator"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// @title AI Fitness Planning System API
// @version 1.0
// @description RESTful API for AI-powered fitness and nutrition planning system. Users can configure their own AI APIs (OpenAI, Wenxin, Tongyi) to generate personalized training and nutrition plans.
// @contact.name API Support
// @contact.email support@ai-fitness-planner.com
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize configuration
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Logger.Sync()

	logger.Info("Starting AI Fitness Planner API",
		zap.String("version", config.GlobalConfig.App.Version),
		zap.String("mode", config.GlobalConfig.App.Mode),
	)

	// Register custom validators with Gin
	if err := registerCustomValidators(); err != nil {
		logger.Fatal("Failed to register custom validators", zap.Error(err))
	}
	logger.Info("Custom validators registered")

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.Close()
	logger.Info("Database connection established")

	// Initialize Redis
	if err := redis.InitRedis(); err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	defer redis.Close()
	logger.Info("Redis connection established")

	// Setup dependencies
	deps, err := setupDependencies()
	if err != nil {
		logger.Fatal("Failed to setup dependencies", zap.Error(err))
	}

	// Initialize router with dependencies
	ginRouter := router.SetupRouter(deps)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.GlobalConfig.App.Port),
		Handler:      ginRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.Int("port", config.GlobalConfig.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// setupDependencies initializes all dependencies for dependency injection
func setupDependencies() (*router.Dependencies, error) {
	db := database.GetDB()
	redisClient := redis.Rdb

	// Initialize utilities
	encryptor, err := crypto.NewEncryptor(config.GlobalConfig.App.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryptor: %w", err)
	}

	jwtManager := jwt.NewJWTManager(
		config.GlobalConfig.JWT.Secret,
		config.GlobalConfig.JWT.AccessTokenExpire,
		config.GlobalConfig.JWT.RefreshTokenExpire,
	)
	sessionManager := session.NewSessionManager(redisClient)

	// Initialize rate limiter
	rateLimitConfig := &middleware.RateLimitConfig{
		UserRequestsPerMinute: config.GlobalConfig.RateLimit.APICallsPerMinute,
		UserRequestsPerHour:   config.GlobalConfig.RateLimit.APICallsPerHour,
		IPRequestsPerMinute:   100,
		AIGenerationPerMinute: 2,
	}
	rateLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfig)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	aiAPIRepo := repository.NewAIAPIRepository(db)
	trainingPlanRepo := repository.NewTrainingPlanRepository(db)
	trainingRecordRepo := repository.NewTrainingRecordRepository(db)
	nutritionPlanRepo := repository.NewNutritionPlanRepository(db)
	nutritionRecordRepo := repository.NewNutritionRecordRepository(db)
	assessmentRepo := repository.NewAssessmentRepository(db)
	bodyDataRepo := repository.NewBodyDataRepository(db)
	fitnessGoalRepo := repository.NewFitnessGoalRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager, sessionManager)
	userService := service.NewUserService(userRepo, bodyDataRepo, fitnessGoalRepo)
	aiService := service.NewAIService(
		aiAPIRepo,
		encryptor,
		config.GlobalConfig.AI.RetryAttempts,
		config.GlobalConfig.AI.RetryDelay,
	)
	aiAPIService := service.NewAIAPIService(aiAPIRepo, encryptor)
	trainingService := service.NewTrainingService(
		trainingPlanRepo,
		trainingRecordRepo,
		aiAPIRepo,
		assessmentRepo,
		bodyDataRepo,
		fitnessGoalRepo,
		aiService,
	)
	nutritionService := service.NewNutritionService(
		nutritionPlanRepo,
		nutritionRecordRepo,
		aiAPIRepo,
		bodyDataRepo,
		fitnessGoalRepo,
		aiService,
	)
	statisticsService := service.NewStatisticsService(
		trainingRecordRepo,
		bodyDataRepo,
	)

	return &router.Dependencies{
		DB:                db,
		RedisClient:       redisClient,
		JWTManager:        jwtManager,
		SessionManager:    sessionManager,
		RateLimiter:       rateLimiter,
		AuthService:       authService,
		UserService:       userService,
		AIAPIService:      aiAPIService,
		TrainingService:   trainingService,
		NutritionService:  nutritionService,
		StatisticsService: statisticsService,
		AssessmentRepo:    assessmentRepo,
	}, nil
}

// registerCustomValidators registers custom validation functions with Gin's validator
func registerCustomValidators() error {
	// Get the validator instance from Gin's binding
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators
		if err := v.RegisterValidation("password_strength", customvalidator.ValidatePasswordStrength); err != nil {
			return fmt.Errorf("failed to register password_strength validator: %w", err)
		}
		if err := v.RegisterValidation("email_format", customvalidator.ValidateEmailFormat); err != nil {
			return fmt.Errorf("failed to register email_format validator: %w", err)
		}
		if err := v.RegisterValidation("macro_ratio", customvalidator.ValidateMacroRatio); err != nil {
			return fmt.Errorf("failed to register macro_ratio validator: %w", err)
		}
		if err := v.RegisterValidation("future_date", customvalidator.ValidateNotFutureDate); err != nil {
			return fmt.Errorf("failed to register future_date validator: %w", err)
		}
	}
	return nil
}
