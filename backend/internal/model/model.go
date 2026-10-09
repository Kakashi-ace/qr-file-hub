package model

import (
	"time"
)

type FileAsset struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	FileID        string     `gorm:"size:64;uniqueIndex;not null" json:"file_id"`
	FileName      string     `gorm:"size:255;not null" json:"file_name"`
	FileSize      int64      `gorm:"default:0" json:"file_size"`
	MimeType      string     `gorm:"size:128;default:''" json:"mime_type"`
	R2Key         string     `gorm:"size:255;default:''" json:"r2_key"`
	Mode          string     `gorm:"size:32;not null;default:'file'" json:"mode"`
	TextContent   string     `gorm:"type:longtext" json:"text_content,omitempty"`
	TargetURL     string     `gorm:"size:1024;default:''" json:"target_url,omitempty"`
	AllowDownload bool       `gorm:"default:false" json:"allow_download"`
	BridgeDomain  string     `gorm:"size:255;default:''" json:"bridge_domain"`
	ExpireType    string     `gorm:"size:32;default:'never'" json:"expire_type"`
	ExpireAt      *time.Time `gorm:"index" json:"expire_at"`
	ViewCount     int64      `gorm:"default:0" json:"view_count"`
	DownloadCount int64      `gorm:"default:0" json:"download_count"`
	IsDeleted     bool       `gorm:"default:false;index" json:"is_deleted"`
	CreatedAt     time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (f *FileAsset) IsExpired() bool {
	if f.ExpireAt == nil {
		return false
	}
	return time.Now().After(*f.ExpireAt)
}

type Bridge struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Domain        string     `gorm:"size:255;uniqueIndex;not null" json:"domain"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	IsPrimary     bool       `gorm:"default:false" json:"is_primary"`
	IsCustom      bool       `gorm:"default:false" json:"is_custom"`
	WeChatStatus  string     `gorm:"size:32;default:'normal'" json:"wechat_status"`
	QQStatus      string     `gorm:"size:32;default:'normal'" json:"qq_status"`
	LatencyMs     int64      `gorm:"default:0" json:"latency_ms"`
	LastCheckedAt *time.Time `json:"last_checked_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type AccessLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FileID    string    `gorm:"size:64;index;not null" json:"file_id"`
	IP        string    `gorm:"size:64;default:''" json:"ip"`
	UserAgent string    `gorm:"size:512;default:''" json:"user_agent"`
	IsWeChat  bool      `gorm:"default:false" json:"is_wechat"`
	IsQQ      bool      `gorm:"default:false" json:"is_qq"`
	Referer   string    `gorm:"size:512;default:''" json:"referer"`
	Action    string    `gorm:"size:32;default:'view'" json:"action"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
