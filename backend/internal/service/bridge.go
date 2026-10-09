package service

import (
	"crypto/tls"
	"net/http"
	"time"

	"gorm.io/gorm"

	"qr-file-hub/backend/internal/model"
)

type BridgeService struct {
	db *gorm.DB
}

func NewBridgeService(db *gorm.DB) *BridgeService {
	return &BridgeService{db: db}
}

func (s *BridgeService) GetActiveBridges() ([]model.Bridge, error) {
	var list []model.Bridge
	err := s.db.Where("is_active = ?", true).Order("is_primary DESC, id ASC").Find(&list).Error
	return list, err
}

func (s *BridgeService) GetAllBridges() ([]model.Bridge, error) {
	var list []model.Bridge
	err := s.db.Order("is_primary DESC, id ASC").Find(&list).Error
	return list, err
}

func (s *BridgeService) ProbeBridge(bridge *model.Bridge) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	start := time.Now()
	resp, err := client.Get(bridge.Domain)
	duration := time.Since(start).Milliseconds()

	now := time.Now()
	bridge.LastCheckedAt = &now
	bridge.LatencyMs = duration

	if err != nil || resp.StatusCode >= 400 {
		bridge.WeChatStatus = "abnormal"
		bridge.QQStatus = "abnormal"
	} else {
		resp.Body.Close()
		bridge.WeChatStatus = "normal"
		bridge.QQStatus = "normal"
	}

	s.db.Save(bridge)
}

func (s *BridgeService) ProbeAll() ([]model.Bridge, error) {
	var bridges []model.Bridge
	if err := s.db.Find(&bridges).Error; err != nil {
		return nil, err
	}

	for i := range bridges {
		s.ProbeBridge(&bridges[i])
	}

	return bridges, nil
}
