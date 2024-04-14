package carinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"car-service/internal/models"
)

const (
	infoEndpoint = "/info"
)

type client struct {
	baseURL string
	cli     *http.Client
}

// GetCarInfo returns car info by reg num
func (c *client) GetCarInfo(ctx context.Context, regNum string) (info *models.CarInfo, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s%s", c.baseURL, infoEndpoint),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new http request: %w", err)
	}
	q := req.URL.Query()
	q.Add("regNum", regNum)
	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	resp, err = c.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer resp.Body.Close()

	info = &models.CarInfo{}
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return nil, fmt.Errorf("failed to decode car info request: %w", err)
	}
	return
}

// NewClient returns a new car info client instance
func NewClient(baseURL string, cli *http.Client) *client {
	return &client{
		cli:     cli,
		baseURL: baseURL,
	}
}
