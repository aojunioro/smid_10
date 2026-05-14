package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type unitRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewUnitRepository(db *sql.DB, alias common.DBAlias) UnitRepository {
	return &unitRepository{
		db:     db,
		alias: alias,
	}
}

func (r *unitRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *unitRepository) Create(ctx context.Context, unit *SystemUnit) error {
	now := time.Now()
	query := `
		INSERT INTO system_units (name, parent_id, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		unit.Name,
		unit.ParentID,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create unit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	unit.ID = id
	unit.CreatedAt = now
	unit.UpdatedAt = now
	return nil
}

func (r *unitRepository) Update(ctx context.Context, unit *SystemUnit) error {
	query := `
		UPDATE system_units
		SET name = ?, parent_id = ?, updated_at = ?
		WHERE id = ?
	`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		unit.Name,
		unit.ParentID,
		now,
		unit.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("unit not found")
	}

	unit.UpdatedAt = now
	return nil
}

func (r *unitRepository) FindByID(ctx context.Context, id int64) (*SystemUnit, error) {
	query := `
		SELECT id, name, parent_id, created_at, updated_at
		FROM system_units
		WHERE id = ?
	`
	
	var unit SystemUnit
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&unit.ID,
		&unit.Name,
		&unit.ParentID,
		&unit.CreatedAt,
		&unit.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unit not found")
		}
		return nil, fmt.Errorf("failed to find unit: %w", err)
	}

	return &unit, nil
}

func (r *unitRepository) List(ctx context.Context, opts ListOptions) ([]SystemUnit, error) {
	query := `
		SELECT id, name, parent_id, created_at, updated_at
		FROM system_units
		ORDER BY created_at DESC
	`
	
	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list units: %w", err)
	}
	defer rows.Close()

	var units []SystemUnit
	for rows.Next() {
		var unit SystemUnit
		err := rows.Scan(
			&unit.ID,
			&unit.Name,
			&unit.ParentID,
			&unit.CreatedAt,
			&unit.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan unit: %w", err)
		}
		units = append(units, unit)
	}

	return units, nil
}

func (r *unitRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM system_units WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("unit not found")
	}

	return nil
}
