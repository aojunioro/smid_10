package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type groupRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewGroupRepository(db *sql.DB, alias common.DBAlias) GroupRepository {
	return &groupRepository{
		db:     db,
		alias: alias,
	}
}

func (r *groupRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *groupRepository) Create(ctx context.Context, group *SystemGroup) error {
	now := time.Now()
	query := `
		INSERT INTO system_groups (name, frontpage_id, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		group.Name,
		group.FrontpageID,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	group.ID = id
	group.CreatedAt = now
	group.UpdatedAt = now
	return nil
}

func (r *groupRepository) Update(ctx context.Context, group *SystemGroup) error {
	query := `
		UPDATE system_groups
		SET name = ?, frontpage_id = ?, updated_at = ?
		WHERE id = ?
	`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		group.Name,
		group.FrontpageID,
		now,
		group.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("group not found")
	}

	group.UpdatedAt = now
	return nil
}

func (r *groupRepository) FindByID(ctx context.Context, id int64) (*SystemGroup, error) {
	query := `
		SELECT id, name, frontpage_id, created_at, updated_at
		FROM system_groups
		WHERE id = ?
	`
	
	var group SystemGroup
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID,
		&group.Name,
		&group.FrontpageID,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to find group: %w", err)
	}

	return &group, nil
}

func (r *groupRepository) List(ctx context.Context, opts ListOptions) ([]SystemGroup, error) {
	query := `
		SELECT id, name, frontpage_id, created_at, updated_at
		FROM system_groups
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
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	defer rows.Close()

	var groups []SystemGroup
	for rows.Next() {
		var group SystemGroup
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.FrontpageID,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (r *groupRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM system_groups WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}
