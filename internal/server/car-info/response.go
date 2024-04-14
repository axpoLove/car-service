package carinfo

import "car-service/internal/models"

// CarInfoResponse is a car info response
type CarInfoResponse []models.CarInfo

type errorResponse struct {
	Error string `json:"error,omitempty"`
}
