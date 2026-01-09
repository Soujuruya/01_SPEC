package httphelper

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func ParseUUIDFromQuery(r *http.Request, param string) (uuid.UUID, error) {
	idStr := r.URL.Query().Get(param)
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("%s is required", param)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s format", param)
	}
	return id, nil
}
