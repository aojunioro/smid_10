package comissoes

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// comissaoRepository é a implementação concreta de ComissaoRepository para MySQL.
type comissaoRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "comissoes"
}

// NewComissaoRepository cria uma nova instância de comissaoRepository.
func NewComissaoRepository(db *sql.DB, alias common.DBAlias) ComissaoRepository {
	return &comissaoRepository{
		db:    db,
		alias: alias,
		table: "comissoes",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *comissaoRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere uma nova comissão no banco de dados.
func (r *comissaoRepository) Create(ctx context.Context, comissao *Comissao) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(ped_id, vlr_comissao, total_pago, vlr_saldo, dt_prevista, stt_comis, observacao)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		comissao.PedID,
		comissao.VlrComissao,
		comissao.TotalPago,
		comissao.VlrSaldo,
		comissao.DtPrevista,
		comissao.SttComis,
		comissao.Observacao,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar comissão: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da comissão criada: %w", err)
	}

	comissao.ID = id
	return nil
}

// FindByID busca uma comissão por ID.
func (r *comissaoRepository) FindByID(ctx context.Context, id int64) (*Comissao, error) {
	query := fmt.Sprintf(`
		SELECT id, ped_id, vlr_comissao, total_pago, vlr_saldo, dt_prevista, stt_comis, observacao, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var comissao Comissao
	var dtPrevista, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comissao.ID,
		&comissao.PedID,
		&comissao.VlrComissao,
		&comissao.TotalPago,
		&comissao.VlrSaldo,
		&dtPrevista,
		&comissao.SttComis,
		&comissao.Observacao,
		&comissao.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar comissão por ID: %w", err)
	}

	if dtPrevista.Valid {
		comissao.DtPrevista = &dtPrevista.Time
	}
	if alteradoEm.Valid {
		comissao.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		comissao.ExcluidoEm = &excluidoEm.Time
	}

	return &comissao, nil
}

// List retorna uma lista paginada de comissões com filtros opcionais.
func (r *comissaoRepository) List(ctx context.Context, opts ListOptions) ([]Comissao, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.PedID != nil {
		where += " AND ped_id = ?"
		args = append(args, *opts.PedID)
	}
	if opts.SttComis != nil {
		where += " AND stt_comis = ?"
		args = append(args, *opts.SttComis)
	}

	query := fmt.Sprintf(`
		SELECT id, ped_id, vlr_comissao, total_pago, vlr_saldo, dt_prevista, stt_comis, observacao, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar comissões: %w", err)
	}
	defer rows.Close()

	var comissoes []Comissao
	for rows.Next() {
		var comissao Comissao
		var dtPrevista, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&comissao.ID,
			&comissao.PedID,
			&comissao.VlrComissao,
			&comissao.TotalPago,
			&comissao.VlrSaldo,
			&dtPrevista,
			&comissao.SttComis,
			&comissao.Observacao,
			&comissao.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear comissão: %w", err)
		}

		if dtPrevista.Valid {
			comissao.DtPrevista = &dtPrevista.Time
		}
		if alteradoEm.Valid {
			comissao.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			comissao.ExcluidoEm = &excluidoEm.Time
		}

		comissoes = append(comissoes, comissao)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar comissões: %w", err)
	}

	return comissoes, nil
}

// Update atualiza uma comissão existente.
func (r *comissaoRepository) Update(ctx context.Context, comissao *Comissao) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET ped_id = ?, vlr_comissao = ?, total_pago = ?, vlr_saldo = ?, dt_prevista = ?, stt_comis = ?, observacao = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	comissao.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		comissao.PedID,
		comissao.VlrComissao,
		comissao.TotalPago,
		comissao.VlrSaldo,
		comissao.DtPrevista,
		comissao.SttComis,
		comissao.Observacao,
		comissao.AlteradoEm,
		comissao.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar comissão: %w", err)
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

// SoftDelete marca uma comissão como excluída (soft delete).
func (r *comissaoRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir comissão: %w", err)
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

// comissItemRepository é a implementação concreta de ComissItemRepository para MySQL.
type comissItemRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "comis_ped_item"
}

// NewComissItemRepository cria uma nova instância de comissItemRepository.
func NewComissItemRepository(db *sql.DB, alias common.DBAlias) ComissItemRepository {
	return &comissItemRepository{
		db:    db,
		alias: alias,
		table: "comis_ped_item",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *comissItemRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo item de comissão no banco de dados.
func (r *comissItemRepository) Create(ctx context.Context, item *ComissItem) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(comis_id, vlr_pago, dt_pgto, obs_pgto)
		VALUES (?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		item.ComisID,
		item.VlrPago,
		item.DtPgto,
		item.ObsPgto,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar item de comissão: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do item de comissão criado: %w", err)
	}

	item.ID = id
	return nil
}

// FindByID busca um item de comissão por ID.
func (r *comissItemRepository) FindByID(ctx context.Context, id int64) (*ComissItem, error) {
	query := fmt.Sprintf(`
		SELECT id, comis_id, vlr_pago, dt_pgto, obs_pgto, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var item ComissItem
	var dtPgto, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.ComisID,
		&item.VlrPago,
		&dtPgto,
		&item.ObsPgto,
		&item.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar item de comissão por ID: %w", err)
	}

	if dtPgto.Valid {
		item.DtPgto = &dtPgto.Time
	}
	if alteradoEm.Valid {
		item.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		item.ExcluidoEm = &excluidoEm.Time
	}

	return &item, nil
}

// List retorna uma lista paginada de itens de comissão.
func (r *comissItemRepository) List(ctx context.Context, comisID int64, limit, offset int) ([]ComissItem, error) {
	where := "excluido_em IS NULL AND comis_id = ?"
	args := []any{comisID}

	query := fmt.Sprintf(`
		SELECT id, comis_id, vlr_pago, dt_pgto, obs_pgto, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar itens de comissão: %w", err)
	}
	defer rows.Close()

	var itens []ComissItem
	for rows.Next() {
		var item ComissItem
		var dtPgto, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&item.ID,
			&item.ComisID,
			&item.VlrPago,
			&dtPgto,
			&item.ObsPgto,
			&item.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear item de comissão: %w", err)
		}

		if dtPgto.Valid {
			item.DtPgto = &dtPgto.Time
		}
		if alteradoEm.Valid {
			item.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			item.ExcluidoEm = &excluidoEm.Time
		}

		itens = append(itens, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar itens de comissão: %w", err)
	}

	return itens, nil
}

// Update atualiza um item de comissão existente.
func (r *comissItemRepository) Update(ctx context.Context, item *ComissItem) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET comis_id = ?, vlr_pago = ?, dt_pgto = ?, obs_pgto = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	item.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		item.ComisID,
		item.VlrPago,
		item.DtPgto,
		item.ObsPgto,
		item.AlteradoEm,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar item de comissão: %w", err)
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

// SoftDelete marca um item de comissão como excluído (soft delete).
func (r *comissItemRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir item de comissão: %w", err)
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
