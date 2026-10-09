package service

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
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

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return false
	}
	if ipv4[0] == 10 {
		return true
	}
	if ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31 {
		return true
	}
	if ipv4[0] == 192 && ipv4[1] == 168 {
		return true
	}
	if ipv4[0] == 169 && ipv4[1] == 254 {
		return true
	}
	return false
}

func (s *BridgeService) IsSafeURL(rawURL string) bool {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return false
	}
	hostname := u.Hostname()
	if hostname == "" {
		return false
	}
	if ip := net.ParseIP(hostname); ip != nil {
		if isPrivateIP(ip) {
			return false
		}
	} else {
		ips, err := net.LookupIP(hostname)
		if err != nil {
			return false
		}
		for _, ip := range ips {
			if isPrivateIP(ip) {
				return false
			}
		}
	}
	return true
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
	if !s.IsSafeURL(bridge.Domain) {
		now := time.Now()
		bridge.LastCheckedAt = &now
		bridge.WeChatStatus = "abnormal"
		bridge.QQStatus = "abnormal"
		s.db.Save(bridge)
		return
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
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
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
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
