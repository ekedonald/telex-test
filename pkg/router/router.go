package router

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/hngprojects/telex_be/internal/config"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/utility"
)

func Setup(logger *utility.Logger, validator *validator.Validate, db *storage.Database, appConfiguration *config.App) *gin.Engine {
	if appConfiguration.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// Middlewares
	// r.Use(gin.Logger())
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(config.GetConfig().Server.TrustedProxies)
	r.Use(middleware.Security())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.MaxMultipartMemory = 1 << 20 // 1MB

	// routers
	ApiVersion := "api/v1"
	Health(r, ApiVersion, validator, db, logger)
	Auth(r, ApiVersion, validator, db, logger)
	Room(r, ApiVersion, validator, db, logger)
	TokenGen(r, ApiVersion, validator, db, logger)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "HNGi Golang Telex BE",
			"status":  http.StatusOK,
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		})
	})

	r.StaticFile("/swagger.yaml", "static/swagger.yaml")
	url := ginSwagger.URL("/swagger.yaml")
	r.GET("/api/docs/*any", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'sha256-2TOI2ugkuROHHfKZr6kdGv+XxhrVUI8uHycXqXUIR4g='; img-src 'self' data:;")
		ginSwagger.WrapHandler(swaggerFiles.Handler, url)(c)
	})

	return r
}
