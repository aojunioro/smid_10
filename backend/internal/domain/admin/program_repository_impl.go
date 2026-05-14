package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type programRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewProgramRepository(db *sql.DB, alias common.DBAlias) ProgramRepository {
	return &programRepository{
		db:     db,
		alias: alias,
	}
}

func (r *programRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *programRepository) Create(ctx context.Context, program *SystemProgram) error {
	now := time.Now()
	query := `
		INSERT INTO system_programs (name, controller, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		program.Name,
		program.Controller,
		program.Description,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create program: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	program.ID = id
	program.CreatedAt = now
	program.UpdatedAt = now
	return nil
}

func (r *programRepository) Update(ctx context.Context, program *SystemProgram) error {
	query := `
		UPDATE system_programs
		SET name = ?, controller = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		program.Name,
		program.Controller,
		program.Description,
		now,
		program.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update program: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("program not found")
	}

	program.UpdatedAt = now
	return nil
}

func (r *programRepository) FindByID(ctx context.Context, id int64) (*SystemProgram, error) {
	query := `
		SELECT id, name, controller, description, created_at, updated_at
		FROM system_programs
		WHERE id = ?
	`
	
	var program SystemProgram
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&program.ID,
		&program.Name,
		&program.Controller,
		&program.Description,
		&program.CreatedAt,
		&program.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("program not found")
		}
		return nil, fmt.Errorf("failed to find program: %w", err)
	}

	return &program, nil
}

func (r *programRepository) List(ctx context.Context, opts ListOptions) ([]SystemProgram, error) {
	query := `
		SELECT id, name, controller, description, created_at, updated_at
		FROM system_programs
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
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}
	defer rows.Close()

	var programs []SystemProgram
	for rows.Next() {
		var program SystemProgram
		err := rows.Scan(
			&program.ID,
			&program.Name,
			&program.Controller,
			&program.Description,
			&program.CreatedAt,
			&program.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan program: %w", err)
		}
		programs = append(programs, program)
	}

	return programs, nil
}

func (r *programRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM system_programs WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete program: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("program not found")
	}

	return nil
}
