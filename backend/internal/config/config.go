package config

import (
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
		dsn = "VYMxcUs7ZvJpiXQ.root:AG251314@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/test?charset=utf8mb4&parseTime=True&loc=Local&tls=true"
	}

	accountID := os.Getenv("R2_ACCOUNT_ID")
	if accountID == "" {
		accountID = "518575a1761a9b2559753fa14343683b"
	}

	keyID := os.Getenv("R2_ACCESS_KEY_ID")
	if keyID == "" {
		keyID = "f24c876ba30372d30157531bc8ea89e2"
	}

	keySecret := os.Getenv("R2_ACCESS_KEY_SECRET")
	if keySecret == "" {
		keySecret = "9fea391ff1c4d8fdc0359b114c97a5470d4d725facd4809c6f2be46895418c1f"
	}

	bucket := os.Getenv("R2_BUCKET")
	if bucket == "" {
		bucket = "media-storage"
	}

	endpoint := os.Getenv("R2_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://518575a1761a9b2559753fa14343683b.r2.cloudflarestorage.com"
	}

	publicDomain := os.Getenv("R2_PUBLIC_DOMAIN")
	if publicDomain == "" {
		publicDomain = "https://share.k52h.cc"
	}

	adminPwd := os.Getenv("ADMIN_PASSWORD")
	if adminPwd == "" {
		adminPwd = "AG251314"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "AG251314_SECURE_TOKEN_KEY_987654321"
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
