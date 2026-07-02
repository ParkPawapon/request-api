package httptransport

import (
	"context"

	"github.com/ParkPawapon/request-api/internal/config"
	"github.com/ParkPawapon/request-api/internal/infrastructure/cache"
	"github.com/ParkPawapon/request-api/internal/infrastructure/database"
	dbrepository "github.com/ParkPawapon/request-api/internal/infrastructure/database/repository"
	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/middleware"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	v1health "github.com/ParkPawapon/request-api/internal/transport/http/v1/health"
	v1petitiontypes "github.com/ParkPawapon/request-api/internal/transport/http/v1/petitiontypes"
	petitiontypeusecase "github.com/ParkPawapon/request-api/internal/usecase/petitiontype"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouterDependencies struct {
	Config *config.Config
	DB     *gorm.DB
	Logger *zap.Logger
	Redis  *redis.Client
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	if deps.Config.App.Env == config.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.MaxMultipartMemory = deps.Config.App.MaxBodyBytes

	router.Use(
		middleware.RequestID(),
		middleware.SecurityHeaders(),
		middleware.CORS(deps.Config.CORS.AllowedOrigins),
		middleware.BodyLimit(deps.Config.App.MaxBodyBytes),
		middleware.Timeout(deps.Config.App.RequestTimeout),
		middleware.Logger(deps.Logger),
		middleware.Recovery(deps.Logger),
	)
	router.NoRoute(func(c *gin.Context) {
		response.AppError(c, apperrors.NotFound("Not Found", nil))
	})

	healthHandler := v1health.NewHandler(v1health.Dependencies{
		CheckDatabase: func(ctx context.Context) error {
			return database.Ping(ctx, deps.DB)
		},
		CheckRedis: func(ctx context.Context) error {
			return cache.Ping(ctx, deps.Redis)
		},
		Logger: deps.Logger,
	})

	v1 := router.Group("/v1")
	v1health.RegisterRoutes(v1, healthHandler)

	petitionTypeRepository := dbrepository.NewPetitionTypeGormRepository(deps.DB)
	petitionTypeUseCase := petitiontypeusecase.NewListPetitionTypesUseCase(petitionTypeRepository)
	petitionTypeHandler := v1petitiontypes.NewHandler(petitionTypeUseCase)
	v1petitiontypes.RegisterRoutes(v1, petitionTypeHandler)

	return router
}
