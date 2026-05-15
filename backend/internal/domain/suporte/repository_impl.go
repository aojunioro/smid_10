package suporte

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// suporteRepository é a implementação concreta de SuporteRepository para MySQL.
type suporteRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "suportes"
}

// NewSuporteRepository cria uma nova instância de suporteRepository.
func NewSuporteRepository(db *sql.DB, alias common.DBAlias) SuporteRepository {
	return &suporteRepository{
		db:    db,
		alias: alias,
		table: "suportes",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *suporteRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo suporte no banco de dados.
func (r *suporteRepository) Create(ctx context.Context, suporte *Suporte) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(ped_id, status_id, fone_sup, solicit_id, depart_id, login, atrib_login, prioridade, dt_sup, dt_limit, relato_cli, dt_resol, relato_tec, img_ordem)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		suporte.PedID,
		suporte.StatusID,
		suporte.FoneSup,
		suporte.SolicitID,
		suporte.DepartID,
		suporte.Login,
		suporte.AtribLogin,
		suporte.Prioridade,
		suporte.DtSup,
		suporte.DtLimit,
		suporte.RelatoCli,
		suporte.DtResol,
		suporte.RelatoTec,
		suporte.ImgOrdem,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar suporte: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do suporte criado: %w", err)
	}

	suporte.ID = id
	return nil
}

// FindByID busca um suporte por ID.
func (r *suporteRepository) FindByID(ctx context.Context, id int64) (*Suporte, error) {
	query := fmt.Sprintf(`
		SELECT id, ped_id, status_id, fone_sup, solicit_id, depart_id, login, atrib_login, prioridade, dt_sup, dt_limit, relato_cli, dt_resol, relato_tec, img_ordem, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var suporte Suporte
	var dtSup, dtLimit, dtResol, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&suporte.ID,
		&suporte.PedID,
		&suporte.StatusID,
		&suporte.FoneSup,
		&suporte.SolicitID,
		&suporte.DepartID,
		&suporte.Login,
		&suporte.AtribLogin,
		&suporte.Prioridade,
		&dtSup,
		&dtLimit,
		&suporte.RelatoCli,
		&dtResol,
		&suporte.RelatoTec,
		&suporte.ImgOrdem,
		&suporte.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar suporte por ID: %w", err)
	}

	if dtSup.Valid {
		suporte.DtSup = &dtSup.Time
	}
	if dtLimit.Valid {
		suporte.DtLimit = &dtLimit.Time
	}
	if dtResol.Valid {
		suporte.DtResol = &dtResol.Time
	}
	if alteradoEm.Valid {
		suporte.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		suporte.ExcluidoEm = &excluidoEm.Time
	}

	return &suporte, nil
}

// List retorna uma lista paginada de suportes com filtros opcionais.
func (r *suporteRepository) List(ctx context.Context, opts ListOptions) ([]Suporte, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.PedID != nil {
		where += " AND ped_id = ?"
		args = append(args, *opts.PedID)
	}
	if opts.Login != nil {
		where += " AND login = ?"
		args = append(args, *opts.Login)
	}
	if opts.AtribLogin != nil {
		where += " AND atrib_login = ?"
		args = append(args, *opts.AtribLogin)
	}
	if opts.StatusID != nil {
		where += " AND status_id = ?"
		args = append(args, *opts.StatusID)
	}

	query := fmt.Sprintf(`
		SELECT id, ped_id, status_id, fone_sup, solicit_id, depart_id, login, atrib_login, prioridade, dt_sup, dt_limit, relato_cli, dt_resol, relato_tec, img_ordem, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar suportes: %w", err)
	}
	defer rows.Close()

	var suportes []Suporte
	for rows.Next() {
		var suporte Suporte
		var dtSup, dtLimit, dtResol, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&suporte.ID,
			&suporte.PedID,
			&suporte.StatusID,
			&suporte.FoneSup,
			&suporte.SolicitID,
			&suporte.DepartID,
			&suporte.Login,
			&suporte.AtribLogin,
			&suporte.Prioridade,
			&dtSup,
			&dtLimit,
			&suporte.RelatoCli,
			&dtResol,
			&suporte.RelatoTec,
			&suporte.ImgOrdem,
			&suporte.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear suporte: %w", err)
		}

		if dtSup.Valid {
			suporte.DtSup = &dtSup.Time
		}
		if dtLimit.Valid {
			suporte.DtLimit = &dtLimit.Time
		}
		if dtResol.Valid {
			suporte.DtResol = &dtResol.Time
		}
		if alteradoEm.Valid {
			suporte.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			suporte.ExcluidoEm = &excluidoEm.Time
		}

		suportes = append(suportes, suporte)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar suportes: %w", err)
	}

	return suportes, nil
}

// Update atualiza um suporte existente.
func (r *suporteRepository) Update(ctx context.Context, suporte *Suporte) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET ped_id = ?, status_id = ?, fone_sup = ?, solicit_id = ?, depart_id = ?, login = ?, atrib_login = ?, prioridade = ?, dt_sup = ?, dt_limit = ?, relato_cli = ?, dt_resol = ?, relato_tec = ?, img_ordem = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	suporte.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		suporte.PedID,
		suporte.StatusID,
		suporte.FoneSup,
		suporte.SolicitID,
		suporte.DepartID,
		suporte.Login,
		suporte.AtribLogin,
		suporte.Prioridade,
		suporte.DtSup,
		suporte.DtLimit,
		suporte.RelatoCli,
		suporte.DtResol,
		suporte.RelatoTec,
		suporte.ImgOrdem,
		suporte.AlteradoEm,
		suporte.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar suporte: %w", err)
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

// SoftDelete marca um suporte como excluído (soft delete).
func (r *suporteRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir suporte: %w", err)
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
