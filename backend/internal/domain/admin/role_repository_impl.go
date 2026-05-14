package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type roleRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewRoleRepository(db *sql.DB, alias common.DBAlias) RoleRepository {
	return &roleRepository{
		db:     db,
		alias: alias,
	}
}

func (r *roleRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *roleRepository) Create(ctx context.Context, role *SystemRole) error {
	now := time.Now()
	query := `
		INSERT INTO system_roles (name, created_at, updated_at)
		VALUES (?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		role.Name,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	role.ID = id
	role.CreatedAt = now
	role.UpdatedAt = now
	return nil
}

func (r *roleRepository) Update(ctx context.Context, role *SystemRole) error {
	query := `
		UPDATE system_roles
		SET name = ?, updated_at = ?
		WHERE id = ?
	`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		role.Name,
		now,
		role.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("role not found")
	}

	role.UpdatedAt = now
	return nil
}

func (r *roleRepository) FindByID(ctx context.Context, id int64) (*SystemRole, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM system_roles
		WHERE id = ?
	`
	
	var role SystemRole
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}

	return &role, nil
}

func (r *roleRepository) List(ctx context.Context, opts ListOptions) ([]SystemRole, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM system_roles
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
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []SystemRole
	for rows.Next() {
		var role SystemRole
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM system_roles WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}
