package handlers

import (
	"context"
	"net/http"
	"shift-scheduler-service/config"
	"shift-scheduler-service/pkg/logger"
	"shift-scheduler-service/pkg/metric"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	API_PREFIX = "/shift-scheduler-service"
	RN_PREFIX  = "cld:::shiftservice:::"
)

// MyShiftService is a struct for auth core
type ShiftService struct {
	inAppCache   *persistence.InMemoryStore
	cache        *redis.Client
	cacheContext context.Context
	db           *gorm.DB
	// s3sess       *session.Session
}

func NewShiftService(
	inAppCache *persistence.InMemoryStore,
	cache *redis.Client,
	cacheContext context.Context,
	db *gorm.DB,
	// s3sess *session.Session,
) *ShiftService {
	return &ShiftService{
		inAppCache:   inAppCache,
		cache:        cache,
		cacheContext: cacheContext,
		db:           db,
		// s3sess:       s3sess,
	}
}

type RespondJson struct {
	Status  bool        `json:"status"`
	Intent  string      `json:"intent"`
	Message interface{} `json:"message"`
}

func respondJson(ctx *gin.Context, code int, intent string, message interface{}, err error) {
	if err == nil {
		ctx.JSON(code, RespondJson{
			Status:  true,
			Intent:  intent,
			Message: message,
		})
	} else {
		ctx.JSON(code, RespondJson{
			Status:  false,
			Intent:  intent,
			Message: err.Error(),
		})
	}
}

func (bs *ShiftService) InitRouter(r *gin.Engine) {
	// Prometheus metrics
	metrics, err := metric.CreateMetrics(config.C.Metric.Url, config.C.Metric.Service)
	if err != nil {
		logger.CLogger.Warn("Cannot create metrics: ", err)
	} else {
		logger.CLogger.Infof("Metrics server running. Metrics: %+v", metrics)
	}

	// -- my service routes (group)
	v1 := r.Group(API_PREFIX)
	health := v1.Group("/health")

	// Get all shift schedules
	v1.GET("/shift-schedules", func(ctx *gin.Context) {
		code, data, err := bs.HandleGetAllShiftSchedules(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules", data, err)
	})

	// Get shift schedule by id
	v1.GET("/shift-schedules/:id", func(ctx *gin.Context) {
		code, data, err := bs.HandleGetShiftScheduleByID(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules/:id", data, err)
	})

	// Get shift schedule by year
	v1.GET("/shift-schedules/year/:year", func(ctx *gin.Context) {
		code, data, err := bs.HandleGetShiftScheduleByYear(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules/year/:year", data, err)
	})

	// Get shift schedule by current week
	v1.GET("/shift-schedules/week", func(ctx *gin.Context) {
		code, data, err := bs.HandleGetShiftScheduleByWeek(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules/week", data, err)
	})

	// Get Get shift schedule by current week with pagination
	v1.GET("/shift-schedules/week/paginated", func(ctx *gin.Context) {
		code, data, err := bs.HandleGetShiftScheduleByWeekWithPagination(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules/week/paginated", data, err)
	})

	// Get deleted shift schedules
	v1.GET("/shift-schedules/deleted", func(ctx *gin.Context) {
		code, data, err := bs.HandleGetOnlyDeletedShiftSchedules(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules/deleted", data, err)
	})

	// Create shift schedule
	v1.POST("/shift-schedules", func(ctx *gin.Context) {
		code, data, err := bs.HandleCreateShiftSchedule(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules", data, err)
	})

	// Update shift schedule
	v1.PUT("/shift-schedules/:id", func(ctx *gin.Context) {
		code, data, err := bs.HandleUpdateShiftSchedule(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules/:id", data, err)
	})

	// Delete shift schedule (Soft delete)
	v1.DELETE("/shift-schedules/:id", func(ctx *gin.Context) {
		code, data, err := bs.HandleDeleteShiftSchedule(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules/:id", data, err)
	})

	// Restore shift schedule
	v1.PATCH("/shift-schedules/:id/restore", func(ctx *gin.Context) {
		code, data, err := bs.HandleRestoreShiftSchedule(ctx)
		respondJson(ctx, code, RN_PREFIX+"/shift-schedules/:id/restore", data, err)
	})

	// Health check
	health.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Shift Scheduler Service " + config.C.App.Version + " is running on port " + config.C.App.Port + ".",
		})
	})
}
