package compras

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// compraRepository é a implementação concreta de CompraRepository para MySQL.
type compraRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "compras"
}

// NewCompraRepository cria uma nova instância de compraRepository.
func NewCompraRepository(db *sql.DB, alias common.DBAlias) CompraRepository {
	return &compraRepository{
		db:    db,
		alias: alias,
		table: "compras",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *compraRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere uma nova compra no banco de dados.
func (r *compraRepository) Create(ctx context.Context, compra *Compra) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(ped_id, fornec_id, status_id, transp_id, login, dt_compr, dt_coleta, dt_chegada, frete, vlr_compr, n_nf, n_parcelas, dt_pgto)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		compra.PedID,
		compra.FornecID,
		compra.StatusID,
		compra.TranspID,
		compra.Login,
		compra.DtCompr,
		compra.DtColeta,
		compra.DtChegada,
		compra.Frete,
		compra.VlrCompr,
		compra.NF,
		compra.NParcelas,
		compra.DtPgto,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar compra: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da compra criada: %w", err)
	}

	compra.ID = id
	return nil
}

// FindByID busca uma compra por ID.
func (r *compraRepository) FindByID(ctx context.Context, id int64) (*Compra, error) {
	query := fmt.Sprintf(`
		SELECT id, ped_id, fornec_id, status_id, transp_id, login, dt_compr, dt_coleta, dt_chegada, frete, vlr_compr, n_nf, n_parcelas, dt_pgto, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var compra Compra
	var dtCompr, dtColeta, dtChegada, dtPgto, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&compra.ID,
		&compra.PedID,
		&compra.FornecID,
		&compra.StatusID,
		&compra.TranspID,
		&compra.Login,
		&dtCompr,
		&dtColeta,
		&dtChegada,
		&compra.Frete,
		&compra.VlrCompr,
		&compra.NF,
		&compra.NParcelas,
		&dtPgto,
		&compra.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar compra por ID: %w", err)
	}

	if dtCompr.Valid {
		compra.DtCompr = &dtCompr.Time
	}
	if dtColeta.Valid {
		compra.DtColeta = &dtColeta.Time
	}
	if dtChegada.Valid {
		compra.DtChegada = &dtChegada.Time
	}
	if dtPgto.Valid {
		compra.DtPgto = &dtPgto.Time
	}
	if alteradoEm.Valid {
		compra.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		compra.ExcluidoEm = &excluidoEm.Time
	}

	return &compra, nil
}

// List retorna uma lista paginada de compras com filtros opcionais.
func (r *compraRepository) List(ctx context.Context, opts ListOptions) ([]Compra, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.PedID != nil {
		where += " AND ped_id = ?"
		args = append(args, *opts.PedID)
	}
	if opts.FornecID != nil {
		where += " AND fornec_id = ?"
		args = append(args, *opts.FornecID)
	}
	if opts.StatusID != nil {
		where += " AND status_id = ?"
		args = append(args, *opts.StatusID)
	}
	if opts.TranspID != nil {
		where += " AND transp_id = ?"
		args = append(args, *opts.TranspID)
	}
	if opts.Login != nil {
		where += " AND login = ?"
		args = append(args, *opts.Login)
	}

	query := fmt.Sprintf(`
		SELECT id, ped_id, fornec_id, status_id, transp_id, login, dt_compr, dt_coleta, dt_chegada, frete, vlr_compr, n_nf, n_parcelas, dt_pgto, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar compras: %w", err)
	}
	defer rows.Close()

	var compras []Compra
	for rows.Next() {
		var compra Compra
		var dtCompr, dtColeta, dtChegada, dtPgto, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&compra.ID,
			&compra.PedID,
			&compra.FornecID,
			&compra.StatusID,
			&compra.TranspID,
			&compra.Login,
			&dtCompr,
			&dtColeta,
			&dtChegada,
			&compra.Frete,
			&compra.VlrCompr,
			&compra.NF,
			&compra.NParcelas,
			&dtPgto,
			&compra.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear compra: %w", err)
		}

		if dtCompr.Valid {
			compra.DtCompr = &dtCompr.Time
		}
		if dtColeta.Valid {
			compra.DtColeta = &dtColeta.Time
		}
		if dtChegada.Valid {
			compra.DtChegada = &dtChegada.Time
		}
		if dtPgto.Valid {
			compra.DtPgto = &dtPgto.Time
		}
		if alteradoEm.Valid {
			compra.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			compra.ExcluidoEm = &excluidoEm.Time
		}

		compras = append(compras, compra)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar compras: %w", err)
	}

	return compras, nil
}

// Update atualiza uma compra existente.
func (r *compraRepository) Update(ctx context.Context, compra *Compra) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET ped_id = ?, fornec_id = ?, status_id = ?, transp_id = ?, login = ?, dt_compr = ?, dt_coleta = ?, dt_chegada = ?, frete = ?, vlr_compr = ?, n_nf = ?, n_parcelas = ?, dt_pgto = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	compra.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		compra.PedID,
		compra.FornecID,
		compra.StatusID,
		compra.TranspID,
		compra.Login,
		compra.DtCompr,
		compra.DtColeta,
		compra.DtChegada,
		compra.Frete,
		compra.VlrCompr,
		compra.NF,
		compra.NParcelas,
		compra.DtPgto,
		compra.AlteradoEm,
		compra.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar compra: %w", err)
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

// SoftDelete marca uma compra como excluída (soft delete).
func (r *compraRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir compra: %w", err)
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
