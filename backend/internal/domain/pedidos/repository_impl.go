package pedidos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// pedidoRepository é a implementação concreta de PedidoRepository para MySQL.
type pedidoRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "pedidos"
}

// NewPedidoRepository cria uma nova instância de pedidoRepository.
func NewPedidoRepository(db *sql.DB, alias common.DBAlias) PedidoRepository {
	return &pedidoRepository{
		db:    db,
		alias: alias,
		table: "pedidos",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *pedidoRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo pedido no banco de dados.
func (r *pedidoRepository) Create(ctx context.Context, pedido *Pedido) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(lead_id, n_ped, dt_ped, dt_prev, dt_quit, dt_entr, status_id, canal_id, fpgto_id, cpgto_id, login_repre, login, total_ped, entrada_ped, taxa_financeira, valor_liquido, obs_ped, obs_ped_ger, img_ped)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		pedido.LeadID,
		pedido.NPed,
		pedido.DtPed,
		pedido.DtPrev,
		pedido.DtQuit,
		pedido.DtEntr,
		pedido.StatusID,
		pedido.CanalID,
		pedido.FpgtoID,
		pedido.CpgtoID,
		pedido.LoginRepre,
		pedido.Login,
		pedido.TotalPed,
		pedido.EntradaPed,
		pedido.TaxaFinanceira,
		pedido.ValorLiquido,
		pedido.ObsPed,
		pedido.ObsPedGer,
		pedido.ImgPed,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar pedido: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do pedido criado: %w", err)
	}

	pedido.ID = id
	return nil
}

// FindByID busca um pedido por ID.
func (r *pedidoRepository) FindByID(ctx context.Context, id int64) (*Pedido, error) {
	query := fmt.Sprintf(`
		SELECT id, lead_id, n_ped, dt_ped, dt_prev, dt_quit, dt_entr, status_id, canal_id, fpgto_id, cpgto_id, login_repre, login, total_ped, entrada_ped, taxa_financeira, valor_liquido, obs_ped, obs_ped_ger, img_ped, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var pedido Pedido
	var dtPed, dtPrev, dtQuit, dtEntr, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pedido.ID,
		&pedido.LeadID,
		&pedido.NPed,
		&dtPed,
		&dtPrev,
		&dtQuit,
		&dtEntr,
		&pedido.StatusID,
		&pedido.CanalID,
		&pedido.FpgtoID,
		&pedido.CpgtoID,
		&pedido.LoginRepre,
		&pedido.Login,
		&pedido.TotalPed,
		&pedido.EntradaPed,
		&pedido.TaxaFinanceira,
		&pedido.ValorLiquido,
		&pedido.ObsPed,
		&pedido.ObsPedGer,
		&pedido.ImgPed,
		&pedido.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar pedido por ID: %w", err)
	}

	if dtPed.Valid {
		pedido.DtPed = &dtPed.Time
	}
	if dtPrev.Valid {
		pedido.DtPrev = &dtPrev.Time
	}
	if dtQuit.Valid {
		pedido.DtQuit = &dtQuit.Time
	}
	if dtEntr.Valid {
		pedido.DtEntr = &dtEntr.Time
	}
	if alteradoEm.Valid {
		pedido.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		pedido.ExcluidoEm = &excluidoEm.Time
	}

	return &pedido, nil
}

// List retorna uma lista paginada de pedidos com filtros opcionais.
func (r *pedidoRepository) List(ctx context.Context, opts ListOptions) ([]Pedido, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.LeadID != nil {
		where += " AND lead_id = ?"
		args = append(args, *opts.LeadID)
	}
	if opts.StatusID != nil {
		where += " AND status_id = ?"
		args = append(args, *opts.StatusID)
	}
	if opts.LoginRepre != nil {
		where += " AND login_repre = ?"
		args = append(args, *opts.LoginRepre)
	}
	if opts.DtPed != nil {
		where += " AND dt_ped = ?"
		args = append(args, *opts.DtPed)
	}

	query := fmt.Sprintf(`
		SELECT id, lead_id, n_ped, dt_ped, dt_prev, dt_quit, dt_entr, status_id, canal_id, fpgto_id, cpgto_id, login_repre, login, total_ped, entrada_ped, taxa_financeira, valor_liquido, obs_ped, obs_ped_ger, img_ped, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY dt_ped DESC, criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pedidos: %w", err)
	}
	defer rows.Close()

	var pedidos []Pedido
	for rows.Next() {
		var pedido Pedido
		var dtPed, dtPrev, dtQuit, dtEntr, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&pedido.ID,
			&pedido.LeadID,
			&pedido.NPed,
			&dtPed,
			&dtPrev,
			&dtQuit,
			&dtEntr,
			&pedido.StatusID,
			&pedido.CanalID,
			&pedido.FpgtoID,
			&pedido.CpgtoID,
			&pedido.LoginRepre,
			&pedido.Login,
			&pedido.TotalPed,
			&pedido.EntradaPed,
			&pedido.TaxaFinanceira,
			&pedido.ValorLiquido,
			&pedido.ObsPed,
			&pedido.ObsPedGer,
			&pedido.ImgPed,
			&pedido.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear pedido: %w", err)
		}

		if dtPed.Valid {
			pedido.DtPed = &dtPed.Time
		}
		if dtPrev.Valid {
			pedido.DtPrev = &dtPrev.Time
		}
		if dtQuit.Valid {
			pedido.DtQuit = &dtQuit.Time
		}
		if dtEntr.Valid {
			pedido.DtEntr = &dtEntr.Time
		}
		if alteradoEm.Valid {
			pedido.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			pedido.ExcluidoEm = &excluidoEm.Time
		}

		pedidos = append(pedidos, pedido)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar pedidos: %w", err)
	}

	return pedidos, nil
}

// Update atualiza um pedido existente.
func (r *pedidoRepository) Update(ctx context.Context, pedido *Pedido) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET lead_id = ?, n_ped = ?, dt_ped = ?, dt_prev = ?, dt_quit = ?, dt_entr = ?, status_id = ?, canal_id = ?, fpgto_id = ?, cpgto_id = ?, login_repre = ?, login = ?, total_ped = ?, entrada_ped = ?, taxa_financeira = ?, valor_liquido = ?, obs_ped = ?, obs_ped_ger = ?, img_ped = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	pedido.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		pedido.LeadID,
		pedido.NPed,
		pedido.DtPed,
		pedido.DtPrev,
		pedido.DtQuit,
		pedido.DtEntr,
		pedido.StatusID,
		pedido.CanalID,
		pedido.FpgtoID,
		pedido.CpgtoID,
		pedido.LoginRepre,
		pedido.Login,
		pedido.TotalPed,
		pedido.EntradaPed,
		pedido.TaxaFinanceira,
		pedido.ValorLiquido,
		pedido.ObsPed,
		pedido.ObsPedGer,
		pedido.ImgPed,
		pedido.AlteradoEm,
		pedido.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar pedido: %w", err)
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

// SoftDelete marca um pedido como excluído (soft delete).
func (r *pedidoRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir pedido: %w", err)
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
