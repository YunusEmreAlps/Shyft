package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"

	docs "shyft/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"shyft/config"
	"shyft/internal/handlers"
	"shyft/pkg/db/postgres"
	"shyft/pkg/db/redis"
	"shyft/pkg/logger"
)

// @title Shift Scheduler Service API
// @description: Shift Scheduler Service API
// @version 1.0.0
// @schemes http https

// @contact.name   Yunus Emre Alpu
// @contact.url    https://yunusemrealpu.netlify.app
// @contact.email  YunusAlpu@icloud.com

// @BasePath /shyft

var isConfigSuccess = false

// var equals string = strings.Repeat("=", 50)

// APP_NAME = "localhost:9097/shyft/"
const (
	APP_NAME = "shyft"
)

func main() {
	mode := config.C.App.Mode
	port := config.C.App.Port

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(gin.DefaultErrorWriter))
	inAppCache := redis.NewInAppCacheStore(time.Minute)
	cacheConn, cacheContext := redis.NewRedisCacheConnection(config.C.Cache.Url)
	dbConn := postgres.NewPostgresDB(config.C.DB.Url)

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: config.C.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           config.C.Jaeger.LogSpans,
			LocalAgentHostPort: config.C.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		logger.CLogger.Fatalf("Cannot create tracer: %v", err)
	}

	// create application service
	shiftsvc := handlers.NewShiftService(
		inAppCache,
		cacheConn,
		cacheContext,
		dbConn,
	)

	// check env and set gin mode
	setApplicationMode(mode, router)
	shiftsvc.InitRouter(router)

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	logger.CLogger.Info("Tracing enabled: Jaeger host=", config.C.Jaeger.Host, " service=", config.C.Jaeger.ServiceName)

	logger.CLogger.Infof("Configuration loaded: service=%s mode=%s port=%s version=%s", config.C.App.Name, config.C.App.Mode, config.C.App.Port, config.C.App.Version)
	logger.CLogger.Infof("Database connected: %s", config.C.DB.Url)
	logger.CLogger.Infof("Cache connected: %s", config.C.Cache.Url)
	logger.CLogger.Infof("Broker connected: %s", config.C.Broker.Url)
	logger.CLogger.Infof("Logger initialized: level=%s encoding=%s development=%v", config.C.Logger.Level, config.C.Logger.Encoding, config.C.Logger.Development)

	// run application
	logger.CLogger.Infof("Application %s started on port %s", APP_NAME, port)
	if gin.Mode() == gin.DebugMode {
		logger.CLogger.Warn("Running in debug mode. Set GIN_MODE=release for production.")
	}
	if err := router.Run(":" + port); err != nil {
		logger.CLogger.Fatalf("Failed to start server: %v", err)
	}
}

// Initialize Application
func init() {
	isConfigSuccess = configureApplication()
	if !isConfigSuccess {
		logger.CLogger.Error("Application configuration failed.")
		os.Exit(1)
	} else {
		logger.CLogger.Info("Application configuration loaded successfully.")
	}
}

// Configure Application with config file
func configureApplication() bool {
	dir, err := os.Getwd()
	if err != nil {
		logger.CLogger.Error("INIT: Cannot get current working directory os.Getwd()")
		return false
	} else {
		config.ReadConfig(dir)
		return true
	}
}

// Set Application Mode
func setApplicationMode(md string, router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	if md == "prod" || md == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		router.Use(gin.Logger())
		// logger time gonna be in UTC+3 (Asia/Istanbul)
		gin.SetMode(gin.DebugMode)
	}

	// check env and set swagger
	if !(md == "prod" || md == "production") {
		docs.SwaggerInfo.BasePath = handlers.API_PREFIX
		// Endpoint for swagger: http://localhost:9097/shyft/swagger/index.html
		router.GET(handlers.API_PREFIX+"/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
