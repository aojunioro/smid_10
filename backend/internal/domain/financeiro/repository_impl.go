package financeiro

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// finContaPagarRepository é a implementação concreta de FinContaPagarRepository para MySQL.
type finContaPagarRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "fin_contas_pagar"
}

// NewFinContaPagarRepository cria uma nova instância de finContaPagarRepository.
func NewFinContaPagarRepository(db *sql.DB, alias common.DBAlias) FinContaPagarRepository {
	return &finContaPagarRepository{
		db:    db,
		alias: alias,
		table: "fin_contas_pagar",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *finContaPagarRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere uma nova conta a pagar no banco de dados.
func (r *finContaPagarRepository) Create(ctx context.Context, conta *FinContaPagar) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(categoria_id, pedido_id, descricao, valor, dt_vencimento, status, observacao, lancamento_automatico, recorrencia_id, dt_venc_orig)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		conta.CategoriaID,
		conta.PedidoID,
		conta.Descricao,
		conta.Valor,
		conta.DtVencimento,
		conta.Status,
		conta.Observacao,
		conta.LancamentoAuto,
		conta.RecorrenciaID,
		conta.DtVencOrig,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar conta a pagar: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da conta a pagar criada: %w", err)
	}

	conta.ID = id
	return nil
}

// FindByID busca uma conta a pagar por ID.
func (r *finContaPagarRepository) FindByID(ctx context.Context, id int64) (*FinContaPagar, error) {
	query := fmt.Sprintf(`
		SELECT id, categoria_id, pedido_id, descricao, valor, dt_vencimento, dt_pagamento, status, observacao, lancamento_automatico, recorrencia_id, dt_venc_orig, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var conta FinContaPagar
	var dtVencimento, dtPagamento, dtVencOrig, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conta.ID,
		&conta.CategoriaID,
		&conta.PedidoID,
		&conta.Descricao,
		&conta.Valor,
		&dtVencimento,
		&dtPagamento,
		&conta.Status,
		&conta.Observacao,
		&conta.LancamentoAuto,
		&conta.RecorrenciaID,
		&dtVencOrig,
		&conta.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar conta a pagar por ID: %w", err)
	}

	if dtVencimento.Valid {
		conta.DtVencimento = &dtVencimento.Time
	}
	if dtPagamento.Valid {
		conta.DtPagamento = &dtPagamento.Time
	}
	if dtVencOrig.Valid {
		conta.DtVencOrig = &dtVencOrig.Time
	}
	if alteradoEm.Valid {
		conta.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		conta.ExcluidoEm = &excluidoEm.Time
	}

	return &conta, nil
}

// List retorna uma lista paginada de contas a pagar com filtros opcionais.
func (r *finContaPagarRepository) List(ctx context.Context, opts ListOptions) ([]FinContaPagar, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.CategoriaID != nil {
		where += " AND categoria_id = ?"
		args = append(args, *opts.CategoriaID)
	}
	if opts.PedidoID != nil {
		where += " AND pedido_id = ?"
		args = append(args, *opts.PedidoID)
	}
	if opts.Status != nil {
		where += " AND status = ?"
		args = append(args, *opts.Status)
	}

	query := fmt.Sprintf(`
		SELECT id, categoria_id, pedido_id, descricao, valor, dt_vencimento, dt_pagamento, status, observacao, lancamento_automatico, recorrencia_id, dt_venc_orig, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas a pagar: %w", err)
	}
	defer rows.Close()

	var contas []FinContaPagar
	for rows.Next() {
		var conta FinContaPagar
		var dtVencimento, dtPagamento, dtVencOrig, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&conta.ID,
			&conta.CategoriaID,
			&conta.PedidoID,
			&conta.Descricao,
			&conta.Valor,
			&dtVencimento,
			&dtPagamento,
			&conta.Status,
			&conta.Observacao,
			&conta.LancamentoAuto,
			&conta.RecorrenciaID,
			&dtVencOrig,
			&conta.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear conta a pagar: %w", err)
		}

		if dtVencimento.Valid {
			conta.DtVencimento = &dtVencimento.Time
		}
		if dtPagamento.Valid {
			conta.DtPagamento = &dtPagamento.Time
		}
		if dtVencOrig.Valid {
			conta.DtVencOrig = &dtVencOrig.Time
		}
		if alteradoEm.Valid {
			conta.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			conta.ExcluidoEm = &excluidoEm.Time
		}

		contas = append(contas, conta)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar contas a pagar: %w", err)
	}

	return contas, nil
}

// Update atualiza uma conta a pagar existente.
func (r *finContaPagarRepository) Update(ctx context.Context, conta *FinContaPagar) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET categoria_id = ?, pedido_id = ?, descricao = ?, valor = ?, dt_vencimento = ?, dt_pagamento = ?, status = ?, observacao = ?, lancamento_automatico = ?, recorrencia_id = ?, dt_venc_orig = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	conta.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		conta.CategoriaID,
		conta.PedidoID,
		conta.Descricao,
		conta.Valor,
		conta.DtVencimento,
		conta.DtPagamento,
		conta.Status,
		conta.Observacao,
		conta.LancamentoAuto,
		conta.RecorrenciaID,
		conta.DtVencOrig,
		conta.AlteradoEm,
		conta.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar conta a pagar: %w", err)
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

// SoftDelete marca uma conta a pagar como excluída (soft delete).
func (r *finContaPagarRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir conta a pagar: %w", err)
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

// finContaReceberRepository é a implementação concreta de FinContaReceberRepository para MySQL.
type finContaReceberRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "fin_contas_receber"
}

// NewFinContaReceberRepository cria uma nova instância de finContaReceberRepository.
func NewFinContaReceberRepository(db *sql.DB, alias common.DBAlias) FinContaReceberRepository {
	return &finContaReceberRepository{
		db:    db,
		alias: alias,
		table: "fin_contas_receber",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *finContaReceberRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere uma nova conta a receber no banco de dados.
func (r *finContaReceberRepository) Create(ctx context.Context, conta *FinContaReceber) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(categoria_id, pedido_id, cliente_nome, cliente_cpf, descricao, valor, dt_vencimento, status, observacao, lancamento_automatico, recorrencia_id, dt_venc_orig)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		conta.CategoriaID,
		conta.PedidoID,
		conta.ClienteNome,
		conta.ClienteCPF,
		conta.Descricao,
		conta.Valor,
		conta.DtVencimento,
		conta.Status,
		conta.Observacao,
		conta.LancamentoAuto,
		conta.RecorrenciaID,
		conta.DtVencOrig,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar conta a receber: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da conta a receber criada: %w", err)
	}

	conta.ID = id
	return nil
}

// FindByID busca uma conta a receber por ID.
func (r *finContaReceberRepository) FindByID(ctx context.Context, id int64) (*FinContaReceber, error) {
	query := fmt.Sprintf(`
		SELECT id, categoria_id, pedido_id, cliente_nome, cliente_cpf, descricao, valor, dt_vencimento, dt_recebimento, status, observacao, lancamento_automatico, recorrencia_id, dt_venc_orig, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var conta FinContaReceber
	var dtVencimento, dtRecebimento, dtVencOrig, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conta.ID,
		&conta.CategoriaID,
		&conta.PedidoID,
		&conta.ClienteNome,
		&conta.ClienteCPF,
		&conta.Descricao,
		&conta.Valor,
		&dtVencimento,
		&dtRecebimento,
		&conta.Status,
		&conta.Observacao,
		&conta.LancamentoAuto,
		&conta.RecorrenciaID,
		&dtVencOrig,
		&conta.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar conta a receber por ID: %w", err)
	}

	if dtVencimento.Valid {
		conta.DtVencimento = &dtVencimento.Time
	}
	if dtRecebimento.Valid {
		conta.DtRecebimento = &dtRecebimento.Time
	}
	if dtVencOrig.Valid {
		conta.DtVencOrig = &dtVencOrig.Time
	}
	if alteradoEm.Valid {
		conta.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		conta.ExcluidoEm = &excluidoEm.Time
	}

	return &conta, nil
}

// List retorna uma lista paginada de contas a receber com filtros opcionais.
func (r *finContaReceberRepository) List(ctx context.Context, opts ListOptions) ([]FinContaReceber, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.CategoriaID != nil {
		where += " AND categoria_id = ?"
		args = append(args, *opts.CategoriaID)
	}
	if opts.PedidoID != nil {
		where += " AND pedido_id = ?"
		args = append(args, *opts.PedidoID)
	}
	if opts.Status != nil {
		where += " AND status = ?"
		args = append(args, *opts.Status)
	}

	query := fmt.Sprintf(`
		SELECT id, categoria_id, pedido_id, cliente_nome, cliente_cpf, descricao, valor, dt_vencimento, dt_recebimento, status, observacao, lancamento_automatico, recorrencia_id, dt_venc_orig, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas a receber: %w", err)
	}
	defer rows.Close()

	var contas []FinContaReceber
	for rows.Next() {
		var conta FinContaReceber
		var dtVencimento, dtRecebimento, dtVencOrig, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&conta.ID,
			&conta.CategoriaID,
			&conta.PedidoID,
			&conta.ClienteNome,
			&conta.ClienteCPF,
			&conta.Descricao,
			&conta.Valor,
			&dtVencimento,
			&dtRecebimento,
			&conta.Status,
			&conta.Observacao,
			&conta.LancamentoAuto,
			&conta.RecorrenciaID,
			&dtVencOrig,
			&conta.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear conta a receber: %w", err)
		}

		if dtVencimento.Valid {
			conta.DtVencimento = &dtVencimento.Time
		}
		if dtRecebimento.Valid {
			conta.DtRecebimento = &dtRecebimento.Time
		}
		if dtVencOrig.Valid {
			conta.DtVencOrig = &dtVencOrig.Time
		}
		if alteradoEm.Valid {
			conta.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			conta.ExcluidoEm = &excluidoEm.Time
		}

		contas = append(contas, conta)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar contas a receber: %w", err)
	}

	return contas, nil
}

// Update atualiza uma conta a receber existente.
func (r *finContaReceberRepository) Update(ctx context.Context, conta *FinContaReceber) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET categoria_id = ?, pedido_id = ?, cliente_nome = ?, cliente_cpf = ?, descricao = ?, valor = ?, dt_vencimento = ?, dt_recebimento = ?, status = ?, observacao = ?, lancamento_automatico = ?, recorrencia_id = ?, dt_venc_orig = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	conta.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		conta.CategoriaID,
		conta.PedidoID,
		conta.ClienteNome,
		conta.ClienteCPF,
		conta.Descricao,
		conta.Valor,
		conta.DtVencimento,
		conta.DtRecebimento,
		conta.Status,
		conta.Observacao,
		conta.LancamentoAuto,
		conta.RecorrenciaID,
		conta.DtVencOrig,
		conta.AlteradoEm,
		conta.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar conta a receber: %w", err)
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

// SoftDelete marca uma conta a receber como excluída (soft delete).
func (r *finContaReceberRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir conta a receber: %w", err)
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
