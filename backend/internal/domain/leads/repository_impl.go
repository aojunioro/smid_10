package leads

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// leadRepository é a implementação concreta de LeadRepository para MySQL.
// Usa database/sql de forma portátil para facilitar migração para PostgreSQL (ADR 0004).
type leadRepository struct {
	db     *sql.DB
	alias  common.DBAlias
	table  string // "leads" (nome da tabela)
}

// NewLeadRepository cria uma nova instância de LeadRepository.
func NewLeadRepository(db *sql.DB, alias common.DBAlias) LeadRepository {
	return &leadRepository{
		db:    db,
		alias: alias,
		table: "leads", // nome da tabela
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *leadRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo lead no banco de dados.
func (r *leadRepository) Create(ctx context.Context, lead *Lead) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(fone1, starttime, fone2, nome, profissao, idade, patologia, nome_acomp, profis_acomp, idd_acomp, pato_acomp, midia_id, tent_id, contato_ok, status_id, unidd_id, meio_id, mot_pend_id, mot_perd_id, email, obs_curta_lead, login, login_recep, login_super, criado_em, alterado_em)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	now := time.Now().UTC()
	lead.CriadoEm = now
	lead.AlteradoEm = now

	result, err := r.db.ExecContext(ctx, query,
		lead.Fone1,
		lead.StartTime,
		lead.Fone2,
		lead.Nome,
		lead.Profissao,
		lead.Idade,
		lead.Patologia,
		lead.NomeAcomp,
		lead.ProfisAcomp,
		lead.IddAcomp,
		lead.PatoAcomp,
		lead.MidiaID,
		lead.TentID,
		lead.ContatoOK,
		lead.StatusID,
		lead.UniddID,
		lead.MeioID,
		lead.MotPendID,
		lead.MotPerdID,
		lead.Email,
		lead.ObsCurtaLead,
		lead.Login,
		lead.LoginRecep,
		lead.LoginSuper,
		lead.CriadoEm,
		lead.AlteradoEm,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar lead: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do lead criado: %w", err)
	}

	lead.ID = id
	return nil
}

// Update atualiza um lead existente.
func (r *leadRepository) Update(ctx context.Context, lead *Lead) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET fone1 = ?, starttime = ?, fone2 = ?, nome = ?, profissao = ?, idade = ?, patologia = ?, nome_acomp = ?, profis_acomp = ?, idd_acomp = ?, pato_acomp = ?, midia_id = ?, tent_id = ?, contato_ok = ?, status_id = ?, unidd_id = ?, meio_id = ?, mot_pend_id = ?, mot_perd_id = ?, email = ?, obs_curta_lead = ?, login = ?, login_recep = ?, login_super = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	lead.AlteradoEm = time.Now().UTC()

	result, err := r.db.ExecContext(ctx, query,
		lead.Fone1,
		lead.StartTime,
		lead.Fone2,
		lead.Nome,
		lead.Profissao,
		lead.Idade,
		lead.Patologia,
		lead.NomeAcomp,
		lead.ProfisAcomp,
		lead.IddAcomp,
		lead.PatoAcomp,
		lead.MidiaID,
		lead.TentID,
		lead.ContatoOK,
		lead.StatusID,
		lead.UniddID,
		lead.MeioID,
		lead.MotPendID,
		lead.MotPerdID,
		lead.Email,
		lead.ObsCurtaLead,
		lead.Login,
		lead.LoginRecep,
		lead.LoginSuper,
		lead.AlteradoEm,
		lead.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar lead: %w", err)
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

// FindByID busca um lead por ID.
func (r *leadRepository) FindByID(ctx context.Context, id int64) (*Lead, error) {
	query := fmt.Sprintf(`
		SELECT id, fone1, starttime, fone2, nome, profissao, idade, patologia, nome_acomp, profis_acomp, idd_acomp, pato_acomp, midia_id, tent_id, contato_ok, status_id, unidd_id, meio_id, mot_pend_id, mot_perd_id, email, obs_curta_lead, login, login_recep, login_super, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var lead Lead
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lead.ID,
		&lead.Fone1,
		&lead.StartTime,
		&lead.Fone2,
		&lead.Nome,
		&lead.Profissao,
		&lead.Idade,
		&lead.Patologia,
		&lead.NomeAcomp,
		&lead.ProfisAcomp,
		&lead.IddAcomp,
		&lead.PatoAcomp,
		&lead.MidiaID,
		&lead.TentID,
		&lead.ContatoOK,
		&lead.StatusID,
		&lead.UniddID,
		&lead.MeioID,
		&lead.MotPendID,
		&lead.MotPerdID,
		&lead.Email,
		&lead.ObsCurtaLead,
		&lead.Login,
		&lead.LoginRecep,
		&lead.LoginSuper,
		&lead.CriadoEm,
		&lead.AlteradoEm,
		&lead.ExcluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar lead por ID: %w", err)
	}

	return &lead, nil
}

// List retorna uma lista paginada de leads com filtros opcionais.
func (r *leadRepository) List(ctx context.Context, opts ListOptions) ([]Lead, error) {
	// Constrói query dinâmica com filtros opcionais
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.UnidadeID != nil {
		where += " AND unidd_id = ?"
		args = append(args, *opts.UnidadeID)
	}
	if opts.AtendenteID != nil {
		where += " AND login_recep = ?"
		args = append(args, *opts.AtendenteID)
	}
	if opts.StatusID != nil {
		where += " AND status_id = ?"
		args = append(args, *opts.StatusID)
	}

	query := fmt.Sprintf(`
		SELECT id, fone1, starttime, fone2, nome, profissao, idade, patologia, nome_acomp, profis_acomp, idd_acomp, pato_acomp, midia_id, tent_id, contato_ok, status_id, unidd_id, meio_id, mot_pend_id, mot_perd_id, email, obs_curta_lead, login, login_recep, login_super, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar leads: %w", err)
	}
	defer rows.Close()

	var leads []Lead
	for rows.Next() {
		var lead Lead
		err := rows.Scan(
			&lead.ID,
			&lead.Fone1,
			&lead.StartTime,
			&lead.Fone2,
			&lead.Nome,
			&lead.Profissao,
			&lead.Idade,
			&lead.Patologia,
			&lead.NomeAcomp,
			&lead.ProfisAcomp,
			&lead.IddAcomp,
			&lead.PatoAcomp,
			&lead.MidiaID,
			&lead.TentID,
			&lead.ContatoOK,
			&lead.StatusID,
			&lead.UniddID,
			&lead.MeioID,
			&lead.MotPendID,
			&lead.MotPerdID,
			&lead.Email,
			&lead.ObsCurtaLead,
			&lead.Login,
			&lead.LoginRecep,
			&lead.LoginSuper,
			&lead.CriadoEm,
			&lead.AlteradoEm,
			&lead.ExcluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear lead: %w", err)
		}
		leads = append(leads, lead)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar leads: %w", err)
	}

	return leads, nil
}

// SoftDelete marca um lead como excluído (soft delete).
// Usa SQL portátil: UPDATE com excluido_em = NOW().
func (r *leadRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	result, err := r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("erro ao soft delete lead: %w", err)
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
