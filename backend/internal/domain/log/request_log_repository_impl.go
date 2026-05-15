package log

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type requestLogRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "system_request_log"
}

func NewRequestLogRepository(db *sql.DB, alias common.DBAlias) RequestLogRepository {
	return &requestLogRepository{
		db:    db,
		alias: alias,
		table: "system_request_log",
	}
}

func (r *requestLogRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *requestLogRepository) List(ctx context.Context, opts ListOptions) ([]SystemRequestLog, error) {
	query := fmt.Sprintf(`
		SELECT id, endpoint, logdate, log_year, log_month, log_day, sessionid, login, access_ip, class_name, class_method, http_host, server_port, request_uri, request_method, query_string, request_headers, request_body, request_duration
		FROM %s
		ORDER BY logdate DESC
		LIMIT ? OFFSET ?
	`, r.table)

	rows, err := r.db.QueryContext(ctx, query, opts.Limit, opts.Offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar request logs: %w", err)
	}
	defer rows.Close()

	var logs []SystemRequestLog
	for rows.Next() {
		var log SystemRequestLog
		err := rows.Scan(
			&log.ID,
			&log.Endpoint,
			&log.LogDate,
			&log.LogYear,
			&log.LogMonth,
			&log.LogDay,
			&log.SessionID,
			&log.Login,
			&log.AccessIP,
			&log.ClassName,
			&log.ClassMethod,
			&log.HttpHost,
			&log.ServerPort,
			&log.RequestURI,
			&log.RequestMethod,
			&log.QueryString,
			&log.RequestHeaders,
			&log.RequestBody,
			&log.RequestDuration,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear request log: %w", err)
		}

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar request logs: %w", err)
	}

	return logs, nil
}
