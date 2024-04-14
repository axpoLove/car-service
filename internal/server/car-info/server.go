package carinfo

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
	GetCarInfo(ctx context.Context, regNum string) (*models.CarInfo, error)
}

type server struct {
	srv     *http.Server
	port    int
	router  *chi.Mux
	logger  *slog.Logger
	service service
}

// Starts the server
func (s *server) Start() error {
	return s.srv.ListenAndServe()
}

// handleGerCarInfo godoc
// @ID get-car-info
// @Summary get car info
// @Description get car info
// @Produce json
// @Param regNum query string true "car reg number"
// @Success 200 {object} models.CarInfo
// @Failure 400 string empty
// @Failure 500 string empty
// @Router /info [get]
func (s *server) handleGetCarInfo(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	regNum := queryValues.Get("regNum")
	getCarInfoRequest := GetCarInfoRequest{RegNum: regNum}
	err := getCarInfoRequest.Validate()
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	carInfo, err := s.service.GetCarInfo(r.Context(), regNum)
	if err != nil {
		s.apiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&carInfo)
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
	server := &server{
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
	router.Get("/info", server.handleGetCarInfo)
	router.Get("/swagger/*", httpswagger.Handler())
	return server
}
