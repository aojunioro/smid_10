package log

import (
	"context"
)

import "github.com/aojunioro/smid_10/backend/internal/domain/common"

// AccessLogRepository define operações de leitura para SystemAccessLog.
type AccessLogRepository interface {
	common.Repository

	List(ctx context.Context, opts ListOptions) ([]SystemAccessLog, error)
}

// ChangeLogRepository define operações de leitura para SystemChangeLog.
type ChangeLogRepository interface {
	common.Repository

	List(ctx context.Context, opts ListOptions) ([]SystemChangeLog, error)
}

// SqlLogRepository define operações de leitura para SystemSqlLog.
type SqlLogRepository interface {
	common.Repository

	List(ctx context.Context, opts ListOptions) ([]SystemSqlLog, error)
}

// RequestLogRepository define operações de leitura para SystemRequestLog.
type RequestLogRepository interface {
	common.Repository

	List(ctx context.Context, opts ListOptions) ([]SystemRequestLog, error)
}

// ListOptions define opções para listagem paginada.
type ListOptions struct {
	Limit  int
	Offset int
}
