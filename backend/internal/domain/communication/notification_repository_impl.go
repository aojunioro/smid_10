package communication

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type notificationRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewNotificationRepository(db *sql.DB, alias common.DBAlias) NotificationRepository {
	return &notificationRepository{
		db:     db,
		alias: alias,
	}
}

func (r *notificationRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *notificationRepository) Create(ctx context.Context, notification *SystemNotification) error {
	query := `
		INSERT INTO system_notification (system_user_id, system_user_to_id, subject, message, dt_message, action_url, action_label, icon, checked)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		notification.SystemUserID,
		notification.SystemUserToID,
		notification.Subject,
		notification.Message,
		notification.DtMessage,
		notification.ActionURL,
		notification.ActionLabel,
		notification.Icon,
		notification.Checked,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	notification.ID = id
	return nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *SystemNotification) error {
	query := `
		UPDATE system_notification
		SET subject = ?, message = ?, action_url = ?, action_label = ?, icon = ?, checked = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		notification.Subject,
		notification.Message,
		notification.ActionURL,
		notification.ActionLabel,
		notification.Icon,
		notification.Checked,
		notification.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (r *notificationRepository) FindByID(ctx context.Context, id int64) (*SystemNotification, error) {
	query := `
		SELECT id, system_user_id, system_user_to_id, subject, message, dt_message, 
		       action_url, action_label, icon, checked
		FROM system_notification
		WHERE id = ?
	`
	
	var notification SystemNotification
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.SystemUserID,
		&notification.SystemUserToID,
		&notification.Subject,
		&notification.Message,
		&notification.DtMessage,
		&notification.ActionURL,
		&notification.ActionLabel,
		&notification.Icon,
		&notification.Checked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}

	return &notification, nil
}

func (r *notificationRepository) List(ctx context.Context, opts ListOptions) ([]SystemNotification, error) {
	query := `
		SELECT id, system_user_id, system_user_to_id, subject, message, dt_message, 
		       action_url, action_label, icon, checked
		FROM system_notification
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
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []SystemNotification
	for rows.Next() {
		var notification SystemNotification
		err := rows.Scan(
			&notification.ID,
			&notification.SystemUserID,
			&notification.SystemUserToID,
			&notification.Subject,
			&notification.Message,
			&notification.DtMessage,
			&notification.ActionURL,
			&notification.ActionLabel,
			&notification.Icon,
			&notification.Checked,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
