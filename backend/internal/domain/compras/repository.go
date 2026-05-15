package compras

import "context"

// CompraRepository define a interface para persistência de compras.
type CompraRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca uma compra por ID.
	FindByID(ctx context.Context, id int64) (*Compra, error)

	// List retorna uma lista paginada de compras.
	List(ctx context.Context, opts ListOptions) ([]Compra, error)

	// Create insere uma nova compra.
	Create(ctx context.Context, compra *Compra) error

	// Update atualiza uma compra existente.
	Update(ctx context.Context, compra *Compra) error

	// SoftDelete marca uma compra como excluída (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de compras.
type ListOptions struct {
	Limit     int
	Offset    int
	PedID     *int64
	FornecID  *int64
	StatusID  *int64
	TranspID  *int64
	Login     *string
}
