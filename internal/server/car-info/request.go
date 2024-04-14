package carinfo

import (
	"fmt"
)

// GetCarInfoRequest is a get car info request
type GetCarInfoRequest struct {
	RegNum string `json:"reg_num,omitempty"`
}

// Validate validates the request
func (c *GetCarInfoRequest) Validate() error {
	if c.RegNum == "" {
		return fmt.Errorf("invalid Regnum")
	}
	return nil
}
