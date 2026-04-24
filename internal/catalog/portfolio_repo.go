package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgPortfolioRepo struct {
	db *sql.DB
}

var _ PortfolioRepository = (*pgPortfolioRepo)(nil)

func NewPgPortfolioRepo(db *sql.DB) *pgPortfolioRepo {
	return &pgPortfolioRepo{db: db}
}

func (r *pgPortfolioRepo) Insert(ctx context.Context, p *PortfolioItem) error {
	const q = `
		INSERT INTO portfolio_items (id, service_id, storage_url, titulo, orden)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, q, p.ID, p.ServiceID, p.StorageURL, p.Titulo, p.Orden)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return ErrServiceNotFound
		}
		return fmt.Errorf("insert portfolio item: %w", err)
	}
	return nil
}

func (r *pgPortfolioRepo) FindByServiceID(ctx context.Context, serviceID uuid.UUID) ([]*PortfolioItem, error) {
	const q = `
		SELECT id, service_id, storage_url, COALESCE(titulo, ''), orden
		FROM portfolio_items
		WHERE service_id = $1
		ORDER BY orden, id
	`
	rows, err := r.db.QueryContext(ctx, q, serviceID)
	if err != nil {
		return nil, fmt.Errorf("find portfolio items: %w", err)
	}
	defer rows.Close()

	var out []*PortfolioItem
	for rows.Next() {
		p := &PortfolioItem{}
		if err := rows.Scan(&p.ID, &p.ServiceID, &p.StorageURL, &p.Titulo, &p.Orden); err != nil {
			return nil, fmt.Errorf("scan portfolio item: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate portfolio items: %w", err)
	}
	return out, nil
}

func (r *pgPortfolioRepo) FindByID(ctx context.Context, id uuid.UUID) (*PortfolioItem, error) {
	const q = `
		SELECT id, service_id, storage_url, COALESCE(titulo, ''), orden
		FROM portfolio_items
		WHERE id = $1
	`
	p := &PortfolioItem{}
	err := r.db.QueryRowContext(ctx, q, id).
		Scan(&p.ID, &p.ServiceID, &p.StorageURL, &p.Titulo, &p.Orden)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPortfolioItemNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find portfolio item by id: %w", err)
	}
	return p, nil
}

func (r *pgPortfolioRepo) Update(ctx context.Context, p *PortfolioItem) error {
	const q = `
		UPDATE portfolio_items
		SET storage_url = $2, titulo = $3, orden = $4
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, q, p.ID, p.StorageURL, p.Titulo, p.Orden)
	if err != nil {
		return fmt.Errorf("update portfolio item: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update portfolio item rows affected: %w", err)
	}
	if n == 0 {
		return ErrPortfolioItemNotFound
	}
	return nil
}

func (r *pgPortfolioRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM portfolio_items WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete portfolio item: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete portfolio item rows affected: %w", err)
	}
	if n == 0 {
		return ErrPortfolioItemNotFound
	}
	return nil
}
