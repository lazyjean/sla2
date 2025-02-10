package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/config"
	_ "github.com/lazyjean/sla2/docs" // 导入 swagger docs
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/infrastructure/cache/redis"
	"github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/lazyjean/sla2/internal/interfaces/http/handler"
	"github.com/lazyjean/sla2/internal/interfaces/http/routes"
	"github.com/lazyjean/sla2/pkg/logger"
	"github.com/lazyjean/sla2/pkg/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title        生词本 API
// @version      1.0
// @description  生词本服务 API 文档

// @contact.name   LazyJean
// @contact.email  lazyjean@foxmail.com

// @host      localhost:9000
// @BasePath  /api/v1
// @schemes   http

// @securityDefinitions.apikey  Bearer
// @in                         header
// @name                       Authorization
// @description               Bearer token for authentication

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		logger.Log.Fatal("Failed to initialize config", zap.Error(err))
	}

	// 加载配置
	cfg := config.GetConfig()

	// 初始化日志
	logger.InitLogger(&cfg.Log)
	defer logger.Log.Sync()

	// 设置 gin 的日志输出
	gin.DefaultWriter = logger.NewGinLogger()
	gin.DefaultErrorWriter = gin.DefaultWriter

	// 禁用 gin 的控制台颜色
	gin.DisableConsoleColor()

	// 设置gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化 Swagger 配置
	swagger.InitSwagger()

	// 创建路由 - 使用 New() 而不是 Default()
	r := gin.New()

	// 使用自定义的日志中间件
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: logger.NewGinLogger(),
	}))

	// 使用 Recovery 中间件
	r.Use(gin.Recovery())

	// 初始化数据库
	db, err := postgres.NewDB(&cfg.Database)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database: " + err.Error())
	}

	// 获取底层 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Fatal("Failed to get underlying *sql.DB: " + err.Error())
	}
	defer sqlDB.Close()

	// 初始化Redis缓存
	redisCache, err := redis.NewRedisCache(&cfg.Redis)
	if err != nil {
		logger.Log.Fatal("Failed to connect to Redis: " + err.Error())
	}
	defer redisCache.Close()

	// 初始化仓储
	baseWordRepo := postgres.NewWordRepository(db)
	wordRepo := postgres.NewCachedWordRepository(baseWordRepo, redisCache)
	learningRepo := postgres.NewLearningRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// 初始化应用服务
	wordService := service.NewWordService(wordRepo)
	learningService := service.NewLearningService(learningRepo)
	authService := service.NewAuthService(userRepo)

	// 初始化处理器
	wordHandler := handler.NewWordHandler(wordService)
	authHandler := handler.NewAuthHandler(authService)
	learningHandler := handler.NewLearningHandler(learningService)
	healthHandler := handler.NewHealthHandler()

	handlers := handler.NewHandlers(wordHandler, authHandler, learningHandler, healthHandler)

	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册业务路由
	routes.SetupRoutes(r, handlers)

	// 注册用户路由
	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)

	// 创建服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		// 监听系统信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Log.Info("Shutting down server...")

		// 创建一个5秒的上下文用于超时控制
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Fatal("Server forced to shutdown: " + err.Error())
		}

		logger.Log.Info("Server exiting")
	}()

	// 启动服务器
	logger.Log.Info("Server starting on port " + cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("Failed to start server: " + err.Error())
	}
}
