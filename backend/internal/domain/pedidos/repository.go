package pedidos

import "context"

// PedidoRepository define a interface para persistência de pedidos.
type PedidoRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca um pedido por ID.
	FindByID(ctx context.Context, id int64) (*Pedido, error)

	// List retorna uma lista paginada de pedidos.
	List(ctx context.Context, opts ListOptions) ([]Pedido, error)

	// Create insere um novo pedido.
	Create(ctx context.Context, pedido *Pedido) error

	// Update atualiza um pedido existente.
	Update(ctx context.Context, pedido *Pedido) error

	// SoftDelete marca um pedido como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de pedidos.
type ListOptions struct {
	Limit      int
	Offset     int
	LeadID     *int64
	StatusID   *int64
	LoginRepre *string
	DtPed      *string
}
