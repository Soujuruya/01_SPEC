package incident

import (
	"math"
	"time"

	"github.com/google/uuid"
)

const EarthRadiusMeters = 6371000

type Incident struct {
	ID        uuid.UUID `json:"id"`         // идентификатор инцидента
	Title     string    `json:"title"`      // описание инцидента
	Lat       float64   `json:"lat"`        // широта зоны инцидента
	Lng       float64   `json:"lng"`        // долгота зоны зоны инцидента
	Radius    float64   `json:"radius"`     // радиус зоны инцидента
	IsActive  bool      `json:"is_active"`  // активность
	CreatedAt time.Time `json:"created_at"` // дата появления
	UpdatedAt time.Time `json:"-"`          // дата изменения информации об инциденте
}

func NewIncident(title string, lat, lng, radius float64, isActive bool) *Incident {
	return &Incident{
		ID:        uuid.New(),
		Title:     title,
		Lat:       lat,
		Lng:       lng,
		Radius:    radius,
		IsActive:  isActive,
		CreatedAt: time.Now(),
	}
}

// IsPointInRadius Вычисление расстояние между двумя точками на сфере
func (i *Incident) IsPointInRadius(lat, lng float64) bool {
	//переводим градусы в радианы
	lat1 := i.Lat * math.Pi / 180.0
	lng1 := i.Lng * math.Pi / 180.0
	lat2 := lat * math.Pi / 180.0
	lng2 := lng * math.Pi / 180.0

	dLat := lat2 - lat1
	dLng := lng2 - lng1
	//формулы
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distanceMeters := EarthRadiusMeters * c

	return distanceMeters <= i.Radius
}
