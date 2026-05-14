package log

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type sqlLogRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewSqlLogRepository(db *sql.DB, alias common.DBAlias) SqlLogRepository {
	return &sqlLogRepository{
		db:     db,
		alias: alias,
	}
}

func (r *sqlLogRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *sqlLogRepository) List(ctx context.Context, opts ListOptions) ([]SystemSqlLog, error) {
	query := `
		SELECT id, logdate, login, database_name, sql_command, statement_type, 
		       access_ip, transaction_id, log_trace, session_id, class_name, 
		       php_sapi, request_id, log_year, log_month, log_day
		FROM system_sql_log
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
		return nil, fmt.Errorf("failed to list sql logs: %w", err)
	}
	defer rows.Close()

	var logs []SystemSqlLog
	for rows.Next() {
		var log SystemSqlLog
		err := rows.Scan(
			&log.ID,
			&log.LogDate,
			&log.Login,
			&log.DatabaseName,
			&log.SqlCommand,
			&log.StatementType,
			&log.AccessIP,
			&log.TransactionID,
			&log.LogTrace,
			&log.SessionID,
			&log.ClassName,
			&log.PhpSapi,
			&log.RequestID,
			&log.LogYear,
			&log.LogMonth,
			&log.LogDay,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sql log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
