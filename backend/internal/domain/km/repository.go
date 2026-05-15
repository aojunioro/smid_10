package km

import "context"

// KmConfigRepository define a interface para persistência de configurações de KM.
type KmConfigRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca uma configuração de KM por ID.
	FindByID(ctx context.Context, id int64) (*KmConfig, error)

	// List retorna uma lista paginada de configurações de KM.
	List(ctx context.Context, limit, offset int) ([]KmConfig, error)

	// Create insere uma nova configuração de KM.
	Create(ctx context.Context, config *KmConfig) error

	// Update atualiza uma configuração de KM existente.
	Update(ctx context.Context, config *KmConfig) error

	// SoftDelete marca uma configuração de KM como excluída (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// KmValorKmVigenciaRepository define a interface para persistência de valor de KM por vigência.
type KmValorKmVigenciaRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca um valor de KM por vigência por ID.
	FindByID(ctx context.Context, id int64) (*KmValorKmVigencia, error)

	// List retorna uma lista paginada de valores de KM por vigência.
	List(ctx context.Context, limit, offset int) ([]KmValorKmVigencia, error)

	// Create insere um novo valor de KM por vigência.
	Create(ctx context.Context, valorKm *KmValorKmVigencia) error

	// Update atualiza um valor de KM por vigência existente.
	Update(ctx context.Context, valorKm *KmValorKmVigencia) error

	// SoftDelete marca um valor de KM por vigência como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// KmReembolsoLoteRepository define a interface para persistência de lotes de reembolso.
type KmReembolsoLoteRepository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error

	// FindByID busca um lote de reembolso por ID.
	FindByID(ctx context.Context, id int64) (*KmReembolsoLote, error)

	// List retorna uma lista paginada de lotes de reembolso.
	List(ctx context.Context, opts ListOptions) ([]KmReembolsoLote, error)

	// Create insere um novo lote de reembolso.
	Create(ctx context.Context, lote *KmReembolsoLote) error

	// Update atualiza um lote de reembolso existente.
	Update(ctx context.Context, lote *KmReembolsoLote) error

	// SoftDelete marca um lote de reembolso como excluído (soft delete).
	SoftDelete(ctx context.Context, id int64) error
}

// ListOptions define opções de filtro e paginação para listagem de lotes de reembolso.
type ListOptions struct {
	Limit          int
	Offset         int
	LoginRepre     *string
	StatusPagamento *string
}
