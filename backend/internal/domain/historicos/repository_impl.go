package historicos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// historicoRepository é a implementação concreta de HistoricoRepository para MySQL.
type historicoRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "historicos"
}

// NewHistoricoRepository cria uma nova instância de historicoRepository.
func NewHistoricoRepository(db *sql.DB, alias common.DBAlias) HistoricoRepository {
	return &historicoRepository{
		db:    db,
		alias: alias,
		table: "historicos",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *historicoRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo histórico no banco de dados.
func (r *historicoRepository) Create(ctx context.Context, historico *HistoricoRepre) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(vis_id, lead_id, motivo_id, ocorrido_id, hist, foto_hist, login)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		historico.VisID,
		historico.LeadID,
		historico.MotivoID,
		historico.OcorridoID,
		historico.Hist,
		historico.FotoHist,
		historico.Login,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar histórico: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do histórico criado: %w", err)
	}

	historico.ID = id
	return nil
}

// FindByID busca um histórico por ID.
func (r *historicoRepository) FindByID(ctx context.Context, id int64) (*HistoricoRepre, error) {
	query := fmt.Sprintf(`
		SELECT id, vis_id, lead_id, motivo_id, ocorrido_id, hist, foto_hist, login, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var historico HistoricoRepre
	var alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&historico.ID,
		&historico.VisID,
		&historico.LeadID,
		&historico.MotivoID,
		&historico.OcorridoID,
		&historico.Hist,
		&historico.FotoHist,
		&historico.Login,
		&historico.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar histórico por ID: %w", err)
	}

	if alteradoEm.Valid {
		historico.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		historico.ExcluidoEm = &excluidoEm.Time
	}

	return &historico, nil
}

// List retorna uma lista paginada de históricos com filtros opcionais.
func (r *historicoRepository) List(ctx context.Context, opts ListOptions) ([]HistoricoRepre, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.VisID != nil {
		where += " AND vis_id = ?"
		args = append(args, *opts.VisID)
	}
	if opts.LeadID != nil {
		where += " AND lead_id = ?"
		args = append(args, *opts.LeadID)
	}
	if opts.Login != nil {
		where += " AND login = ?"
		args = append(args, *opts.Login)
	}
	if opts.MotivoID != nil {
		where += " AND motivo_id = ?"
		args = append(args, *opts.MotivoID)
	}
	if opts.OcorridoID != nil {
		where += " AND ocorrido_id = ?"
		args = append(args, *opts.OcorridoID)
	}

	query := fmt.Sprintf(`
		SELECT id, vis_id, lead_id, motivo_id, ocorrido_id, hist, foto_hist, login, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar históricos: %w", err)
	}
	defer rows.Close()

	var historicos []HistoricoRepre
	for rows.Next() {
		var historico HistoricoRepre
		var alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&historico.ID,
			&historico.VisID,
			&historico.LeadID,
			&historico.MotivoID,
			&historico.OcorridoID,
			&historico.Hist,
			&historico.FotoHist,
			&historico.Login,
			&historico.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear histórico: %w", err)
		}

		if alteradoEm.Valid {
			historico.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			historico.ExcluidoEm = &excluidoEm.Time
		}

		historicos = append(historicos, historico)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar históricos: %w", err)
	}

	return historicos, nil
}

// Update atualiza um histórico existente.
func (r *historicoRepository) Update(ctx context.Context, historico *HistoricoRepre) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET vis_id = ?, lead_id = ?, motivo_id = ?, ocorrido_id = ?, hist = ?, foto_hist = ?, login = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	historico.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		historico.VisID,
		historico.LeadID,
		historico.MotivoID,
		historico.OcorridoID,
		historico.Hist,
		historico.FotoHist,
		historico.Login,
		historico.AlteradoEm,
		historico.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar histórico: %w", err)
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

// SoftDelete marca um histórico como excluído (soft delete).
func (r *historicoRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir histórico: %w", err)
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
