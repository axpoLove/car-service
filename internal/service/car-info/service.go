package carinfo

import (
	"context"

	"github.com/jaswdr/faker/v2"

	"car-service/internal/models"
)

type service struct {
	fake faker.Faker
}

// GetCarInfo returns car info by reg num
func (s *service) GetCarInfo(ctx context.Context, regNum string) (*models.CarInfo, error) {
	output := &models.CarInfo{
		RegNum: regNum,
		Mark:   s.fake.Car().Maker(),
		Model:  s.fake.Car().Model(),
		Year:   s.fake.Time().Year(),
		Owner: models.Person{
			Name:       s.fake.Person().FirstName(),
			Surname:    s.fake.Person().LastName(),
			Patronymic: s.fake.Person().FirstNameMale(),
		},
	}
	return output, nil
}

// NewService returns a new service instance
func NewService(fake faker.Faker) *service {
	return &service{
		fake: fake,
	}
}
