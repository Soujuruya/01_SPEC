package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Soujuruya/01_SPEC/internal/domain/location"
	"github.com/Soujuruya/01_SPEC/internal/pkg/errs"
	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepo struct {
	pgxPool *pgxpool.Pool
	builder squirrel.StatementBuilderType
	lg      *logger.Logger
}

func NewLocationRepo(pgxPool *pgxpool.Pool, lg *logger.Logger) *LocationRepo {
	return &LocationRepo{
		pgxPool: pgxPool,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		lg:      lg,
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
		r.lg.Error("LocationRepo.Save: error building query", "error", err, "user_id", loc.UserID)
		return err
	}

	_, err = r.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			r.lg.Warn("LocationRepo.Save: duplicate entry", "user_id", loc.UserID)
			return errs.ErrDuplicate
		}
		r.lg.Error("LocationRepo.Save: error executing query", "error", err, "user_id", loc.UserID)
		return err
	}

	r.lg.Debug("LocationRepo.Save: location saved successfully", "user_id", loc.UserID, "location_id", loc.ID)
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
		r.lg.Error("LocationRepo.ListByUser: error building query", "error", err, "user_id", userID)
		return nil, err
	}

	rows, err := r.pgxPool.Query(ctx, query, args...)
	if err != nil {
		r.lg.Error("LocationRepo.ListByUser: error executing query", "error", err, "user_id", userID)
		return nil, err
	}
	defer rows.Close()

	var locations []*location.Location
	for rows.Next() {
		loc := &location.Location{}
		if err := rows.Scan(&loc.ID, &loc.UserID, &loc.Lat, &loc.Lng, &loc.Timestamp, &loc.IsCheck, &loc.IncidentIDs); err != nil {
			r.lg.Error("LocationRepo.ListByUser: error scanning row", "error", err, "user_id", userID)
			return nil, err
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		r.lg.Error("LocationRepo.ListByUser: rows error", "error", err, "user_id", userID)
		return nil, err
	}

	if len(locations) == 0 {
		r.lg.Debug("LocationRepo.ListByUser: no locations found", "user_id", userID)
		return nil, errs.ErrNotFound
	}

	r.lg.Debug("LocationRepo.ListByUser: locations fetched successfully", "user_id", userID, "count", len(locations))
	return locations, nil
}

// CountUniqueUsers возвращает количество уникальных пользователей с указанного времени
func (r *LocationRepo) CountUniqueUsers(ctx context.Context, since time.Time) (int, error) {
	query, args, err := r.builder.
		Select("COUNT(DISTINCT user_id)").
		From("locations").
		Where(squirrel.GtOrEq{"timestamp": since}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		r.lg.Error("LocationRepo.CountUniqueUsers: error building query", "error", err)
		return 0, err
	}

	var count int
	err = r.pgxPool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		r.lg.Error("LocationRepo.CountUniqueUsers: error executing query", "error", err)
		return 0, err
	}

	r.lg.Debug("LocationRepo.CountUniqueUsers: unique user count fetched", "count", count)
	return count, nil
}
