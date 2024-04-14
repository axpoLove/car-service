package car

import "car-service/internal/models"

// GetCarsResponse is a get cars response structure
type GetCarsResponse struct {
	Cars  []models.CarInfo `json:"cars,omitempty"`
	Total int              `json:"total,omitempty"`
	Page  int              `json:"page,omitempty"`
}

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

type emptyResponse struct{}
