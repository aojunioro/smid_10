package log

import "context"

type AccessLogService struct {
	accessLogRepo AccessLogRepository
}

func NewAccessLogService(accessLogRepo AccessLogRepository) *AccessLogService {
	return &AccessLogService{
		accessLogRepo: accessLogRepo,
	}
}

func (s *AccessLogService) List(ctx context.Context, limit, offset int) ([]SystemAccessLog, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	logs, err := s.accessLogRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

type ChangeLogService struct {
	changeLogRepo ChangeLogRepository
}

func NewChangeLogService(changeLogRepo ChangeLogRepository) *ChangeLogService {
	return &ChangeLogService{
		changeLogRepo: changeLogRepo,
	}
}

func (s *ChangeLogService) List(ctx context.Context, limit, offset int) ([]SystemChangeLog, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	logs, err := s.changeLogRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

type SqlLogService struct {
	sqlLogRepo SqlLogRepository
}

func NewSqlLogService(sqlLogRepo SqlLogRepository) *SqlLogService {
	return &SqlLogService{
		sqlLogRepo: sqlLogRepo,
	}
}

func (s *SqlLogService) List(ctx context.Context, limit, offset int) ([]SystemSqlLog, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	logs, err := s.sqlLogRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
