package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"qr-file-hub/backend/internal/config"
	"qr-file-hub/backend/internal/model"
	"qr-file-hub/backend/internal/service"
)

type AdminHandler struct {
	cfg     *config.Config
	db      *gorm.DB
	storage *service.StorageService
	bridge  *service.BridgeService
}

func NewAdminHandler(cfg *config.Config, db *gorm.DB, storage *service.StorageService, bridge *service.BridgeService) *AdminHandler {
	return &AdminHandler{
		cfg:     cfg,
		db:      db,
		storage: storage,
		bridge:  bridge,
	}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入管理密码"})
		return
	}

	if req.Password != h.cfg.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理密码错误"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "admin",
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成凭证失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func (h *AdminHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证凭证"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "凭证已过期或无效"})
			return
		}

		c.Next()
	}
}

func (h *AdminHandler) CheckAuth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	var totalFiles int64
	var totalStorage int64
	var totalDownloads int64
	var expiredCount int64

	h.db.Model(&model.FileAsset{}).Where("is_deleted = ?", false).Count(&totalFiles)

	type StorageSum struct {
		Total int64
	}
	var sum StorageSum
	h.db.Model(&model.FileAsset{}).Where("is_deleted = ?", false).Select("COALESCE(SUM(file_size), 0) as total").Scan(&sum)
	totalStorage = sum.Total

	var dlSum StorageSum
	h.db.Model(&model.FileAsset{}).Where("is_deleted = ?", false).Select("COALESCE(SUM(download_count), 0) as total").Scan(&dlSum)
	totalDownloads = dlSum.Total

	now := time.Now()
	h.db.Model(&model.FileAsset{}).Where("is_deleted = ? AND expire_at IS NOT NULL AND expire_at < ?", false, now).Count(&expiredCount)

	c.JSON(http.StatusOK, gin.H{
		"total_files":     totalFiles,
		"total_storage":   totalStorage,
		"total_downloads": totalDownloads,
		"expired_count":   expiredCount,
	})
}

func (h *AdminHandler) GetFiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	mode := strings.TrimSpace(c.Query("mode"))
	status := strings.TrimSpace(c.Query("status"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := h.db.Model(&model.FileAsset{}).Where("is_deleted = ?", false)

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("file_name LIKE ? OR file_id LIKE ?", like, like)
	}

	if mode != "" && mode != "all" {
		query = query.Where("mode = ?", mode)
	}

	now := time.Now()
	if status == "expired" {
		query = query.Where("expire_at IS NOT NULL AND expire_at < ?", now)
	} else if status == "active" {
		query = query.Where("expire_at IS NULL OR expire_at >= ?", now)
	}

	var total int64
	query.Count(&total)

	var list []model.FileAsset
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     list,
	})
}

func (h *AdminHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("fileId")

	var asset model.FileAsset
	if err := h.db.Where("file_id = ?", fileID).First(&asset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	if asset.R2Key != "" {
		_ = h.storage.DeleteObject(context.Background(), asset.R2Key)
	}

	h.db.Model(&asset).Update("is_deleted", true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) BatchDelete(c *gin.Context) {
	var req struct {
		FileIDs []string `json:"file_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.FileIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供要删除的文件ID列表"})
		return
	}

	var assets []model.FileAsset
	h.db.Where("file_id IN (?)", req.FileIDs).Find(&assets)

	for _, a := range assets {
		if a.R2Key != "" {
			_ = h.storage.DeleteObject(context.Background(), a.R2Key)
		}
	}

	h.db.Model(&model.FileAsset{}).Where("file_id IN (?)", req.FileIDs).Update("is_deleted", true)
	c.JSON(http.StatusOK, gin.H{"success": true, "deleted_count": len(assets)})
}

func (h *AdminHandler) CleanExpired(c *gin.Context) {
	now := time.Now()
	var expiredAssets []model.FileAsset
	h.db.Where("is_deleted = ? AND expire_at IS NOT NULL AND expire_at < ?", false, now).Find(&expiredAssets)

	for _, a := range expiredAssets {
		if a.R2Key != "" {
			_ = h.storage.DeleteObject(context.Background(), a.R2Key)
		}
	}

	var ids []string
	for _, a := range expiredAssets {
		ids = append(ids, a.FileID)
	}

	if len(ids) > 0 {
		h.db.Model(&model.FileAsset{}).Where("file_id IN (?)", ids).Update("is_deleted", true)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"cleaned_count": len(ids),
	})
}

func (h *AdminHandler) GetBridges(c *gin.Context) {
	list, err := h.bridge.GetAllBridges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取跳板列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *AdminHandler) ProbeBridges(c *gin.Context) {
	bridges, err := h.bridge.ProbeAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "跳板体检失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bridges})
}

func (h *AdminHandler) CreateBridge(c *gin.Context) {
	var req struct {
		Domain string `json:"domain" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入有效的域名或网址"})
		return
	}

	domain := strings.TrimSpace(req.Domain)
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}

	bridge := model.Bridge{
		Domain:       domain,
		IsActive:     true,
		IsPrimary:    false,
		IsCustom:     true,
		WeChatStatus: "normal",
		QQStatus:     "normal",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.db.Create(&bridge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加跳板失败，可能已存在"})
		return
	}

	go h.bridge.ProbeBridge(&bridge)
	c.JSON(http.StatusOK, gin.H{"data": bridge})
}

func (h *AdminHandler) SetPrimaryBridge(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	h.db.Model(&model.Bridge{}).Where("id != ?", id).Update("is_primary", false)
	h.db.Model(&model.Bridge{}).Where("id = ?", id).Update("is_primary", true)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) DeleteBridge(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	h.db.Delete(&model.Bridge{}, id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
