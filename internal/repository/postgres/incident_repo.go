package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Soujuruya/01_SPEC/internal/domain/incident"
	"github.com/Soujuruya/01_SPEC/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepo struct {
	pgxPool *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewIncidentRepo(pgxPool *pgxpool.Pool) *IncidentRepo {
	return &IncidentRepo{
		pgxPool: pgxPool,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *IncidentRepo) CountActiveIncidents(ctx context.Context) (int, error) {
	query, args, err := r.builder.
		Select("COUNT(*)").
		From("incidents").
		Where(squirrel.Eq{"is_active": true}).
		ToSql()
	if err != nil {
		return 0, err
	}

	var total int
	err = r.pgxPool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *IncidentRepo) GetActiveIncidents(ctx context.Context) ([]*incident.Incident, error) {
	query, args, err := r.builder.
		Select("id", "title", "lat", "lng", "radius", "is_active", "created_at", "updated_at").
		From("incidents").
		Where(squirrel.Eq{"is_active": true}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pgxPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*incident.Incident
	for rows.Next() {
		i := &incident.Incident{}
		if err := rows.Scan(
			&i.ID, &i.Title, &i.Lat, &i.Lng, &i.Radius,
			&i.IsActive, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}

func (r *IncidentRepo) Create(ctx context.Context, inc *incident.Incident) error {
	if inc.ID == uuid.Nil {
		inc.ID = uuid.New()
	}
	if inc.CreatedAt.IsZero() {
		inc.CreatedAt = time.Now()
	}
	inc.UpdatedAt = inc.CreatedAt

	query, args, err := r.builder.
		Insert("incidents").
		Columns("id", "title", "lat", "lng", "radius", "is_active", "created_at", "updated_at").
		Values(inc.ID, inc.Title, inc.Lat, inc.Lng, inc.Radius, inc.IsActive, inc.CreatedAt, inc.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errs.ErrDuplicate
		}
		return err
	}

	return nil
}

func (r *IncidentRepo) GetByID(ctx context.Context, id uuid.UUID) (*incident.Incident, error) {
	query, args, err := r.builder.
		Select("id", "title", "lat", "lng", "radius", "is_active", "created_at", "updated_at").
		From("incidents").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.pgxPool.QueryRow(ctx, query, args...)
	inc := &incident.Incident{}
	err = row.Scan(&inc.ID, &inc.Title, &inc.Lat, &inc.Lng, &inc.Radius, &inc.IsActive, &inc.CreatedAt, &inc.UpdatedAt)
	if err != nil {
		return nil, errs.ErrNotFound
	}
	return inc, nil
}

func (r *IncidentRepo) Update(ctx context.Context, inc *incident.Incident) error {
	inc.UpdatedAt = time.Now()

	query, args, err := r.builder.
		Update("incidents").
		Set("title", inc.Title).
		Set("lat", inc.Lat).
		Set("lng", inc.Lng).
		Set("radius", inc.Radius).
		Set("is_active", inc.IsActive).
		Set("updated_at", inc.UpdatedAt).
		Where(squirrel.Eq{"id": inc.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *IncidentRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.builder.
		Update("incidents").
		Set("is_active", false).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *IncidentRepo) List(ctx context.Context, offset, limit int) ([]*incident.Incident, error) {
	query, args, err := r.builder.
		Select("id", "title", "lat", "lng", "radius", "is_active", "created_at", "updated_at").
		From("incidents").
		OrderBy("created_at DESC").
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pgxPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*incident.Incident
	for rows.Next() {
		i := &incident.Incident{}
		if err := rows.Scan(&i.ID, &i.Title, &i.Lat, &i.Lng, &i.Radius, &i.IsActive, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		incidents = append(incidents, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}

func (r *IncidentRepo) ListWithTotal(ctx context.Context, offset, limit int) ([]*incident.Incident, int, error) {
	incs, err := r.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := r.CountActiveIncidents(ctx)
	if err != nil {
		return nil, 0, err
	}

	return incs, total, nil
}
