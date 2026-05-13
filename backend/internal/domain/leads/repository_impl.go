package leads

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// leadRepository é a implementação concreta de LeadRepository para MySQL.
// Usa database/sql de forma portátil para facilitar migração para PostgreSQL (ADR 0004).
type leadRepository struct {
	db     *sql.DB
	alias  common.DBAlias
	table  string // "lead" (sem prefixo, schema definido pelo DSN)
}

// NewLeadRepository cria uma nova instância de LeadRepository.
func NewLeadRepository(db *sql.DB, alias common.DBAlias) LeadRepository {
	return &leadRepository{
		db:    db,
		alias: alias,
		table: "lead", // nome da tabela sem prefixo de banco
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *leadRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo lead no banco de dados.
// Usa SQL portátil: INSERT com placeholders ?, COALESCE para valores nulos.
func (r *leadRepository) Create(ctx context.Context, lead *Lead) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(nome, telefone, email, status_id, unidade_id, atendente_id, meio_id, midia_id, criado_em, atualizado_em)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	now := time.Now().UTC()
	lead.CriadoEm = now
	lead.AtualizadoEm = now

	result, err := r.db.ExecContext(ctx, query,
		lead.Nome,
		lead.Telefone,
		lead.Email,         // sql.NullString tratado automaticamente por *string
		lead.StatusID,
		lead.UnidadeID,
		lead.AtendenteID,
		lead.MeioID,
		lead.MidiaID,       // sql.NullString tratado automaticamente por *string
		lead.CriadoEm,
		lead.AtualizadoEm,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar lead: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do lead criado: %w", err)
	}

	lead.ID = id
	return nil
}

// Update atualiza um lead existente.
// Usa SQL portátil: UPDATE com placeholders, atualiza atualizado_em.
func (r *leadRepository) Update(ctx context.Context, lead *Lead) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET nome = ?, telefone = ?, email = ?, status_id = ?, unidade_id = ?, atendente_id = ?, meio_id = ?, midia_id = ?, atualizado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	lead.AtualizadoEm = time.Now().UTC()

	result, err := r.db.ExecContext(ctx, query,
		lead.Nome,
		lead.Telefone,
		lead.Email,
		lead.StatusID,
		lead.UnidadeID,
		lead.AtendenteID,
		lead.MeioID,
		lead.MidiaID,
		lead.AtualizadoEm,
		lead.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar lead: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindByID busca um lead por ID.
// Usa SQL portátil: SELECT com WHERE id = ? AND excluido_em IS NULL.
func (r *leadRepository) FindByID(ctx context.Context, id int64) (*Lead, error) {
	query := fmt.Sprintf(`
		SELECT id, nome, telefone, email, status_id, unidade_id, atendente_id, meio_id, midia_id, criado_em, atualizado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var lead Lead
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lead.ID,
		&lead.Nome,
		&lead.Telefone,
		&lead.Email,
		&lead.StatusID,
		&lead.UnidadeID,
		&lead.AtendenteID,
		&lead.MeioID,
		&lead.MidiaID,
		&lead.CriadoEm,
		&lead.AtualizadoEm,
		&lead.ExcluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar lead por ID: %w", err)
	}

	return &lead, nil
}

// List retorna uma lista paginada de leads com filtros opcionais.
// Usa SQL portátil: LIMIT n OFFSET m (equivalente a LIMIT x, y do MySQL),
// COALESCE para valores nulos em filtros.
func (r *leadRepository) List(ctx context.Context, opts ListOptions) ([]Lead, error) {
	// Constrói query dinâmica com filtros opcionais
	where := "excluido_em IS NULL"
	args := []any{}
	argPos := 1

	if opts.UnidadeID != nil {
		where += fmt.Sprintf(" AND unidade_id = $%d", argPos)
		args = append(args, *opts.UnidadeID)
		argPos++
	}
	if opts.AtendenteID != nil {
		where += fmt.Sprintf(" AND atendente_id = $%d", argPos)
		args = append(args, *opts.AtendenteID)
		argPos++
	}
	if opts.StatusID != nil {
		where += fmt.Sprintf(" AND status_id = $%d", argPos)
		args = append(args, *opts.StatusID)
		argPos++
	}

	// Paginação portátil: LIMIT n OFFSET m
	query := fmt.Sprintf(`
		SELECT id, nome, telefone, email, status_id, unidade_id, atendente_id, meio_id, midia_id, criado_em, atualizado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT $%d OFFSET $%d
	`, r.table, where, argPos, argPos+1)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar leads: %w", err)
	}
	defer rows.Close()

	var leads []Lead
	for rows.Next() {
		var lead Lead
		err := rows.Scan(
			&lead.ID,
			&lead.Nome,
			&lead.Telefone,
			&lead.Email,
			&lead.StatusID,
			&lead.UnidadeID,
			&lead.AtendenteID,
			&lead.MeioID,
			&lead.MidiaID,
			&lead.CriadoEm,
			&lead.AtualizadoEm,
			&lead.ExcluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear lead: %w", err)
		}
		leads = append(leads, lead)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar leads: %w", err)
	}

	return leads, nil
}

// SoftDelete marca um lead como excluído (soft delete).
// Usa SQL portátil: UPDATE com excluido_em = NOW().
func (r *leadRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	result, err := r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("erro ao soft delete lead: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
