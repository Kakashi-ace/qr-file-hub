package handler

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"qr-file-hub/backend/internal/config"
	"qr-file-hub/backend/internal/model"
)

type AuthLimiter struct {
	sync.Mutex
	attempts map[string][]time.Time
}

func NewAuthLimiter() *AuthLimiter {
	return &AuthLimiter{
		attempts: make(map[string][]time.Time),
	}
}

func (l *AuthLimiter) Allow(ip string) bool {
	l.Lock()
	defer l.Unlock()
	now := time.Now()
	cutoff := now.Add(-15 * time.Minute)
	var valid []time.Time
	for _, t := range l.attempts[ip] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	l.attempts[ip] = valid
	return len(valid) < 5
}

func (l *AuthLimiter) RecordFail(ip string) {
	l.Lock()
	defer l.Unlock()
	l.attempts[ip] = append(l.attempts[ip], time.Now())
}

func (l *AuthLimiter) Reset(ip string) {
	l.Lock()
	defer l.Unlock()
	delete(l.attempts, ip)
}

type AuthHandler struct {
	cfg     *config.Config
	db      *gorm.DB
	limiter *AuthLimiter
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		cfg:     cfg,
		db:      db,
		limiter: NewAuthLimiter(),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入登录账号与密码"})
		return
	}

	clientIP := c.ClientIP()
	if !h.limiter.Allow(clientIP) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "密码错误次数过多，请 15 分钟后再试"})
		return
	}

	username := strings.TrimSpace(req.Username)
	if username == "" {
		username = "admin"
	}

	var user model.User
	if err := h.db.Where("username = ?", username).First(&user).Error; err != nil {
		if username == "admin" && subtle.ConstantTimeCompare([]byte(req.Password), []byte(h.cfg.AdminPassword)) == 1 {
			hash, _ := bcrypt.GenerateFromPassword([]byte(h.cfg.AdminPassword), bcrypt.DefaultCost)
			user = model.User{
				Username:     "admin",
				PasswordHash: string(hash),
				Role:         "admin",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			h.db.Create(&user)
		} else {
			h.limiter.RecordFail(clientIP)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
			return
		}
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			if username == "admin" && subtle.ConstantTimeCompare([]byte(req.Password), []byte(h.cfg.AdminPassword)) == 1 {
				hash, _ := bcrypt.GenerateFromPassword([]byte(h.cfg.AdminPassword), bcrypt.DefaultCost)
				h.db.Model(&user).Update("password_hash", string(hash))
			} else {
				h.limiter.RecordFail(clientIP)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
				return
			}
		}
	}

	h.limiter.Reset(clientIP)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成凭证失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "请先登录后再访问"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(h.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录凭证已失效，请重新登录"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的身份凭证"})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func (h *AuthHandler) AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足，需要管理员权限"})
			return
		}
		c.Next()
	}
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
		"role":     role,
	})
}
