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

// IsPointInRadius Вычисление расстояние между двумя точками на сфере
func (i *Incident) IsPointInRadius(lat, lng float64) bool {
	radLat := math.Pi * i.Lat / 180.0
	radLng := math.Pi * i.Lng / 180.0
	radLatVerified := math.Pi * lat / 180.0
	radLngVerified := math.Pi * lng / 180.0

	deltaRadLat := radLatVerified - radLat
	deltaRadLng := radLngVerified - radLng

	halfChord := math.Pow(math.Sin(deltaRadLat/2), 2) +
		math.Cos(radLat)*math.Cos(radLatVerified)*
			math.Pow(math.Sin(deltaRadLng/2), 2)

	angularDistance := 2 * math.Atan2(math.Sqrt(halfChord), math.Sqrt(1-halfChord))
	distanceMeters := EarthRadiusMeters * angularDistance // расстояние в метрах

	return distanceMeters <= i.Radius
}
