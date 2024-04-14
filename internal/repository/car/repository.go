package car

import (
	"context"

	"github.com/gocraft/dbr/v2"

	"car-service/internal/models"
)

type repository struct {
	db *dbr.Connection
}

// UpdateCar updates a car
func (r *repository) UpdateCar(ctx context.Context, regNum string, car *models.CarInfo) (err error) {
	session := r.db.NewSession(nil)
	stmt := session.Update("cars")
	if car.RegNum != "" {
		stmt = stmt.Set("reg_num", car.RegNum)
	}
	if car.Mark != "" {
		stmt = stmt.Set("mark", car.Mark)
	}
	if car.Model != "" {
		stmt = stmt.Set("model", car.Model)
	}
	if car.Year != 0 {
		stmt = stmt.Set("year", car.Year)
	}
	if car.Owner.Name != "" {
		stmt = stmt.Set("owner_name", car.Owner.Name)
	}
	if car.Owner.Surname != "" {
		stmt = stmt.Set("owner_surname", car.Owner.Surname)
	}
	if car.Owner.Surname != "" {
		stmt = stmt.Set("owner_patronymic", car.Owner.Patronymic)
	}
	_, err = stmt.ExecContext(ctx)
	return
}

// AddCar adds a car
func (r *repository) AddCar(
	ctx context.Context,
	car *models.CarInfo,
) (err error) {
	session := r.db.NewSession(nil)

	_, err = session.InsertBySql(
		`INSERT INTO cars (
			reg_num,
			mark,
			model,
			year,
			owner_name,
			owner_surname,
			owner_patronymic
			)
		VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT (reg_num) DO NOTHING`,
		car.RegNum,
		car.Mark,
		car.Model,
		car.Year,
		car.Owner.Name,
		car.Owner.Surname,
		car.Owner.Patronymic,
	).ExecContext(ctx)
	return err
}

// GetCars returns cars list by filter
func (r *repository) GetCars(
	ctx context.Context,
	filter *models.CarInfo,
	page, limit int,
) (
	cars []models.CarInfo,
	total int,
	err error,
) {
	session := r.db.NewSession(nil)
	stmt := session.Select("count(*)").From("cars")
	stmt = r.buildGetCarsStatement(stmt, filter)
	err = stmt.LoadOneContext(ctx, &total)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	stmt = session.Select("*").From("cars")
	stmt = r.buildGetCarsStatement(stmt, filter)
	stmt = stmt.Limit(uint64(limit)).Offset(uint64(offset))
	cars = make([]models.CarInfo, 0)
	_, err = stmt.LoadContext(ctx, &cars)
	if err != nil {
		return nil, 0, err
	}
	return cars, total, nil
}

// DeleteCar deletes a car by reg num
func (r *repository) DeleteCar(ctx context.Context, regNum string) (err error) {
	session := r.db.NewSession(nil)
	_, err = session.DeleteFrom("cars").Where("reg_num = $1", regNum).ExecContext(ctx)
	return err
}

func (r *repository) buildGetCarsStatement(stmt *dbr.SelectStmt, filter *models.CarInfo) *dbr.SelectStmt {
	if filter.RegNum != "" {
		stmt = stmt.Where("reg_num = ?", filter.RegNum)
	}
	if filter.Mark != "" {
		stmt = stmt.Where("mark = ?", filter.Mark)
	}
	if filter.Model != "" {
		stmt = stmt.Where("model = ?", filter.Model)
	}
	if filter.Year != 0 {
		stmt = stmt.Where("year = ?", filter.Year)
	}
	if filter.Owner.Name != "" {
		stmt = stmt.Where("owner_name = ?", filter.Owner.Name)
	}
	if filter.Owner.Surname != "" {
		stmt = stmt.Where("owner_surname = ?", filter.Owner.Surname)
	}
	if filter.Owner.Patronymic != "" {
		stmt = stmt.Where("owner_patronymic = ?", filter.Owner.Patronymic)
	}
	return stmt
}

// NewRepository returns a new repository instance
func NewRepository(db *dbr.Connection) *repository {
	return &repository{
		db: db,
	}
}
