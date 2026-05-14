package tarefas

import (
	"context"
)

import "github.com/aojunioro/smid_10/backend/internal/domain/common"

// TarefaRepository define operações de CRUD para Tarefa.
type TarefaRepository interface {
	common.Repository

	Create(ctx context.Context, tarefa *Tarefa) error
	Update(ctx context.Context, tarefa *Tarefa) error
	FindByID(ctx context.Context, id int64) (*Tarefa, error)
	List(ctx context.Context, opts ListOptions) ([]Tarefa, error)
}

// ListOptions define opções para listagem paginada.
type ListOptions struct {
	Limit  int
	Offset int
}
