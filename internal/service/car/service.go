package car

import (
	"context"
	"fmt"

	"car-service/internal/models"
)

const (
	defaultPage  = 1
	defaultLimit = 10
)

type carInfoClient interface {
	GetCarInfo(ctx context.Context, regNum string) (info *models.CarInfo, err error)
}

type repository interface {
	GetCars(ctx context.Context, filter *models.CarInfo, page, limit int) (cars []models.CarInfo, total int, err error)
	DeleteCar(ctx context.Context, regNum string) error
	UpdateCar(ctx context.Context, regNum string, car *models.CarInfo) error
	AddCar(ctx context.Context, carInfo *models.CarInfo) error
}

type service struct {
	repository    repository
	carInfoClient carInfoClient
}

// GetCars returns cars list by filter
func (s *service) GetCars(
	ctx context.Context,
	filter *models.CarInfo,
	page, limit int,
) (
	cars []models.CarInfo,
	total, outputPage int,
	err error,
) {
	if page < 1 {
		page = defaultPage
	}
	outputPage = page
	if limit < 1 {
		limit = defaultLimit
	}
	cars, total, err = s.repository.GetCars(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get cars: %w", err)
	}
	return cars, total, outputPage, nil
}

// DeleteCar deletes a car by reg num
func (s *service) DeleteCar(ctx context.Context, regNum string) error {
	err := s.repository.DeleteCar(ctx, regNum)
	if err != nil {
		return fmt.Errorf("failed go delete car: %w", err)
	}
	return nil
}

// UpdateCar updates a car
func (s *service) UpdateCar(ctx context.Context, regNum string, car *models.CarInfo) error {
	err := s.repository.UpdateCar(ctx, regNum, car)
	if err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}
	return nil
}

// AddCar adds cars by reg num list
func (s *service) AddCar(ctx context.Context, regNums []string) error {
	for _, regNum := range regNums {
		carInfo, err := s.carInfoClient.GetCarInfo(ctx, regNum)
		if err != nil {
			return fmt.Errorf("failed to get car info: %w", err)
		}
		err = s.repository.AddCar(ctx, carInfo)
		if err != nil {
			return fmt.Errorf("failed to add car to database: %w", err)
		}
	}
	return nil
}

// NewService returns a new service instance
func NewService(repository repository, carInfoClient carInfoClient) *service {
	return &service{
		repository:    repository,
		carInfoClient: carInfoClient,
	}
}
