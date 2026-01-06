package incident

import (
	"context"

	"github.com/google/uuid"
)

type IncidentRepository interface {
	Create(ctx context.Context, inc *Incident) error
	GetByID(ctx context.Context, id uuid.UUID) (*Incident, error)
	Update(ctx context.Context, inc *Incident) error
	Deactivate(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Incident, error)
	GetActiveIncidents(ctx context.Context) ([]*Incident, error)
}
