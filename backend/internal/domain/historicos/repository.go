package historicos

import "context"

// HistoricoRepository define a interface para persistência de históricos.
type HistoricoRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca um histórico por ID.
	FindByID(ctx context.Context, id int64) (*HistoricoRepre, error)

	// List retorna uma lista paginada de históricos.
	List(ctx context.Context, opts ListOptions) ([]HistoricoRepre, error)

	// Create insere um novo histórico.
	Create(ctx context.Context, historico *HistoricoRepre) error

	// Update atualiza um histórico existente.
	Update(ctx context.Context, historico *HistoricoRepre) error

	// SoftDelete marca um histórico como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de históricos.
type ListOptions struct {
	Limit       int
	Offset      int
	VisID       *int64
	LeadID      *int64
	Login       *string
	MotivoID    *int64
	OcorridoID  *int64
}
