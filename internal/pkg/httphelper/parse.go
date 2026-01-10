package httphelper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func ParseUUIDFromPath(r *http.Request, prefix string) (uuid.UUID, error) {
	idStr := strings.TrimPrefix(r.URL.Path, prefix)
	idStr = strings.Trim(idStr, "/")
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("id is required in path")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid id format")
	}

	return id, nil
}
