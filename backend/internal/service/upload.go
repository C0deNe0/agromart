package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/C0deNe0/agromart/internal/lib/aws"
	"github.com/C0deNe0/agromart/internal/model/upload"
	"github.com/google/uuid"
)

type UploadService struct {
	s3             *aws.S3Client
	ProductService *ProductService
}

func NewUploadService(s3 *aws.S3Client, productService *ProductService) *UploadService {
	return &UploadService{
		s3:             s3,
		ProductService: productService,
	}
}

func (s *UploadService) PresignUpload(ctx context.Context, req upload.UploadRequest) (*upload.PresignedUpload, error) {

	if !strings.HasPrefix(req.ContentType, "image/") {
		return nil, errors.New("invalid content type")
	}

	switch req.Type {
	case upload.UploadProductImage:
		if req.ProductID == nil {
			return nil, errors.New("product id is required")
		}

		if err := s.ProductService.AuthorizeProductMutation(ctx, req.UserID, *req.ProductID); err != nil {
			return nil, err
		}

		//server generated key
		key := fmt.Sprintf("products/%s/images/%s", req.ProductID, uuid.New().String())

		//presigned 
		presigned,err := s.s3.PreSignedUpload(
			ctx,
			key,
			req.ContentType,
		)
		if err != nil {
			return nil, err
		}

		return &upload.PresignedUpload{
			URL:      presigned.URL,
			Key:      presigned.Key,
			ExpireAt: presigned.Expires,
		}, nil
	default:
		return nil, errors.New("unsupported upload type")
	}
}
