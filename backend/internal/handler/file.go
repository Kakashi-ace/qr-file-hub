package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qr-file-hub/backend/internal/config"
	"qr-file-hub/backend/internal/model"
	"qr-file-hub/backend/internal/service"
)

const (
	maxUploadFileSize   = 100 * 1024 * 1024
	maxMarkdownSize     = 512 * 1024
	maxMarkdownTitleLen = 128
	maxTargetUrlLen     = 2048
)

type FileHandler struct {
	cfg     *config.Config
	db      *gorm.DB
	storage *service.StorageService
	bridge  *service.BridgeService
}

func NewFileHandler(cfg *config.Config, db *gorm.DB, storage *service.StorageService, bridge *service.BridgeService) *FileHandler {
	return &FileHandler{
		cfg:     cfg,
		db:      db,
		storage: storage,
		bridge:  bridge,
	}
}

func generateRandomID(length int) string {
	bytes := make([]byte, length/2+1)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

func calculateExpireAt(expireType string, customSeconds int64) *time.Time {
	now := time.Now()
	var expireAt time.Time
	switch expireType {
	case "1h":
		expireAt = now.Add(time.Hour)
	case "24h":
		expireAt = now.Add(24 * time.Hour)
	case "7d":
		expireAt = now.Add(7 * 24 * time.Hour)
	case "30d":
		expireAt = now.Add(30 * 24 * time.Hour)
	case "custom":
		if customSeconds > 0 && customSeconds <= 365*24*3600 {
			expireAt = now.Add(time.Duration(customSeconds) * time.Second)
		} else {
			return nil
		}
	default:
		return nil
	}
	return &expireAt
}

func (h *FileHandler) InitUpload(c *gin.Context) {
	var req struct {
		FileName string `json:"file_name" binding:"required"`
		FileSize int64  `json:"file_size"`
		MimeType string `json:"mime_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}

	cleanName := filepath.Base(strings.TrimSpace(req.FileName))
	if cleanName == "." || cleanName == "/" || cleanName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法文件名"})
		return
	}

	if len(cleanName) > 200 {
		cleanName = cleanName[:200]
	}

	if req.FileSize > maxUploadFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件过大，单文件上传限制最大为 100MB"})
		return
	}

	fileID := generateRandomID(8)
	safeName := strings.ReplaceAll(cleanName, " ", "_")
	r2Key := fmt.Sprintf("uploads/%s/%s", fileID, safeName)

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	uploadURL, err := h.storage.GeneratePresignedPutURL(context.Background(), r2Key, mimeType, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成上传凭证失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":    fileID,
		"r2_key":     r2Key,
		"upload_url": uploadURL,
	})
}

func (h *FileHandler) CompleteUpload(c *gin.Context) {
	var req struct {
		FileID        string `json:"file_id" binding:"required"`
		FileName      string `json:"file_name" binding:"required"`
		FileSize      int64  `json:"file_size"`
		MimeType      string `json:"mime_type"`
		R2Key         string `json:"r2_key" binding:"required"`
		AllowDownload bool   `json:"allow_download"`
		BridgeDomain  string `json:"bridge_domain"`
		ExpireType    string `json:"expire_type"`
		CustomSeconds int64  `json:"custom_seconds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}

	expectedPrefix := fmt.Sprintf("uploads/%s/", req.FileID)
	if !strings.HasPrefix(req.R2Key, expectedPrefix) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法资源存储标识"})
		return
	}

	if req.FileSize > maxUploadFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件超出系统上限"})
		return
	}

	cleanName := filepath.Base(strings.TrimSpace(req.FileName))
	if cleanName == "." || cleanName == "/" || cleanName == "" {
		cleanName = "uploaded_file"
	}

	expireAt := calculateExpireAt(req.ExpireType, req.CustomSeconds)

	asset := model.FileAsset{
		FileID:        req.FileID,
		FileName:      cleanName,
		FileSize:      req.FileSize,
		MimeType:      req.MimeType,
		R2Key:         req.R2Key,
		Mode:          "file",
		AllowDownload: req.AllowDownload,
		BridgeDomain:  req.BridgeDomain,
		ExpireType:    req.ExpireType,
		ExpireAt:      expireAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.db.Create(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存记录失败", "details": err.Error()})
		return
	}

	directURL := fmt.Sprintf("%s/s/%s", strings.TrimRight(h.cfg.R2PublicDomain, "/"), asset.FileID)
	bridgeURL := directURL
	if asset.BridgeDomain != "" {
		b := strings.TrimRight(asset.BridgeDomain, "/")
		bridgeURL = fmt.Sprintf("%s/?f=%s", b, asset.FileID)
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":    asset.FileID,
		"file_name":  asset.FileName,
		"direct_url": directURL,
		"bridge_url": bridgeURL,
		"expire_at":  asset.ExpireAt,
	})
}

func (h *FileHandler) CreateMarkdown(c *gin.Context) {
	var req struct {
		Title         string `json:"title"`
		Content       string `json:"content" binding:"required"`
		AllowDownload bool   `json:"allow_download"`
		BridgeDomain  string `json:"bridge_domain"`
		ExpireType    string `json:"expire_type"`
		CustomSeconds int64  `json:"custom_seconds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}

	if len(req.Content) > maxMarkdownSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Markdown 内容过大，单次最大限制为 512KB"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if len(title) > maxMarkdownTitleLen {
		title = title[:maxMarkdownTitleLen]
	}
	if title == "" {
		title = fmt.Sprintf("Markdown_%s.md", time.Now().Format("20060102_150405"))
	}
	if !strings.HasSuffix(strings.ToLower(title), ".md") {
		title = title + ".md"
	}
	title = filepath.Base(title)

	fileID := generateRandomID(8)
	r2Key := fmt.Sprintf("markdown/%s/%s", fileID, title)

	_ = h.storage.PutString(context.Background(), r2Key, req.Content, "text/markdown; charset=utf-8")

	expireAt := calculateExpireAt(req.ExpireType, req.CustomSeconds)

	asset := model.FileAsset{
		FileID:        fileID,
		FileName:      title,
		FileSize:      int64(len(req.Content)),
		MimeType:      "text/markdown",
		R2Key:         r2Key,
		Mode:          "text",
		TextContent:   req.Content,
		AllowDownload: req.AllowDownload,
		BridgeDomain:  req.BridgeDomain,
		ExpireType:    req.ExpireType,
		ExpireAt:      expireAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.db.Create(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存记录失败", "details": err.Error()})
		return
	}

	directURL := fmt.Sprintf("%s/s/%s", strings.TrimRight(h.cfg.R2PublicDomain, "/"), asset.FileID)
	bridgeURL := directURL
	if asset.BridgeDomain != "" {
		b := strings.TrimRight(asset.BridgeDomain, "/")
		bridgeURL = fmt.Sprintf("%s/?f=%s", b, asset.FileID)
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":    asset.FileID,
		"file_name":  asset.FileName,
		"direct_url": directURL,
		"bridge_url": bridgeURL,
		"expire_at":  asset.ExpireAt,
	})
}

func (h *FileHandler) CreateUrl(c *gin.Context) {
	var req struct {
		TargetURL     string `json:"target_url" binding:"required"`
		Title         string `json:"title"`
		BridgeDomain  string `json:"bridge_domain"`
		ExpireType    string `json:"expire_type"`
		CustomSeconds int64  `json:"custom_seconds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}

	rawTarget := strings.TrimSpace(req.TargetURL)
	if len(rawTarget) > maxTargetUrlLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目标网址过长"})
		return
	}

	parsed, err := url.ParseRequestURI(rawTarget)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的网址格式"})
		return
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持以 http:// 或 https:// 开头的网址"})
		return
	}

	if parsed.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "网址必须包含有效的主机名"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if len(title) > maxMarkdownTitleLen {
		title = title[:maxMarkdownTitleLen]
	}
	if title == "" {
		title = rawTarget
	}

	fileID := generateRandomID(8)
	expireAt := calculateExpireAt(req.ExpireType, req.CustomSeconds)

	asset := model.FileAsset{
		FileID:       fileID,
		FileName:     title,
		TargetURL:    rawTarget,
		Mode:         "url",
		BridgeDomain: req.BridgeDomain,
		ExpireType:   req.ExpireType,
		ExpireAt:     expireAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.db.Create(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存记录失败", "details": err.Error()})
		return
	}

	directURL := fmt.Sprintf("%s/s/%s", strings.TrimRight(h.cfg.R2PublicDomain, "/"), asset.FileID)
	bridgeURL := directURL
	if asset.BridgeDomain != "" {
		b := strings.TrimRight(asset.BridgeDomain, "/")
		bridgeURL = fmt.Sprintf("%s/?f=%s", b, asset.FileID)
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":    asset.FileID,
		"target_url": asset.TargetURL,
		"direct_url": directURL,
		"bridge_url": bridgeURL,
		"expire_at":  asset.ExpireAt,
	})
}

func (h *FileHandler) GetShareInfo(c *gin.Context) {
	fileID := c.Param("fileId")

	var asset model.FileAsset
	if err := h.db.Where("file_id = ? AND is_deleted = ?", fileID, false).First(&asset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在或已被删除"})
		return
	}

	if asset.IsExpired() {
		c.JSON(http.StatusGone, gin.H{"error": "此分享已过期失效"})
		return
	}

	h.db.Model(&asset).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	ua := c.Request.UserAgent()
	isWeChat := strings.Contains(strings.ToLower(ua), "micromessenger")
	isQQ := strings.Contains(strings.ToLower(ua), "qq/")

	h.db.Create(&model.AccessLog{
		FileID:    asset.FileID,
		IP:        c.ClientIP(),
		UserAgent: ua,
		IsWeChat:  isWeChat,
		IsQQ:      isQQ,
		Referer:   c.Request.Referer(),
		Action:    "view",
		CreatedAt: time.Now(),
	})

	downloadURL := ""
	if asset.AllowDownload && asset.R2Key != "" {
		dUrl, err := h.storage.GeneratePresignedGetURL(context.Background(), asset.R2Key, asset.FileName, 1*time.Hour)
		if err == nil {
			downloadURL = dUrl
		}
	}

	publicURL := ""
	if asset.R2Key != "" {
		publicURL = h.storage.GetPublicURL(asset.R2Key)
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":        asset.FileID,
		"file_name":      asset.FileName,
		"file_size":      asset.FileSize,
		"mime_type":      asset.MimeType,
		"mode":           asset.Mode,
		"text_content":   asset.TextContent,
		"target_url":     asset.TargetURL,
		"allow_download": asset.AllowDownload,
		"download_url":   downloadURL,
		"public_url":     publicURL,
		"view_count":     asset.ViewCount + 1,
		"created_at":     asset.CreatedAt,
		"expire_at":      asset.ExpireAt,
	})
}

func (h *FileHandler) DownloadFile(c *gin.Context) {
	fileID := c.Param("fileId")

	var asset model.FileAsset
	if err := h.db.Where("file_id = ? AND is_deleted = ?", fileID, false).First(&asset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	if asset.IsExpired() {
		c.JSON(http.StatusGone, gin.H{"error": "此分享已过期失效"})
		return
	}

	if !asset.AllowDownload {
		c.JSON(http.StatusForbidden, gin.H{"error": "发布者未开放文件下载权限"})
		return
	}

	h.db.Model(&asset).UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))

	ua := c.Request.UserAgent()
	h.db.Create(&model.AccessLog{
		FileID:    asset.FileID,
		IP:        c.ClientIP(),
		UserAgent: ua,
		IsWeChat:  strings.Contains(strings.ToLower(ua), "micromessenger"),
		IsQQ:      strings.Contains(strings.ToLower(ua), "qq/"),
		Referer:   c.Request.Referer(),
		Action:    "download",
		CreatedAt: time.Now(),
	})

	url, err := h.storage.GeneratePresignedGetURL(context.Background(), asset.R2Key, asset.FileName, 30*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成下载链接失败"})
		return
	}

	c.Redirect(http.StatusFound, url)
}

func (h *FileHandler) GetBridges(c *gin.Context) {
	list, err := h.bridge.GetActiveBridges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取跳板列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
