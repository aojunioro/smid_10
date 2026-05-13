package leads

import (
	"context"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// Lead representa um lead no sistema.
// Tipos de domínio independentes do driver (time.Time, *string, bool)
// para facilitar migração para PostgreSQL pós-cutover (ADR 0004).
type Lead struct {
	ID          int64
	Nome        string
	Telefone    string
	Email       *string // nulável
	StatusID    int64
	UnidadeID   int64
	AtendenteID int64
	MeioID      int64
	MidiaID     *string // nulável
	CriadoEm    time.Time
	AtualizadoEm time.Time
	ExcluidoEm   *time.Time // soft delete
}

// LeadRepository define a interface para operações de persistência de leads.
// Uso de interface permite trocar a implementação (MySQL → PostgreSQL)
// sem alterar handlers/services (ADR 0004).
type LeadRepository interface {
	common.Repository

	// Create insere um novo lead no banco de dados.
	Create(ctx context.Context, lead *Lead) error

	// Update atualiza um lead existente.
	Update(ctx context.Context, lead *Lead) error

	// FindByID busca um lead por ID. Retorna nil se não encontrado.
	FindByID(ctx context.Context, id int64) (*Lead, error)

	// List retorna uma lista paginada de leads com filtros opcionais.
	// O parâmetro limit controla o tamanho da página, offset o deslocamento.
	List(ctx context.Context, opts ListOptions) ([]Lead, error)

	// SoftDelete marca um lead como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define filtros e paginação para listagem de leads.
type ListOptions struct {
	UnidadeID   *int64 // filtro opcional por unidade
	AtendenteID *int64 // filtro opcional por atendente
	StatusID    *int64 // filtro opcional por status
	Limit       int    // tamanho da página
	Offset      int    // deslocamento para paginação
}
