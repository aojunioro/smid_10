package admin

import (
	"context"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// UserRepository define operações de CRUD para SystemUser.
type UserRepository interface {
	common.Repository

	Create(ctx context.Context, user *SystemUser) error
	Update(ctx context.Context, user *SystemUser) error
	FindByID(ctx context.Context, id int64) (*SystemUser, error)
	FindByLogin(ctx context.Context, login string) (*SystemUser, error)
	List(ctx context.Context, opts ListOptions) ([]SystemUser, error)
	Delete(ctx context.Context, id int64) error
}

// GroupRepository define operações de CRUD para SystemGroup.
type GroupRepository interface {
	common.Repository

	Create(ctx context.Context, group *SystemGroup) error
	Update(ctx context.Context, group *SystemGroup) error
	FindByID(ctx context.Context, id int64) (*SystemGroup, error)
	List(ctx context.Context, opts ListOptions) ([]SystemGroup, error)
	Delete(ctx context.Context, id int64) error
}

// RoleRepository define operações de CRUD para SystemRole.
type RoleRepository interface {
	common.Repository

	Create(ctx context.Context, role *SystemRole) error
	Update(ctx context.Context, role *SystemRole) error
	FindByID(ctx context.Context, id int64) (*SystemRole, error)
	List(ctx context.Context, opts ListOptions) ([]SystemRole, error)
	Delete(ctx context.Context, id int64) error
}

// UnitRepository define operações de CRUD para SystemUnit.
type UnitRepository interface {
	common.Repository

	Create(ctx context.Context, unit *SystemUnit) error
	Update(ctx context.Context, unit *SystemUnit) error
	FindByID(ctx context.Context, id int64) (*SystemUnit, error)
	List(ctx context.Context, opts ListOptions) ([]SystemUnit, error)
	Delete(ctx context.Context, id int64) error
}

// ProgramRepository define operações de CRUD para SystemProgram.
type ProgramRepository interface {
	common.Repository

	Create(ctx context.Context, program *SystemProgram) error
	Update(ctx context.Context, program *SystemProgram) error
	FindByID(ctx context.Context, id int64) (*SystemProgram, error)
	List(ctx context.Context, opts ListOptions) ([]SystemProgram, error)
	Delete(ctx context.Context, id int64) error
}

// ListOptions define opções para listagem paginada.
type ListOptions struct {
	Limit  int
	Offset int
}
