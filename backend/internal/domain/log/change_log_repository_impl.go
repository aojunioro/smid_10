package log

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type changeLogRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewChangeLogRepository(db *sql.DB, alias common.DBAlias) ChangeLogRepository {
	return &changeLogRepository{
		db:     db,
		alias: alias,
	}
}

func (r *changeLogRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *changeLogRepository) List(ctx context.Context, opts ListOptions) ([]SystemChangeLog, error) {
	query := `
		SELECT id, logdate, login, tablename, primarykey, pkvalue, operation, 
		       columnname, oldvalue, newvalue, access_ip, transaction_id, 
		       log_trace, session_id, class_name, php_sapi, log_year, log_month, log_day
		FROM system_change_log
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
		return nil, fmt.Errorf("failed to list change logs: %w", err)
	}
	defer rows.Close()

	var logs []SystemChangeLog
	for rows.Next() {
		var log SystemChangeLog
		err := rows.Scan(
			&log.ID,
			&log.LogDate,
			&log.Login,
			&log.TableName,
			&log.PrimaryKey,
			&log.PKValue,
			&log.Operation,
			&log.ColumnName,
			&log.OldValue,
			&log.NewValue,
			&log.AccessIP,
			&log.TransactionID,
			&log.LogTrace,
			&log.SessionID,
			&log.ClassName,
			&log.PhpSapi,
			&log.LogYear,
			&log.LogMonth,
			&log.LogDay,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan change log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
