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

type pgCategoryRepo struct {
	db *sql.DB
}

var _ CategoryRepository = (*pgCategoryRepo)(nil)

func NewPgCategoryRepo(db *sql.DB) *pgCategoryRepo {
	return &pgCategoryRepo{db: db}
}

func (r *pgCategoryRepo) Insert(ctx context.Context, c *Category) error {
	const q = `
		INSERT INTO categories (id, nombre, slug, parent_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	parent := uuid.NullUUID{}
	if c.ParentID != nil {
		parent = uuid.NullUUID{UUID: *c.ParentID, Valid: true}
	}
	_, err := r.db.ExecContext(ctx, q, c.ID, c.Nombre, c.Slug, parent, c.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrDuplicateSlug
		}
		return fmt.Errorf("insert category: %w", err)
	}
	return nil
}

func (r *pgCategoryRepo) FindAll(ctx context.Context) ([]*Category, error) {
	const q = `
		SELECT id, nombre, slug, parent_id, created_at
		FROM categories
		ORDER BY nombre
	`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("find all categories: %w", err)
	}
	defer rows.Close()

	var out []*Category
	for rows.Next() {
		c := &Category{}
		var parent uuid.NullUUID
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Slug, &parent, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		if parent.Valid {
			id := parent.UUID
			c.ParentID = &id
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}
	return out, nil
}

func (r *pgCategoryRepo) FindByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	const q = `
		SELECT id, nombre, slug, parent_id, created_at
		FROM categories
		WHERE id = $1
	`
	c := &Category{}
	var parent uuid.NullUUID
	err := r.db.QueryRowContext(ctx, q, id).
		Scan(&c.ID, &c.Nombre, &c.Slug, &parent, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCategoryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find category by id: %w", err)
	}
	if parent.Valid {
		pid := parent.UUID
		c.ParentID = &pid
	}
	return c, nil
}

func (r *pgCategoryRepo) Update(ctx context.Context, c *Category) error {
	const q = `
		UPDATE categories
		SET nombre = $2, slug = $3, parent_id = $4
		WHERE id = $1
	`
	parent := uuid.NullUUID{}
	if c.ParentID != nil {
		parent = uuid.NullUUID{UUID: *c.ParentID, Valid: true}
	}
	res, err := r.db.ExecContext(ctx, q, c.ID, c.Nombre, c.Slug, parent)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrDuplicateSlug
		}
		return fmt.Errorf("update category: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update category rows affected: %w", err)
	}
	if n == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

func (r *pgCategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM categories WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			(pgErr.Code == pgerrcode.ForeignKeyViolation || pgErr.Code == pgerrcode.RestrictViolation) {
			return ErrCategoryHasServices
		}
		return fmt.Errorf("delete category: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete category rows affected: %w", err)
	}
	if n == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
