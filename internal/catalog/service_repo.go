package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgServiceRepo struct {
	db *sql.DB
}

var _ ServiceRepository = (*pgServiceRepo)(nil)

func NewPgServiceRepo(db *sql.DB) *pgServiceRepo {
	return &pgServiceRepo{db: db}
}

func (r *pgServiceRepo) Insert(ctx context.Context, s *Service) error {
	const q = `
		INSERT INTO services
			(id, provider_id, category_id, titulo, descripcion, precio_base,
			 lat, lng, radio_km, is_active, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, q,
		s.ID, s.ProviderID, s.CategoryID,
		s.Titulo, s.Descripcion, s.PrecioBase,
		toNullFloat(s.Lat), toNullFloat(s.Lng), toNullFloat(s.RadioKm),
		s.IsActive, s.CreatedAt)
	if err != nil {
		return mapServiceFKError(err, "insert service")
	}
	return nil
}

func (r *pgServiceRepo) FindAll(ctx context.Context, filter ServiceFilter) ([]*Service, error) {
	const base = `
		SELECT id, provider_id, category_id, titulo, COALESCE(descripcion, ''),
		       precio_base, lat, lng, radio_km, is_active, created_at
		FROM services
	`
	var (
		conds []string
		args  []any
	)
	if filter.CategoryID != nil {
		args = append(args, *filter.CategoryID)
		conds = append(conds, fmt.Sprintf("category_id = $%d", len(args)))
	}
	if filter.IsActive != nil {
		args = append(args, *filter.IsActive)
		conds = append(conds, fmt.Sprintf("is_active = $%d", len(args)))
	}
	q := base
	if len(conds) > 0 {
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("find all services: %w", err)
	}
	defer rows.Close()

	var out []*Service
	for rows.Next() {
		s, err := scanService(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate services: %w", err)
	}
	return out, nil
}

func (r *pgServiceRepo) FindByID(ctx context.Context, id uuid.UUID) (*Service, error) {
	const q = `
		SELECT id, provider_id, category_id, titulo, COALESCE(descripcion, ''),
		       precio_base, lat, lng, radio_km, is_active, created_at
		FROM services
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, q, id)
	s, err := scanService(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrServiceNotFound
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *pgServiceRepo) Update(ctx context.Context, s *Service) error {
	const q = `
		UPDATE services
		SET provider_id = $2, category_id = $3, titulo = $4, descripcion = $5,
		    precio_base = $6, lat = $7, lng = $8, radio_km = $9, is_active = $10
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, q,
		s.ID, s.ProviderID, s.CategoryID,
		s.Titulo, s.Descripcion, s.PrecioBase,
		toNullFloat(s.Lat), toNullFloat(s.Lng), toNullFloat(s.RadioKm),
		s.IsActive)
	if err != nil {
		return mapServiceFKError(err, "update service")
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update service rows affected: %w", err)
	}
	if n == 0 {
		return ErrServiceNotFound
	}
	return nil
}

func (r *pgServiceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM services WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete service: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete service rows affected: %w", err)
	}
	if n == 0 {
		return ErrServiceNotFound
	}
	return nil
}

// rowScanner abstrae *sql.Row y *sql.Rows para reusar el scan en FindByID y FindAll.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanService(r rowScanner) (*Service, error) {
	s := &Service{}
	var lat, lng, radio sql.NullFloat64
	if err := r.Scan(
		&s.ID, &s.ProviderID, &s.CategoryID,
		&s.Titulo, &s.Descripcion, &s.PrecioBase,
		&lat, &lng, &radio, &s.IsActive, &s.CreatedAt,
	); err != nil {
		return nil, err
	}
	if lat.Valid {
		v := lat.Float64
		s.Lat = &v
	}
	if lng.Valid {
		v := lng.Float64
		s.Lng = &v
	}
	if radio.Valid {
		v := radio.Float64
		s.RadioKm = &v
	}
	return s, nil
}

func toNullFloat(v *float64) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *v, Valid: true}
}

func mapServiceFKError(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
		switch {
		case strings.Contains(pgErr.ConstraintName, "provider"):
			return ErrProviderNotFound
		case strings.Contains(pgErr.ConstraintName, "category"):
			return ErrCategoryNotFound
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}
