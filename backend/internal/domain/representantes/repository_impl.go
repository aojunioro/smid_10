package representantes

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// repreDespesaExtraRepository é a implementação concreta de RepreDespesaExtraRepository para MySQL.
type repreDespesaExtraRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "repre_despesa_extra"
}

// NewRepreDespesaExtraRepository cria uma nova instância de repreDespesaExtraRepository.
func NewRepreDespesaExtraRepository(db *sql.DB, alias common.DBAlias) RepreDespesaExtraRepository {
	return &repreDespesaExtraRepository{
		db:    db,
		alias: alias,
		table: "repre_despesa_extra",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *repreDespesaExtraRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere uma nova despesa extra no banco de dados.
func (r *repreDespesaExtraRepository) Create(ctx context.Context, despesa *RepreDespesaExtra) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(login, data_despesa, categ_id, valor, descricao, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		despesa.Login,
		despesa.DataDespesa,
		despesa.CategID,
		despesa.Valor,
		despesa.Descricao,
		despesa.Status,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar despesa extra: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da despesa extra criada: %w", err)
	}

	despesa.ID = id
	return nil
}

// FindByID busca uma despesa extra por ID.
func (r *repreDespesaExtraRepository) FindByID(ctx context.Context, id int64) (*RepreDespesaExtra, error) {
	query := fmt.Sprintf(`
		SELECT id, login, data_despesa, categ_id, valor, descricao, status, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var despesa RepreDespesaExtra
	var dataDespesa, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&despesa.ID,
		&despesa.Login,
		&dataDespesa,
		&despesa.CategID,
		&despesa.Valor,
		&despesa.Descricao,
		&despesa.Status,
		&despesa.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar despesa extra por ID: %w", err)
	}

	if dataDespesa.Valid {
		despesa.DataDespesa = &dataDespesa.Time
	}
	if alteradoEm.Valid {
		despesa.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		despesa.ExcluidoEm = &excluidoEm.Time
	}

	return &despesa, nil
}

// List retorna uma lista paginada de despesas extras com filtros opcionais.
func (r *repreDespesaExtraRepository) List(ctx context.Context, opts ListOptions) ([]RepreDespesaExtra, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.Login != nil {
		where += " AND login = ?"
		args = append(args, *opts.Login)
	}
	if opts.CategID != nil {
		where += " AND categ_id = ?"
		args = append(args, *opts.CategID)
	}
	if opts.Status != nil {
		where += " AND status = ?"
		args = append(args, *opts.Status)
	}

	query := fmt.Sprintf(`
		SELECT id, login, data_despesa, categ_id, valor, descricao, status, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar despesas extras: %w", err)
	}
	defer rows.Close()

	var despesas []RepreDespesaExtra
	for rows.Next() {
		var despesa RepreDespesaExtra
		var dataDespesa, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&despesa.ID,
			&despesa.Login,
			&dataDespesa,
			&despesa.CategID,
			&despesa.Valor,
			&despesa.Descricao,
			&despesa.Status,
			&despesa.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear despesa extra: %w", err)
		}

		if dataDespesa.Valid {
			despesa.DataDespesa = &dataDespesa.Time
		}
		if alteradoEm.Valid {
			despesa.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			despesa.ExcluidoEm = &excluidoEm.Time
		}

		despesas = append(despesas, despesa)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar despesas extras: %w", err)
	}

	return despesas, nil
}

// Update atualiza uma despesa extra existente.
func (r *repreDespesaExtraRepository) Update(ctx context.Context, despesa *RepreDespesaExtra) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET login = ?, data_despesa = ?, categ_id = ?, valor = ?, descricao = ?, status = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	despesa.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		despesa.Login,
		despesa.DataDespesa,
		despesa.CategID,
		despesa.Valor,
		despesa.Descricao,
		despesa.Status,
		despesa.AlteradoEm,
		despesa.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar despesa extra: %w", err)
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

// SoftDelete marca uma despesa extra como excluída (soft delete).
func (r *repreDespesaExtraRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir despesa extra: %w", err)
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
