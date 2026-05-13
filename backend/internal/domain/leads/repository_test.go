package leads_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
	"github.com/aojunioro/smid_10/backend/internal/domain/leads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// setupMySQLContainer inicia um container MySQL para testes.
// Esta função demonstra o uso de testcontainers para testes de repositório,
// permitindo trocar a imagem de mysql para postgres no dia da migração (ADR 0004).
func setupMySQLContainer(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// Configura container MySQL com MariaDB 10.11 (compatível com produção)
	container, err := mysql.RunContainer(ctx,
		mysql.WithDatabase("smid_test"),
		mysql.WithUsername("test"),
		mysql.WithPassword("test"),
	)
	require.NoError(t, err)

	// Obtém string de conexão
	connStr, err := container.ConnectionString(ctx, "multiStatements=true")
	require.NoError(t, err)

	// Abre conexão com o banco
	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err)

	// Cria tabela de leads para o teste
	_, err = db.Exec(`
		CREATE TABLE lead (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			nome VARCHAR(255) NOT NULL,
			telefone VARCHAR(20) NOT NULL,
			email VARCHAR(255),
			status_id BIGINT NOT NULL,
			unidade_id BIGINT NOT NULL,
			atendente_id BIGINT NOT NULL,
			meio_id BIGINT NOT NULL,
			midia_id VARCHAR(50),
			criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			excluido_em TIMESTAMP NULL
		)
	`)
	require.NoError(t, err)

	// Função de cleanup
	cleanup := func() {
		db.Close()
		container.Terminate(ctx)
	}

	return db, cleanup
}

func TestLeadRepository_Create(t *testing.T) {
	db, cleanup := setupMySQLContainer(t)
	defer cleanup()

	repo := leads.NewLeadRepository(db, common.AliasSmid)

	lead := &leads.Lead{
		Nome:        "João Silva",
		Telefone:    "11999999999",
		Email:       ptrString("joao@example.com"),
		StatusID:    1,
		UnidadeID:   1,
		AtendenteID: 1,
		MeioID:      1,
		MidiaID:     ptrString("google"),
	}

	err := repo.Create(context.Background(), lead)
	assert.NoError(t, err)
	assert.Greater(t, lead.ID, int64(0))
	assert.False(t, lead.CriadoEm.IsZero())
}

func TestLeadRepository_FindByID(t *testing.T) {
	db, cleanup := setupMySQLContainer(t)
	defer cleanup()

	repo := leads.NewLeadRepository(db, common.AliasSmid)

	// Cria lead
	lead := &leads.Lead{
		Nome:        "Maria Santos",
		Telefone:    "11888888888",
		Email:       ptrString("maria@example.com"),
		StatusID:    1,
		UnidadeID:   1,
		AtendenteID: 1,
		MeioID:      1,
	}
	err := repo.Create(context.Background(), lead)
	require.NoError(t, err)

	// Busca por ID
	found, err := repo.FindByID(context.Background(), lead.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, lead.Nome, found.Nome)
	assert.Equal(t, lead.Telefone, found.Telefone)

	// ID inexistente
	notFound, err := repo.FindByID(context.Background(), 99999)
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestLeadRepository_Update(t *testing.T) {
	db, cleanup := setupMySQLContainer(t)
	defer cleanup()

	repo := leads.NewLeadRepository(db, common.AliasSmid)

	// Cria lead
	lead := &leads.Lead{
		Nome:        "Pedro Costa",
		Telefone:    "11777777777",
		Email:       ptrString("pedro@example.com"),
		StatusID:    1,
		UnidadeID:   1,
		AtendenteID: 1,
		MeioID:      1,
	}
	err := repo.Create(context.Background(), lead)
	require.NoError(t, err)

	// Atualiza
	lead.Nome = "Pedro Costa Silva"
	lead.StatusID = 2
	err = repo.Update(context.Background(), lead)
	assert.NoError(t, err)

	// Verifica atualização
	updated, err := repo.FindByID(context.Background(), lead.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Pedro Costa Silva", updated.Nome)
	assert.Equal(t, int64(2), updated.StatusID)
}

func TestLeadRepository_SoftDelete(t *testing.T) {
	db, cleanup := setupMySQLContainer(t)
	defer cleanup()

	repo := leads.NewLeadRepository(db, common.AliasSmid)

	// Cria lead
	lead := &leads.Lead{
		Nome:        "Ana Lima",
		Telefone:    "11666666666",
		Email:       ptrString("ana@example.com"),
		StatusID:    1,
		UnidadeID:   1,
		AtendenteID: 1,
		MeioID:      1,
	}
	err := repo.Create(context.Background(), lead)
	require.NoError(t, err)

	// Soft delete
	err = repo.SoftDelete(context.Background(), lead.ID)
	assert.NoError(t, err)

	// Verifica que não é mais encontrado
	deleted, err := repo.FindByID(context.Background(), lead.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted) // soft delete não deve ser retornado
}

func TestLeadRepository_List(t *testing.T) {
	db, cleanup := setupMySQLContainer(t)
	defer cleanup()

	repo := leads.NewLeadRepository(db, common.AliasSmid)

	// Cria múltiplos leads
	for i := 0; i < 5; i++ {
		lead := &leads.Lead{
			Nome:        "Lead Teste",
			Telefone:    "1199999999",
			Email:       ptrString("test@example.com"),
			StatusID:    1,
			UnidadeID:   1,
			AtendenteID: 1,
			MeioID:      1,
		}
		err := repo.Create(context.Background(), lead)
		require.NoError(t, err)
	}

	// Lista com paginação
	opts := leads.ListOptions{
		Limit:  3,
		Offset: 0,
	}
	results, err := repo.List(context.Background(), opts)
	assert.NoError(t, err)
	assert.Len(t, results, 3)

	// Segunda página
	opts.Offset = 3
	results, err = repo.List(context.Background(), opts)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestLeadRepository_Ping(t *testing.T) {
	db, cleanup := setupMySQLContainer(t)
	defer cleanup()

	repo := leads.NewLeadRepository(db, common.AliasSmid)

	err := repo.Ping(context.Background())
	assert.NoError(t, err)
}

// ptrString é um helper para criar ponteiro de string
func ptrString(s string) *string {
	return &s
}
