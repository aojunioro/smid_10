package televendas

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// televendasContatoRepository é a implementação concreta de TelevendasContatoRepository para MySQL.
type televendasContatoRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "televendas_contatos"
}

// NewTelevendasContatoRepository cria uma nova instância de televendasContatoRepository.
func NewTelevendasContatoRepository(db *sql.DB, alias common.DBAlias) TelevendasContatoRepository {
	return &televendasContatoRepository{
		db:    db,
		alias: alias,
		table: "televendas_contatos",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *televendasContatoRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo contato no banco de dados.
func (r *televendasContatoRepository) Create(ctx context.Context, contato *TelevendasContato) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(lead_id, login, status_id, orcam_id, ped_id, visita_id, historico_id, midia_id, unidd_id, login_recep, login_repre, dt_visita, hr_visita, cidade, bairro, status_visita_id, motivo_id, observacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		contato.LeadID,
		contato.Login,
		contato.StatusID,
		contato.OrcamID,
		contato.PedID,
		contato.VisitaID,
		contato.HistoricoID,
		contato.MidiaID,
		contato.UniddID,
		contato.LoginRecep,
		contato.LoginRepre,
		contato.DtVisita,
		contato.HrVisita,
		contato.Cidade,
		contato.Bairro,
		contato.StatusVisitaID,
		contato.MotivoID,
		contato.Observacao,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar contato: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do contato criado: %w", err)
	}

	contato.ID = id
	return nil
}

// FindByID busca um contato por ID.
func (r *televendasContatoRepository) FindByID(ctx context.Context, id int64) (*TelevendasContato, error) {
	query := fmt.Sprintf(`
		SELECT id, lead_id, login, status_id, orcam_id, ped_id, visita_id, historico_id, midia_id, unidd_id, login_recep, login_repre, dt_visita, hr_visita, cidade, bairro, status_visita_id, motivo_id, observacao, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var contato TelevendasContato
	var dtVisita, hrVisita, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&contato.ID,
		&contato.LeadID,
		&contato.Login,
		&contato.StatusID,
		&contato.OrcamID,
		&contato.PedID,
		&contato.VisitaID,
		&contato.HistoricoID,
		&contato.MidiaID,
		&contato.UniddID,
		&contato.LoginRecep,
		&contato.LoginRepre,
		&dtVisita,
		&hrVisita,
		&contato.Cidade,
		&contato.Bairro,
		&contato.StatusVisitaID,
		&contato.MotivoID,
		&contato.Observacao,
		&contato.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar contato por ID: %w", err)
	}

	if dtVisita.Valid {
		contato.DtVisita = &dtVisita.Time
	}
	if hrVisita.Valid {
		contato.HrVisita = &hrVisita.Time
	}
	if alteradoEm.Valid {
		contato.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		contato.ExcluidoEm = &excluidoEm.Time
	}

	return &contato, nil
}

// List retorna uma lista paginada de contatos com filtros opcionais.
func (r *televendasContatoRepository) List(ctx context.Context, opts ListOptions) ([]TelevendasContato, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.LeadID != nil {
		where += " AND lead_id = ?"
		args = append(args, *opts.LeadID)
	}
	if opts.Login != nil {
		where += " AND login = ?"
		args = append(args, *opts.Login)
	}
	if opts.StatusID != nil {
		where += " AND status_id = ?"
		args = append(args, *opts.StatusID)
	}

	query := fmt.Sprintf(`
		SELECT id, lead_id, login, status_id, orcam_id, ped_id, visita_id, historico_id, midia_id, unidd_id, login_recep, login_repre, dt_visita, hr_visita, cidade, bairro, status_visita_id, motivo_id, observacao, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contatos: %w", err)
	}
	defer rows.Close()

	var contatos []TelevendasContato
	for rows.Next() {
		var contato TelevendasContato
		var dtVisita, hrVisita, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&contato.ID,
			&contato.LeadID,
			&contato.Login,
			&contato.StatusID,
			&contato.OrcamID,
			&contato.PedID,
			&contato.VisitaID,
			&contato.HistoricoID,
			&contato.MidiaID,
			&contato.UniddID,
			&contato.LoginRecep,
			&contato.LoginRepre,
			&dtVisita,
			&hrVisita,
			&contato.Cidade,
			&contato.Bairro,
			&contato.StatusVisitaID,
			&contato.MotivoID,
			&contato.Observacao,
			&contato.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear contato: %w", err)
		}

		if dtVisita.Valid {
			contato.DtVisita = &dtVisita.Time
		}
		if hrVisita.Valid {
			contato.HrVisita = &hrVisita.Time
		}
		if alteradoEm.Valid {
			contato.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			contato.ExcluidoEm = &excluidoEm.Time
		}

		contatos = append(contatos, contato)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar contatos: %w", err)
	}

	return contatos, nil
}

// Update atualiza um contato existente.
func (r *televendasContatoRepository) Update(ctx context.Context, contato *TelevendasContato) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET lead_id = ?, login = ?, status_id = ?, orcam_id = ?, ped_id = ?, visita_id = ?, historico_id = ?, midia_id = ?, unidd_id = ?, login_recep = ?, login_repre = ?, dt_visita = ?, hr_visita = ?, cidade = ?, bairro = ?, status_visita_id = ?, motivo_id = ?, observacao = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	contato.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		contato.LeadID,
		contato.Login,
		contato.StatusID,
		contato.OrcamID,
		contato.PedID,
		contato.VisitaID,
		contato.HistoricoID,
		contato.MidiaID,
		contato.UniddID,
		contato.LoginRecep,
		contato.LoginRepre,
		contato.DtVisita,
		contato.HrVisita,
		contato.Cidade,
		contato.Bairro,
		contato.StatusVisitaID,
		contato.MotivoID,
		contato.Observacao,
		contato.AlteradoEm,
		contato.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar contato: %w", err)
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

// SoftDelete marca um contato como excluído (soft delete).
func (r *televendasContatoRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir contato: %w", err)
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
