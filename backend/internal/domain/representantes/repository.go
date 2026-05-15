package representantes

import "context"

// RepreDespesaExtraRepository define a interface para persistência de despesas extras de representantes.
type RepreDespesaExtraRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca uma despesa extra por ID.
	FindByID(ctx context.Context, id int64) (*RepreDespesaExtra, error)

	// List retorna uma lista paginada de despesas extras.
	List(ctx context.Context, opts ListOptions) ([]RepreDespesaExtra, error)

	// Create insere uma nova despesa extra.
	Create(ctx context.Context, despesa *RepreDespesaExtra) error

	// Update atualiza uma despesa extra existente.
	Update(ctx context.Context, despesa *RepreDespesaExtra) error

	// SoftDelete marca uma despesa extra como excluída (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de despesas extras.
type ListOptions struct {
	Limit   int
	Offset  int
	Login   *string
	CategID *int64
	Status  *string
}
