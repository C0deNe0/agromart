package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PreSignedUpload struct {
	URL     string `json:"url"`
	Key     string `json:"key"`
	Expires int64  `json:"expires_in"`
}
type S3Client struct {
	client *s3.Client
	bucket string
}

func NewS3Client(cfg aws.Config, bucket string) *S3Client {
	return &S3Client{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}
}

func (s *S3Client) PreSignedUpload(ctx context.Context, key string, contentType string) (*PreSignedUpload, error) {
	presigner := s3.NewPresignClient(s.client)

	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		ContentType: &contentType,
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return nil, err
	}

	return &PreSignedUpload{
		URL:     req.URL,
		Key:     key,
		Expires: 300,
	}, nil
}

func (s *S3Client) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}
	return nil
}
