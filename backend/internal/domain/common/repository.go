package common

import (
	"context"
)

// DBAlias representa um dos 4 bancos canônicos do SMID 10.
type DBAlias string

const (
	AliasSmid          DBAlias = "smid"
	AliasPermission    DBAlias = "permission"
	AliasLog           DBAlias = "log"
	AliasCommunication DBAlias = "communication"
)

// Repository é a interface base para todos os repositórios do sistema.
// Define métodos comuns que todo repositório deve implementar, como
// verificação de saúde do banco de dados subjacente.
type Repository interface {
	// Ping verifica se a conexão com o banco de dados está ativa.
	Ping(ctx context.Context) error
}
