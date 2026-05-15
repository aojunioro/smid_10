package suporte

import "context"

// SuporteRepository define a interface para persistência de chamados de suporte.
type SuporteRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca um suporte por ID.
	FindByID(ctx context.Context, id int64) (*Suporte, error)

	// List retorna uma lista paginada de suportes.
	List(ctx context.Context, opts ListOptions) ([]Suporte, error)

	// Create insere um novo suporte.
	Create(ctx context.Context, suporte *Suporte) error

	// Update atualiza um suporte existente.
	Update(ctx context.Context, suporte *Suporte) error

	// SoftDelete marca um suporte como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de suportes.
type ListOptions struct {
	Limit      int
	Offset     int
	PedID      *int64
	Login      *string
	AtribLogin *string
	StatusID   *int64
}
