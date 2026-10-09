package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"qr-file-hub/backend/internal/config"
)

type StorageService struct {
	cfg           *config.Config
	s3Client      *s3.Client
	presignClient *s3.PresignClient
}

func NewStorageService(cfg *config.Config) (*StorageService, error) {
	creds := credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2AccessKeySecret, "")

	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithCredentialsProvider(creds),
		awsConfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.R2Endpoint)
		o.UsePathStyle = true
	})
	presign := s3.NewPresignClient(client)

	return &StorageService{
		cfg:           cfg,
		s3Client:      client,
		presignClient: presign,
	}, nil
}

func (s *StorageService) GeneratePresignedPutURL(ctx context.Context, r2Key string, mimeType string, expireDuration time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.R2Bucket),
		Key:         aws.String(r2Key),
		ContentType: aws.String(mimeType),
	}

	req, err := s.presignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(expireDuration))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *StorageService) GeneratePresignedGetURL(ctx context.Context, r2Key string, fileName string, expireDuration time.Duration) (string, error) {
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", fileName)
	input := &s3.GetObjectInput{
		Bucket:                     aws.String(s.cfg.R2Bucket),
		Key:                        aws.String(r2Key),
		ResponseContentDisposition: aws.String(contentDisposition),
	}

	req, err := s.presignClient.PresignGetObject(ctx, input, s3.WithPresignExpires(expireDuration))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *StorageService) PutString(ctx context.Context, r2Key string, content string, mimeType string) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.R2Bucket),
		Key:         aws.String(r2Key),
		Body:        strings.NewReader(content),
		ContentType: aws.String(mimeType),
	}
	_, err := s.s3Client.PutObject(ctx, input)
	return err
}

func (s *StorageService) DeleteObject(ctx context.Context, r2Key string) error {
	if r2Key == "" {
		return nil
	}
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.R2Bucket),
		Key:    aws.String(r2Key),
	}
	_, err := s.s3Client.DeleteObject(ctx, input)
	return err
}

func (s *StorageService) GetPublicURL(r2Key string) string {
	domain := strings.TrimRight(s.cfg.R2PublicDomain, "/")
	return fmt.Sprintf("%s/%s", domain, r2Key)
}
