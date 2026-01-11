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
type S3Service struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Service(client *s3.Client, bucket string, region string) *S3Service {
	return &S3Service{
		client: client,
		bucket: bucket,
		region: region,
	}
}

func (s *S3Service) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiresInSeconds int) (string, error) {
	presigner := s3.NewPresignClient(s.client)

	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expiresInSeconds) * time.Second
	})

	if err != nil {
		return "", fmt.Errorf("failure to generated presiged URL: %w", err)
	}

	return req.URL, nil
}

// GET PUBLIC URL
func (s *S3Service) GetPublicURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
}

func (s *S3Service) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}
	return nil
}
