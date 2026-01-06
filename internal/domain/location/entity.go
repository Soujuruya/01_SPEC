package location

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID          uuid.UUID   `json:"id"`           // уникальный идентификатор проверки
	UserID      uuid.UUID   `json:"user_id"`      // пользователь, который прислал координаты
	Lat         float64     `json:"lat"`          // широта пользователя
	Lng         float64     `json:"lng"`          // долгота пользователя
	Timestamp   time.Time   `json:"timestamp"`    // когда пришли координаты
	IsCheck     bool        `json:"is_check"`     // попали ли координаты в зону инцидента
	IncidentIDs []uuid.UUID `json:"incident_ids"` // список инцидентов, которые попали
}

func (l *Location) HasIncidents() bool {
	return len(l.IncidentIDs) > 0
}
