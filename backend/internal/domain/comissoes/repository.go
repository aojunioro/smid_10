package comissoes

import "context"

// ComissaoRepository define a interface para persistência de comissões.
type ComissaoRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca uma comissão por ID.
	FindByID(ctx context.Context, id int64) (*Comissao, error)

	// List retorna uma lista paginada de comissões.
	List(ctx context.Context, opts ListOptions) ([]Comissao, error)

	// Create insere uma nova comissão.
	Create(ctx context.Context, comissao *Comissao) error

	// Update atualiza uma comissão existente.
	Update(ctx context.Context, comissao *Comissao) error

	// SoftDelete marca uma comissão como excluída (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ComissItemRepository define a interface para persistência de itens de comissão.
type ComissItemRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca um item de comissão por ID.
	FindByID(ctx context.Context, id int64) (*ComissItem, error)

	// List retorna uma lista paginada de itens de comissão.
	List(ctx context.Context, comisID int64, limit, offset int) ([]ComissItem, error)

	// Create insere um novo item de comissão.
	Create(ctx context.Context, item *ComissItem) error

	// Update atualiza um item de comissão existente.
	Update(ctx context.Context, item *ComissItem) error

	// SoftDelete marca um item de comissão como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de comissões.
type ListOptions struct {
	Limit     int
	Offset    int
	PedID     *int64
	SttComis  *string
}
