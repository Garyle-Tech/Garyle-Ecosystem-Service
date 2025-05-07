package ota

import (
	"context"

	otaModel "ecosystem.garyle/service/internal/domain/model/ota"
)

type OTARepository interface {
	Create(ctx context.Context, ota *otaModel.OTA) (*otaModel.OTA, error)
	GetByAppID(ctx context.Context, appID string) (*otaModel.OTA, error)
	List(ctx context.Context, limit, page int) ([]*otaModel.OTA, error)
	Count(ctx context.Context) (int, error)
	UpdateByAppID(ctx context.Context, ota *otaModel.OTA, appID string) error
	DeleteByAppID(ctx context.Context, appID string) error
}
