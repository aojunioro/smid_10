package financeiro

import "context"

// FinContaPagarRepository define a interface para persistência de contas a pagar.
type FinContaPagarRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca uma conta a pagar por ID.
	FindByID(ctx context.Context, id int64) (*FinContaPagar, error)

	// List retorna uma lista paginada de contas a pagar.
	List(ctx context.Context, opts ListOptions) ([]FinContaPagar, error)

	// Create insere uma nova conta a pagar.
	Create(ctx context.Context, conta *FinContaPagar) error

	// Update atualiza uma conta a pagar existente.
	Update(ctx context.Context, conta *FinContaPagar) error

	// SoftDelete marca uma conta a pagar como excluída (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// FinContaReceberRepository define a interface para persistência de contas a receber.
type FinContaReceberRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca uma conta a receber por ID.
	FindByID(ctx context.Context, id int64) (*FinContaReceber, error)

	// List retorna uma lista paginada de contas a receber.
	List(ctx context.Context, opts ListOptions) ([]FinContaReceber, error)

	// Create insere uma nova conta a receber.
	Create(ctx context.Context, conta *FinContaReceber) error

	// Update atualiza uma conta a receber existente.
	Update(ctx context.Context, conta *FinContaReceber) error

	// SoftDelete marca uma conta a receber como excluída (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem financeira.
type ListOptions struct {
	Limit      int
	Offset     int
	CategoriaID *int64
	PedidoID   *int64
	Status     *string
}
