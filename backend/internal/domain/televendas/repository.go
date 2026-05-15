package televendas

import "context"

// TelevendasContatoRepository define a interface para persistência de contatos de televendas.
type TelevendasContatoRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca um contato por ID.
	FindByID(ctx context.Context, id int64) (*TelevendasContato, error)

	// List retorna uma lista paginada de contatos.
	List(ctx context.Context, opts ListOptions) ([]TelevendasContato, error)

	// Create insere um novo contato.
	Create(ctx context.Context, contato *TelevendasContato) error

	// Update atualiza um contato existente.
	Update(ctx context.Context, contato *TelevendasContato) error

	// SoftDelete marca um contato como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de contatos.
type ListOptions struct {
	Limit    int
	Offset   int
	LeadID   *int64
	Login    *string
	StatusID *int64
}
