package incident

import (
	"errors"
	"fmt"
)

func ValidateCreateIncident(req *CreateIncidentRequest) error {
	if req.Title == "" {
		return errors.New("title cannot be empty")
	}
	if req.Lat < -90 || req.Lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %f", req.Lat)
	}
	if req.Lng < -180 || req.Lng > 180 {

		return fmt.Errorf("longitude must be between -180 and 180, got %f", req.Lng)
	}
	if req.Radius <= 0 {
		return fmt.Errorf("radius must be non-negative or null, got %f", req.Radius)
	}
	return nil
}

func ValidateUpdateIncident(req *UpdateIncidentRequest) error {
	if req.Title != nil && *req.Title == "" {
		return errors.New("title cannot be empty")
	}
	if req.Lat != nil && (*req.Lat < -90 || *req.Lat > 90) {
		return fmt.Errorf("latitude must be between -90 and 90, got %f", *req.Lat)
	}
	if req.Lng != nil && (*req.Lng < -180 || *req.Lng > 180) {
		return fmt.Errorf("longitude must be between -180 and 180, got %f", *req.Lng)
	}
	if req.Radius != nil && *req.Radius < 0 {
		return fmt.Errorf("radius must be non-negative, got %f", *req.Radius)
	}
	return nil
}

func ValidateLimitOffset(limit, offset int) error {
	if limit < 1 || limit > 200 {
		return fmt.Errorf("limit must be between 1 and 200, got %d", limit)
	}
	if offset < 0 {
		return fmt.Errorf("offset must be >= 0, got %d", offset)
	}
	return nil
}
