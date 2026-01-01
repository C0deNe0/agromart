package upload

import (
	"github.com/google/uuid"
)

type UploadType string

const (
	UploadProductImage UploadType = "product_image"
)

type UploadRequest struct {
	Type        UploadType
	UserID      uuid.UUID
	ProductID   *uuid.UUID
	ContentType string
}

type PresignedUpload struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	ExpireAt int64  `json:"expireAt"`
}
