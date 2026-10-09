package config

import (
	"log"
	"os"
)

type Config struct {
	Port              string
	DatabaseDSN       string
	R2AccountID       string
	R2AccessKeyID     string
	R2AccessKeySecret string
	R2Bucket          string
	R2Endpoint        string
	R2PublicDomain    string
	AdminPassword     string
	JWTSecret         string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN environment variable is required")
	}

	accountID := os.Getenv("R2_ACCOUNT_ID")
	if accountID == "" {
		log.Fatal("R2_ACCOUNT_ID environment variable is required")
	}

	keyID := os.Getenv("R2_ACCESS_KEY_ID")
	if keyID == "" {
		log.Fatal("R2_ACCESS_KEY_ID environment variable is required")
	}

	keySecret := os.Getenv("R2_ACCESS_KEY_SECRET")
	if keySecret == "" {
		log.Fatal("R2_ACCESS_KEY_SECRET environment variable is required")
	}

	bucket := os.Getenv("R2_BUCKET")
	if bucket == "" {
		bucket = "media-storage"
	}

	endpoint := os.Getenv("R2_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://" + accountID + ".r2.cloudflarestorage.com"
	}

	publicDomain := os.Getenv("R2_PUBLIC_DOMAIN")
	if publicDomain == "" {
		publicDomain = "https://qr.k52h.cc"
	}

	adminPwd := os.Getenv("ADMIN_PASSWORD")
	if adminPwd == "" {
		log.Fatal("ADMIN_PASSWORD environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &Config{
		Port:              port,
		DatabaseDSN:       dsn,
		R2AccountID:       accountID,
		R2AccessKeyID:     keyID,
		R2AccessKeySecret: keySecret,
		R2Bucket:          bucket,
		R2Endpoint:        endpoint,
		R2PublicDomain:    publicDomain,
		AdminPassword:     adminPwd,
		JWTSecret:         jwtSecret,
	}
}
