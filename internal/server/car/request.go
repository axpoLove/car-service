package car

import (
	"fmt"

	"car-service/internal/models"
)

// GetCarsRequest is a get cars request structure
type GetCarsRequest struct {
	RegNum string        `json:"reg_num,omitempty"`
	Mark   string        `json:"mark,omitempty"`
	Model  string        `json:"model,omitempty"`
	Year   int           `json:"year,omitempty"`
	Owner  models.Person `json:"owner,omitempty"`
	Page   int           `json:"page,omitempty"`
	Limit  int           `json:"limit,omitempty"`
}

// UpdateCarRequest is an update car request structure
type UpdateCarRequest struct {
	RegNum string        `json:"reg_num,omitempty"`
	Mark   string        `json:"mark,omitempty"`
	Model  string        `json:"model,omitempty"`
	Year   int           `json:"year,omitempty"`
	Owner  models.Person `json:"owner,omitempty"`
}

// Validate validates the request
func (r *UpdateCarRequest) Validate() error {
	if r.Model == "" {
		return fmt.Errorf("invalid model")
	}
	if r.Year == 0 {
		return fmt.Errorf("invalid year")
	}
	if r.Owner.Name == "" {
		return fmt.Errorf("invalid owner name")
	}
	if r.Owner.Surname == "" {
		return fmt.Errorf("invalid owner surname")
	}
	if r.Owner.Patronymic == "" {
		return fmt.Errorf("invalid owner patronymic")
	}
	return nil
}

// AddCarRequest is an add car request structure
type AddCarRequest struct {
	RegNums []string `json:"reg_nums,omitempty"`
}

// Validate validates the request
func (r *AddCarRequest) Validate() error {
	if len(r.RegNums) == 0 {
		return fmt.Errorf("empty reg num list")
	}
	return nil
}

// DeleteCarRequest is a delete car request structure
type DeleteCarRequest struct {
	RegNum string `json:"reg_num,omitempty"`
}

// Validate validates the request
func (r *DeleteCarRequest) Validate() error {
	if r.RegNum == "" {
		return fmt.Errorf("invalid reg num")
	}
	return nil
}
