package log

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type accessLogRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewAccessLogRepository(db *sql.DB, alias common.DBAlias) AccessLogRepository {
	return &accessLogRepository{
		db:     db,
		alias: alias,
	}
}

func (r *accessLogRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *accessLogRepository) List(ctx context.Context, opts ListOptions) ([]SystemAccessLog, error) {
	query := `
		SELECT id, sessionid, login, login_time, login_year, login_month, login_day, 
		       logout_time, impersonated, access_ip, impersonated_by
		FROM system_access_log
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
		return nil, fmt.Errorf("failed to list access logs: %w", err)
	}
	defer rows.Close()

	var logs []SystemAccessLog
	for rows.Next() {
		var log SystemAccessLog
		err := rows.Scan(
			&log.ID,
			&log.SessionID,
			&log.Login,
			&log.LoginTime,
			&log.LoginYear,
			&log.LoginMonth,
			&log.LoginDay,
			&log.LogoutTime,
			&log.Impersonated,
			&log.AccessIP,
			&log.ImpersonatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan access log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
