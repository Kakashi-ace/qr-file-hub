package router

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qr-file-hub/backend/internal/config"
	"qr-file-hub/backend/internal/handler"
	"qr-file-hub/backend/internal/service"
)

type rateLimiter struct {
	sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *rateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rl.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)
		var valid []time.Time
		for _, t := range rl.requests[ip] {
			if t.After(cutoff) {
				valid = append(valid, t)
			}
		}
		if len(valid) >= rl.limit {
			rl.requests[ip] = valid
			rl.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁，请稍后再试"})
			return
		}
		rl.requests[ip] = append(valid, now)
		rl.Unlock()
		c.Next()
	}
}

func SetupRouter(cfg *config.Config, db *gorm.DB, storage *service.StorageService, bridge *service.BridgeService) *gin.Engine {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins: []string{
			"https://qr.k52h.cc",
			"https://qr-file-hub.pages.dev",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	fileH := handler.NewFileHandler(cfg, db, storage, bridge)
	adminH := handler.NewAdminHandler(cfg, db, storage, bridge)

	creationLimiter := newRateLimiter(60, time.Minute)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "qr-file-hub-backend"})
	})

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.POST("/upload/presign", creationLimiter.Middleware(), fileH.InitUpload)
		api.POST("/upload/complete", fileH.CompleteUpload)
		api.POST("/markdown/create", creationLimiter.Middleware(), fileH.CreateMarkdown)
		api.POST("/url/create", creationLimiter.Middleware(), fileH.CreateUrl)

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
