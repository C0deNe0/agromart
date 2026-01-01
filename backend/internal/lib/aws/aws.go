package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
)

type AWS struct {
	S3 *S3Client
}

func NewAWS(region, bucket string) (*AWS, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return &AWS{
		S3: NewS3Client(cfg, bucket),
	}, nil
}
