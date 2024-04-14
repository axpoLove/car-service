// TODO: DELETE PACKAGE
package carinfo

import (
	"context"
	"database/sql"

	"car-service/internal/models"
)

type repository struct {
	db *sql.DB
}

// GetCarInfo returns car info by reg num
func (r *repository) GetCarInfo(ctx context.Context, regNum string) (*models.CarInfo, error) {
	carInfo := &models.CarInfo{}
	err := r.db.QueryRow(`SELECT * FROM cars WHERE reg_num = ?`, regNum).Scan(
		&carInfo.RegNum,
		&carInfo.Mark,
		&carInfo.Model,
		&carInfo.Year,
		&carInfo.Owner,
	)
	return carInfo, err
}

// NewRepository returns a new repository instance
func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
