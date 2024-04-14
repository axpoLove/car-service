package car

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpswagger "github.com/swaggo/http-swagger/v2"

	"car-service/internal/models"
)

type service interface {
	GetCars(ctx context.Context, filter *models.CarInfo, page, limit int) ([]models.CarInfo, int, int, error)
	DeleteCar(ctx context.Context, regNum string) error
	UpdateCar(ctx context.Context, regNum string, car *models.CarInfo) error
	AddCar(ctx context.Context, regNums []string) error
}

type server struct {
	srv     *http.Server
	port    int
	router  *chi.Mux
	service service
	logger  *slog.Logger
}

// Starts the server
func (s *server) Start() error {
	return s.srv.ListenAndServe()
}

// handleGetCars godoc
// @ID get-cars
// @Summary get cars
// @Description get cars
// @Accept json
// @Produce json
// @Param Filter body GetCarsRequest true "filter for car list"
// @Success 200 {object} GetCarsResponse
// @Failure 400 string empty
// @Failure 500 string empty
// @Router /cars [get]
func (s *server) handleGetCars(w http.ResponseWriter, r *http.Request) {
	var request GetCarsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var (
		total int
		page  int
		cars  []models.CarInfo
	)
	cars, total, page, err = s.service.GetCars(
		r.Context(),
		&models.CarInfo{
			RegNum: request.RegNum,
			Mark:   request.Mark,
			Model:  request.Model,
			Year:   request.Year,
			Owner:  request.Owner,
		},
		request.Page,
		request.Limit,
	)
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&GetCarsResponse{
		Cars:  cars,
		Total: total,
		Page:  page,
	})
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUpdateCar godoc
// @ID update-car
// @Summary update car
// @Description update car
// @Accept json
// @Produce json
// @Param CarInfo body UpdateCarRequest true "car info for update"
// @Success 200 {object} emptyResponse
// @Failure 400 string empty
// @Failure 500 string empty
// @Router /car [patch]
func (s *server) handleUpdateCar(w http.ResponseWriter, r *http.Request) {
	var request UpdateCarRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = request.Validate()
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.service.UpdateCar(r.Context(), request.RegNum, &models.CarInfo{
		Mark:  request.Mark,
		Model: request.Model,
		Year:  request.Year,
		Owner: request.Owner,
	})
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&emptyResponse{})
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddCar godoc
// @ID add-cars
// @Summary add cars
// @Description add cars
// @Accept json
// @Produce json
// @Param regNums body AddCarRequest true "reg num list"
// @Success 200 {object} emptyResponse
// @Failure 400 string empty
// @Failure 500 string empty
// @Router /cars [post]
func (s *server) handleAddCar(w http.ResponseWriter, r *http.Request) {
	var request AddCarRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = request.Validate()
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.service.AddCar(r.Context(), request.RegNums)
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&emptyResponse{})
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleDeleteCar godoc
// @ID delete-car
// @Summary delete car
// @Description delete car
// @Accept json
// @Produce json
// @Param regNum body DeleteCarRequest true "reg num"
// @Success 200 {object} emptyResponse
// @Failure 400 string empty
// @Failure 500 string empty
// @Router /car [delete]
func (s *server) handleDeleteCar(w http.ResponseWriter, r *http.Request) {
	var request DeleteCarRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = request.Validate()
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.service.DeleteCar(r.Context(), request.RegNum)
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&emptyResponse{})
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) apiError(w http.ResponseWriter, msg string, status int) {
	s.logger.Error("unexpected error", "error", msg)
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(&errorResponse{Error: msg})
	if err != nil {
		s.logger.Error("failed to encode response error", "errot", err)
	}
}

// NewServer returns a new server instance
func NewServer(
	port int,
	readHeaderTimeout time.Duration,
	service service,
	logger *slog.Logger,
) *server {
	router := chi.NewRouter()
	s := &server{
		port:    port,
		router:  router,
		service: service,
		logger:  logger,
		srv: &http.Server{
			Addr:              net.JoinHostPort("", strconv.Itoa(port)),
			Handler:           router,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
	router.Use(middleware.Logger)
	router.Get("/cars", s.handleGetCars)
	router.Delete("/car", s.handleDeleteCar)
	router.Patch("/car", s.handleUpdateCar)
	router.Post("/cars", s.handleAddCar)
	router.Get("/swagger/*", httpswagger.Handler())
	return s
}
