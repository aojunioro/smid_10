package produtos

import "context"

// ProdutoRepository define a interface para persistência de produtos.
type ProdutoRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca um produto por ID.
	FindByID(ctx context.Context, id int64) (*Produto, error)

	// List retorna uma lista paginada de produtos.
	List(ctx context.Context, opts ListOptions) ([]Produto, error)

	// Create insere um novo produto.
	Create(ctx context.Context, produto *Produto) error

	// Update atualiza um produto existente.
	Update(ctx context.Context, produto *Produto) error

	// SoftDelete marca um produto como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de produtos.
type ListOptions struct {
	Limit     int
	Offset    int
	CategID   *int64
	MedID     *int64
	Ativo     *string
	Televendas *string
}
