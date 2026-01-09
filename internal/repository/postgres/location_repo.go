package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Soujuruya/01_SPEC/internal/domain/location"
	"github.com/Soujuruya/01_SPEC/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepo struct {
	pgxPool *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewLocationRepo(pgxPool *pgxpool.Pool) *LocationRepo {
	return &LocationRepo{
		pgxPool: pgxPool,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Save сохраняет новую проверку локации
func (r *LocationRepo) Save(ctx context.Context, loc *location.Location) error {
	if loc.ID == uuid.Nil {
		loc.ID = uuid.New()
	}
	if loc.Timestamp.IsZero() {
		loc.Timestamp = time.Now()
	}

	query, args, err := r.builder.
		Insert("locations").
		Columns("id", "user_id", "lat", "lng", "timestamp", "is_check", "incident_ids").
		Values(loc.ID, loc.UserID, loc.Lat, loc.Lng, loc.Timestamp, loc.IsCheck, loc.IncidentIDs).
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

// ListByUser возвращает последние проверки пользователя
func (r *LocationRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*location.Location, error) {
	query, args, err := r.builder.
		Select("id", "user_id", "lat", "lng", "timestamp", "is_check", "incident_ids").
		From("locations").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("timestamp DESC").
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

	var locations []*location.Location
	for rows.Next() {
		loc := &location.Location{}
		if err := rows.Scan(&loc.ID, &loc.UserID, &loc.Lat, &loc.Lng, &loc.Timestamp, &loc.IsCheck, &loc.IncidentIDs); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, errs.ErrNotFound
	}

	return locations, nil
}

func (r *LocationRepo) CountUniqueUsers(ctx context.Context, since time.Time) (int, error) {
	query, args, err := r.builder.
		Select("COUNT(DISTINCT user_id)").
		From("locations").
		Where(squirrel.GtOrEq{"timestamp": since}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = r.pgxPool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
