package ota

import (
	"context"
	"errors"
	"fmt"

	otaModel "ecosystem.garyle/service/internal/domain/model/ota"
	otaRepo "ecosystem.garyle/service/internal/domain/repository/ota"
)

// Service defines OTA business logic operations
type Service interface {
	Create(ctx context.Context, ota *otaModel.OTA) (*otaModel.OTA, error)
	GetByAppID(ctx context.Context, appID string) (*otaModel.OTA, error)
	List(ctx context.Context, limit, page int) ([]*otaModel.OTA, error)
	Count(ctx context.Context) (int, error)
	UpdateByAppID(ctx context.Context, ota *otaModel.OTA, appID string) error
	DeleteByAppID(ctx context.Context, appID string) error
}

type service struct {
	otaRepo otaRepo.OTARepository
}

// NewService creates a new OTA service
func NewService(otaRepo otaRepo.OTARepository) Service {
	return &service{
		otaRepo: otaRepo,
	}
}

func (s *service) Create(ctx context.Context, ota *otaModel.OTA) (*otaModel.OTA, error) {
	if err := validateOTA(ota); err != nil {
		return nil, err
	}

	// Check if there's already an OTA for this app
	existingOTA, err := s.otaRepo.GetByAppID(ctx, ota.AppID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing OTA: %w", err)
	}

	if existingOTA != nil {
		return nil, errors.New("an OTA update already exists for this app")
	}

	return s.otaRepo.Create(ctx, ota)
}

func (s *service) GetByAppID(ctx context.Context, appID string) (*otaModel.OTA, error) {
	if appID == "" {
		return nil, errors.New("invalid app ID")
	}

	return s.otaRepo.GetByAppID(ctx, appID)
}

func (s *service) List(ctx context.Context, limit, page int) ([]*otaModel.OTA, error) {

	return s.otaRepo.List(ctx, limit, page)
}

func (s *service) Count(ctx context.Context) (int, error) {
	return s.otaRepo.Count(ctx)
}

func (s *service) UpdateByAppID(ctx context.Context, ota *otaModel.OTA, appID string) error {
	if appID == "" {
		return errors.New("invalid app ID")
	}

	if err := validateOTA(ota); err != nil {
		return err
	}

	// Get existing OTA to check if it exists
	existing, err := s.otaRepo.GetByAppID(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to get existing OTA: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("OTA for app ID %s not found", appID)
	}

	return s.otaRepo.UpdateByAppID(ctx, ota, appID)
}

func (s *service) DeleteByAppID(ctx context.Context, appID string) error {
	if appID == "" {
		return errors.New("invalid app ID")
	}

	return s.otaRepo.DeleteByAppID(ctx, appID)
}

// validateOTA validates OTA fields
func validateOTA(ota *otaModel.OTA) error {
	if ota.AppID == "" {
		return errors.New("app ID is required")
	}

	if ota.VersionName == "" {
		return errors.New("version name is required")
	}

	if ota.VersionCode <= 0 {
		return errors.New("version code must be a positive number")
	}

	if ota.URL == "" {
		return errors.New("URL is required")
	}

	return nil
}
