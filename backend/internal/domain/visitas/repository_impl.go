package visitas

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// visitaRepository é a implementação concreta de VisitaRepository para MySQL.
type visitaRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "visitas"
}

// NewVisitaRepository cria uma nova instância de visitaRepository.
func NewVisitaRepository(db *sql.DB, alias common.DBAlias) VisitaRepository {
	return &visitaRepository{
		db:    db,
		alias: alias,
		table: "visitas",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *visitaRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere uma nova visita no banco de dados.
func (r *visitaRepository) Create(ctx context.Context, visita *Visita) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(lead_id, status_id, login_recep, login_repre, dt_visita, hr_visita, confirm, login_conf, dt_confirm, interesse, hist_feito, pos_feito, stts_lead)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		visita.LeadID,
		visita.StatusID,
		visita.LoginRecep,
		visita.LoginRepre,
		visita.DtVisita,
		visita.HrVisita,
		visita.Confirm,
		visita.LoginConf,
		visita.DtConfirm,
		visita.Interesse,
		visita.HistFeito,
		visita.PosFeito,
		visita.SttsLead,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar visita: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da visita criada: %w", err)
	}

	visita.ID = id
	return nil
}

// FindByID busca uma visita por ID.
func (r *visitaRepository) FindByID(ctx context.Context, id int64) (*Visita, error) {
	query := fmt.Sprintf(`
		SELECT id, lead_id, status_id, login_recep, login_repre, dt_visita, hr_visita, confirm, login_conf, dt_confirm, interesse, hist_feito, pos_feito, stts_lead, criado_em, COALESCE(alterado_em, '1970-01-01 00:00:00'), COALESCE(excluido_em, '1970-01-01 00:00:00')
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var visita Visita
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&visita.ID,
		&visita.LeadID,
		&visita.StatusID,
		&visita.LoginRecep,
		&visita.LoginRepre,
		&visita.DtVisita,
		&visita.HrVisita,
		&visita.Confirm,
		&visita.LoginConf,
		&visita.DtConfirm,
		&visita.Interesse,
		&visita.HistFeito,
		&visita.PosFeito,
		&visita.SttsLead,
		&visita.CriadoEm,
		&visita.AlteradoEm,
		&visita.ExcluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar visita por ID: %w", err)
	}

	// Converter data padrão para nil se for 1970-01-01
	if visita.AlteradoEm != nil && visita.AlteradoEm.Year() == 1970 {
		visita.AlteradoEm = nil
	}
	if visita.ExcluidoEm != nil && visita.ExcluidoEm.Year() == 1970 {
		visita.ExcluidoEm = nil
	}

	return &visita, nil
}

// List retorna uma lista paginada de visitas com filtros opcionais.
func (r *visitaRepository) List(ctx context.Context, opts ListOptions) ([]Visita, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.LeadID != nil {
		where += " AND lead_id = ?"
		args = append(args, *opts.LeadID)
	}
	if opts.LoginRepre != nil {
		where += " AND login_repre = ?"
		args = append(args, *opts.LoginRepre)
	}
	if opts.LoginRecep != nil {
		where += " AND login_recep = ?"
		args = append(args, *opts.LoginRecep)
	}
	if opts.StatusID != nil {
		where += " AND status_id = ?"
		args = append(args, *opts.StatusID)
	}
	if opts.DtVisita != nil {
		where += " AND dt_visita = ?"
		args = append(args, *opts.DtVisita)
	}

	query := fmt.Sprintf(`
		SELECT id, lead_id, status_id, login_recep, login_repre, dt_visita, hr_visita, confirm, login_conf, dt_confirm, interesse, hist_feito, pos_feito, stts_lead, criado_em, COALESCE(alterado_em, '1970-01-01 00:00:00'), COALESCE(excluido_em, '1970-01-01 00:00:00')
		FROM %s
		WHERE %s
		ORDER BY dt_visita DESC, hr_visita DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar visitas: %w", err)
	}
	defer rows.Close()

	var visitas []Visita
	for rows.Next() {
		var visita Visita
		err := rows.Scan(
			&visita.ID,
			&visita.LeadID,
			&visita.StatusID,
			&visita.LoginRecep,
			&visita.LoginRepre,
			&visita.DtVisita,
			&visita.HrVisita,
			&visita.Confirm,
			&visita.LoginConf,
			&visita.DtConfirm,
			&visita.Interesse,
			&visita.HistFeito,
			&visita.PosFeito,
			&visita.SttsLead,
			&visita.CriadoEm,
			&visita.AlteradoEm,
			&visita.ExcluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear visita: %w", err)
		}

		// Converter data padrão para nil se for 1970-01-01
		if visita.AlteradoEm != nil && visita.AlteradoEm.Year() == 1970 {
			visita.AlteradoEm = nil
		}
		if visita.ExcluidoEm != nil && visita.ExcluidoEm.Year() == 1970 {
			visita.ExcluidoEm = nil
		}

		visitas = append(visitas, visita)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar visitas: %w", err)
	}

	return visitas, nil
}

// Update atualiza uma visita existente.
func (r *visitaRepository) Update(ctx context.Context, visita *Visita) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET lead_id = ?, status_id = ?, login_recep = ?, login_repre = ?, dt_visita = ?, hr_visita = ?, confirm = ?, login_conf = ?, dt_confirm = ?, interesse = ?, hist_feito = ?, pos_feito = ?, stts_lead = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	visita.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		visita.LeadID,
		visita.StatusID,
		visita.LoginRecep,
		visita.LoginRepre,
		visita.DtVisita,
		visita.HrVisita,
		visita.Confirm,
		visita.LoginConf,
		visita.DtConfirm,
		visita.Interesse,
		visita.HistFeito,
		visita.PosFeito,
		visita.SttsLead,
		visita.AlteradoEm,
		visita.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar visita: %w", err)
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

// SoftDelete marca uma visita como excluída (soft delete).
func (r *visitaRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir visita: %w", err)
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
