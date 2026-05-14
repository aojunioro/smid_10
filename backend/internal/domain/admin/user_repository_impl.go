package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type userRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewUserRepository(db *sql.DB, alias common.DBAlias) UserRepository {
	return &userRepository{
		db:     db,
		alias: alias,
	}
}

func (r *userRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *userRepository) Create(ctx context.Context, user *SystemUser) error {
	now := time.Now()
	query := `
		INSERT INTO system_users (login, name, email, password_hash, system_unit_id, active, frontpage_id, equipe_id, phone, address, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		user.Login,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.SystemUnitID,
		user.Active,
		user.FrontpageID,
		user.EquipeID,
		user.Phone,
		user.Address,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *SystemUser) error {
	query := `
		UPDATE system_users
		SET login = ?, name = ?, email = ?, password_hash = ?, system_unit_id = ?, active = ?, frontpage_id = ?, equipe_id = ?, phone = ?, address = ?, updated_at = ?
		WHERE id = ?
	`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		user.Login,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.SystemUnitID,
		user.Active,
		user.FrontpageID,
		user.EquipeID,
		user.Phone,
		user.Address,
		now,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	user.UpdatedAt = now
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*SystemUser, error) {
	query := `
		SELECT id, login, name, email, password_hash, system_unit_id, active, frontpage_id, equipe_id, phone, address, created_at, updated_at
		FROM system_users
		WHERE id = ?
	`
	
	var user SystemUser
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.SystemUnitID,
		&user.Active,
		&user.FrontpageID,
		&user.EquipeID,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByLogin(ctx context.Context, login string) (*SystemUser, error) {
	query := `
		SELECT id, login, name, email, password_hash, system_unit_id, active, frontpage_id, equipe_id, phone, address, created_at, updated_at
		FROM system_users
		WHERE login = ?
	`
	
	var user SystemUser
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.SystemUnitID,
		&user.Active,
		&user.FrontpageID,
		&user.EquipeID,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by login: %w", err)
	}

	return &user, nil
}

func (r *userRepository) List(ctx context.Context, opts ListOptions) ([]SystemUser, error) {
	query := `
		SELECT id, login, name, email, password_hash, system_unit_id, active, frontpage_id, equipe_id, phone, address, created_at, updated_at
		FROM system_users
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
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []SystemUser
	for rows.Next() {
		var user SystemUser
		err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.SystemUnitID,
			&user.Active,
			&user.FrontpageID,
			&user.EquipeID,
			&user.Phone,
			&user.Address,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM system_users WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
