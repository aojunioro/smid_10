package visitas

import "context"

// VisitaRepository define a interface para persistência de visitas.
type VisitaRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca uma visita por ID.
	FindByID(ctx context.Context, id int64) (*Visita, error)

	// List retorna uma lista paginada de visitas.
	List(ctx context.Context, opts ListOptions) ([]Visita, error)

	// Create insere uma nova visita.
	Create(ctx context.Context, visita *Visita) error

	// Update atualiza uma visita existente.
	Update(ctx context.Context, visita *Visita) error

	// SoftDelete marca uma visita como excluída (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de visitas.
type ListOptions struct {
	Limit       int
	Offset      int
	LeadID      *int64
	LoginRepre  *string
	LoginRecep  *string
	StatusID    *int64
	DtVisita    *string
	UnidadeID   *int64
}
