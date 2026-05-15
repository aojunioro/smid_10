package produtos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// produtoRepository é a implementação concreta de ProdutoRepository para MySQL.
type produtoRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "produtos"
}

// NewProdutoRepository cria uma nova instância de produtoRepository.
func NewProdutoRepository(db *sql.DB, alias common.DBAlias) ProdutoRepository {
	return &produtoRepository{
		db:    db,
		alias: alias,
		table: "produtos",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *produtoRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo produto no banco de dados.
func (r *produtoRepository) Create(ctx context.Context, produto *Produto) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(nome_prod, categ_id, fornec_id, med_id, modelo_id, vlr_prod_compra, vlr_prod_venda, estoq_min, estoq_max, ativo, televendas)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		produto.NomeProd,
		produto.CategID,
		produto.FornecID,
		produto.MedID,
		produto.ModeloID,
		produto.VlrProdCompra,
		produto.VlrProdVenda,
		produto.EstoqMin,
		produto.EstoqMax,
		produto.Ativo,
		produto.Televendas,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar produto: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do produto criado: %w", err)
	}

	produto.ID = id
	return nil
}

// FindByID busca um produto por ID.
func (r *produtoRepository) FindByID(ctx context.Context, id int64) (*Produto, error) {
	query := fmt.Sprintf(`
		SELECT id, nome_prod, categ_id, fornec_id, med_id, modelo_id, vlr_prod_compra, vlr_prod_venda, estoq_min, estoq_max, ativo, televendas, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var produto Produto
	var alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&produto.ID,
		&produto.NomeProd,
		&produto.CategID,
		&produto.FornecID,
		&produto.MedID,
		&produto.ModeloID,
		&produto.VlrProdCompra,
		&produto.VlrProdVenda,
		&produto.EstoqMin,
		&produto.EstoqMax,
		&produto.Ativo,
		&produto.Televendas,
		&produto.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar produto por ID: %w", err)
	}

	if alteradoEm.Valid {
		produto.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		produto.ExcluidoEm = &excluidoEm.Time
	}

	return &produto, nil
}

// List retorna uma lista paginada de produtos com filtros opcionais.
func (r *produtoRepository) List(ctx context.Context, opts ListOptions) ([]Produto, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.CategID != nil {
		where += " AND categ_id = ?"
		args = append(args, *opts.CategID)
	}
	if opts.MedID != nil {
		where += " AND med_id = ?"
		args = append(args, *opts.MedID)
	}
	if opts.Ativo != nil {
		where += " AND ativo = ?"
		args = append(args, *opts.Ativo)
	}
	if opts.Televendas != nil {
		where += " AND televendas = ?"
		args = append(args, *opts.Televendas)
	}

	query := fmt.Sprintf(`
		SELECT id, nome_prod, categ_id, fornec_id, med_id, modelo_id, vlr_prod_compra, vlr_prod_venda, estoq_min, estoq_max, ativo, televendas, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY nome_prod ASC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos: %w", err)
	}
	defer rows.Close()

	var produtos []Produto
	for rows.Next() {
		var produto Produto
		var alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&produto.ID,
			&produto.NomeProd,
			&produto.CategID,
			&produto.FornecID,
			&produto.MedID,
			&produto.ModeloID,
			&produto.VlrProdCompra,
			&produto.VlrProdVenda,
			&produto.EstoqMin,
			&produto.EstoqMax,
			&produto.Ativo,
			&produto.Televendas,
			&produto.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %w", err)
		}

		if alteradoEm.Valid {
			produto.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			produto.ExcluidoEm = &excluidoEm.Time
		}

		produtos = append(produtos, produto)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar produtos: %w", err)
	}

	return produtos, nil
}

// Update atualiza um produto existente.
func (r *produtoRepository) Update(ctx context.Context, produto *Produto) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET nome_prod = ?, categ_id = ?, fornec_id = ?, med_id = ?, modelo_id = ?, vlr_prod_compra = ?, vlr_prod_venda = ?, estoq_min = ?, estoq_max = ?, ativo = ?, televendas = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	produto.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		produto.NomeProd,
		produto.CategID,
		produto.FornecID,
		produto.MedID,
		produto.ModeloID,
		produto.VlrProdCompra,
		produto.VlrProdVenda,
		produto.EstoqMin,
		produto.EstoqMax,
		produto.Ativo,
		produto.Televendas,
		produto.AlteradoEm,
		produto.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar produto: %w", err)
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

// SoftDelete marca um produto como excluído (soft delete).
func (r *produtoRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir produto: %w", err)
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
