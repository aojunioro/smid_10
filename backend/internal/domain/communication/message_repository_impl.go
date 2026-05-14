package communication

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type messageRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewMessageRepository(db *sql.DB, alias common.DBAlias) MessageRepository {
	return &messageRepository{
		db:     db,
		alias: alias,
	}
}

func (r *messageRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *messageRepository) Create(ctx context.Context, message *SystemMessage) error {
	query := `
		INSERT INTO system_message (system_user_id, system_user_to_id, subject, message, dt_message, checked, removed, viewed, attachments)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		message.SystemUserID,
		message.SystemUserToID,
		message.Subject,
		message.Message,
		message.DtMessage,
		message.Checked,
		message.Removed,
		message.Viewed,
		message.Attachments,
	)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	message.ID = id
	return nil
}

func (r *messageRepository) Update(ctx context.Context, message *SystemMessage) error {
	query := `
		UPDATE system_message
		SET subject = ?, message = ?, checked = ?, removed = ?, viewed = ?, attachments = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		message.Subject,
		message.Message,
		message.Checked,
		message.Removed,
		message.Viewed,
		message.Attachments,
		message.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

func (r *messageRepository) FindByID(ctx context.Context, id int64) (*SystemMessage, error) {
	query := `
		SELECT id, system_user_id, system_user_to_id, subject, message, dt_message, 
		       checked, removed, viewed, attachments
		FROM system_message
		WHERE id = ?
	`
	
	var message SystemMessage
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.ID,
		&message.SystemUserID,
		&message.SystemUserToID,
		&message.Subject,
		&message.Message,
		&message.DtMessage,
		&message.Checked,
		&message.Removed,
		&message.Viewed,
		&message.Attachments,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to find message: %w", err)
	}

	return &message, nil
}

func (r *messageRepository) List(ctx context.Context, opts ListOptions) ([]SystemMessage, error) {
	query := `
		SELECT id, system_user_id, system_user_to_id, subject, message, dt_message, 
		       checked, removed, viewed, attachments
		FROM system_message
		ORDER BY id DESC
	`
	
	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	var messages []SystemMessage
	for rows.Next() {
		var message SystemMessage
		err := rows.Scan(
			&message.ID,
			&message.SystemUserID,
			&message.SystemUserToID,
			&message.Subject,
			&message.Message,
			&message.DtMessage,
			&message.Checked,
			&message.Removed,
			&message.Viewed,
			&message.Attachments,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}
