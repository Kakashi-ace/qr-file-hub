package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"qr-file-hub/backend/internal/config"
	"qr-file-hub/backend/internal/handler"
	"qr-file-hub/backend/internal/service"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB, storage *service.StorageService, bridge *service.BridgeService) *gin.Engine {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	fileH := handler.NewFileHandler(cfg, db, storage, bridge)
	adminH := handler.NewAdminHandler(cfg, db, storage, bridge)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "qr-file-hub-backend"})
	})

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.POST("/upload/presign", fileH.InitUpload)
		api.POST("/upload/complete", fileH.CompleteUpload)
		api.POST("/markdown/create", fileH.CreateMarkdown)
		api.POST("/url/create", fileH.CreateUrl)

		api.GET("/share/:fileId", fileH.GetShareInfo)
		api.GET("/share/:fileId/download", fileH.DownloadFile)
		api.GET("/bridges", fileH.GetBridges)

		admin := api.Group("/admin")
		{
			admin.POST("/login", adminH.Login)

			authed := admin.Group("")
			authed.Use(adminH.AuthMiddleware())
			{
				authed.GET("/check", adminH.CheckAuth)
				authed.GET("/stats", adminH.GetStats)
				authed.GET("/files", adminH.GetFiles)
				authed.DELETE("/files/:fileId", adminH.DeleteFile)
				authed.POST("/files/batch-delete", adminH.BatchDelete)
				authed.POST("/files/clean-expired", adminH.CleanExpired)

				authed.GET("/bridges", adminH.GetBridges)
				authed.POST("/bridges/probe", adminH.ProbeBridges)
				authed.POST("/bridges", adminH.CreateBridge)
				authed.PUT("/bridges/:id/primary", adminH.SetPrimaryBridge)
				authed.DELETE("/bridges/:id", adminH.DeleteBridge)
			}
		}
	}

	return r
}
