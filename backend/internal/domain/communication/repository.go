package communication

import (
	"context"
)

import "github.com/aojunioro/smid_10/backend/internal/domain/common"

// NotificationRepository define operações de CRUD para SystemNotification.
type NotificationRepository interface {
	common.Repository

	Create(ctx context.Context, notification *SystemNotification) error
	Update(ctx context.Context, notification *SystemNotification) error
	FindByID(ctx context.Context, id int64) (*SystemNotification, error)
	List(ctx context.Context, opts ListOptions) ([]SystemNotification, error)
}

// MessageRepository define operações de CRUD para SystemMessage.
type MessageRepository interface {
	common.Repository

	Create(ctx context.Context, message *SystemMessage) error
	Update(ctx context.Context, message *SystemMessage) error
	FindByID(ctx context.Context, id int64) (*SystemMessage, error)
	List(ctx context.Context, opts ListOptions) ([]SystemMessage, error)
}

// ListOptions define opções para listagem paginada.
type ListOptions struct {
	Limit  int
	Offset int
}
