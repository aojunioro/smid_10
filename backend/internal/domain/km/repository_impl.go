package km

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

// kmConfigRepository é a implementação concreta de KmConfigRepository para MySQL.
type kmConfigRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "km_config"
}

// NewKmConfigRepository cria uma nova instância de kmConfigRepository.
func NewKmConfigRepository(db *sql.DB, alias common.DBAlias) KmConfigRepository {
	return &kmConfigRepository{
		db:    db,
		alias: alias,
		table: "km_config",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *kmConfigRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere uma nova configuração de KM no banco de dados.
func (r *kmConfigRepository) Create(ctx context.Context, config *KmConfig) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(gps_accuracy_max_m, gps_distancia_max_lead_m, map_provider, cache_enabled)
		VALUES (?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		config.GpsAccuracyMaxM,
		config.GpsDistanciaMaxLeadM,
		config.MapProvider,
		config.CacheEnabled,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar configuração de KM: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da configuração de KM criada: %w", err)
	}

	config.ID = id
	return nil
}

// FindByID busca uma configuração de KM por ID.
func (r *kmConfigRepository) FindByID(ctx context.Context, id int64) (*KmConfig, error) {
	query := fmt.Sprintf(`
		SELECT id, gps_accuracy_max_m, gps_distancia_max_lead_m, map_provider, cache_enabled, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var config KmConfig
	var alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID,
		&config.GpsAccuracyMaxM,
		&config.GpsDistanciaMaxLeadM,
		&config.MapProvider,
		&config.CacheEnabled,
		&config.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar configuração de KM por ID: %w", err)
	}

	if alteradoEm.Valid {
		config.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		config.ExcluidoEm = &excluidoEm.Time
	}

	return &config, nil
}

// List retorna uma lista paginada de configurações de KM.
func (r *kmConfigRepository) List(ctx context.Context, limit, offset int) ([]KmConfig, error) {
	query := fmt.Sprintf(`
		SELECT id, gps_accuracy_max_m, gps_distancia_max_lead_m, map_provider, cache_enabled, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE excluido_em IS NULL
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar configurações de KM: %w", err)
	}
	defer rows.Close()

	var configs []KmConfig
	for rows.Next() {
		var config KmConfig
		var alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&config.ID,
			&config.GpsAccuracyMaxM,
			&config.GpsDistanciaMaxLeadM,
			&config.MapProvider,
			&config.CacheEnabled,
			&config.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear configuração de KM: %w", err)
		}

		if alteradoEm.Valid {
			config.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			config.ExcluidoEm = &excluidoEm.Time
		}

		configs = append(configs, config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar configurações de KM: %w", err)
	}

	return configs, nil
}

// Update atualiza uma configuração de KM existente.
func (r *kmConfigRepository) Update(ctx context.Context, config *KmConfig) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET gps_accuracy_max_m = ?, gps_distancia_max_lead_m = ?, map_provider = ?, cache_enabled = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	config.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		config.GpsAccuracyMaxM,
		config.GpsDistanciaMaxLeadM,
		config.MapProvider,
		config.CacheEnabled,
		config.AlteradoEm,
		config.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar configuração de KM: %w", err)
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

// SoftDelete marca uma configuração de KM como excluída (soft delete).
func (r *kmConfigRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir configuração de KM: %w", err)
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

// kmValorKmVigenciaRepository é a implementação concreta de KmValorKmVigenciaRepository para MySQL.
type kmValorKmVigenciaRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "km_valor_km_vigencia"
}

// NewKmValorKmVigenciaRepository cria uma nova instância de kmValorKmVigenciaRepository.
func NewKmValorKmVigenciaRepository(db *sql.DB, alias common.DBAlias) KmValorKmVigenciaRepository {
	return &kmValorKmVigenciaRepository{
		db:    db,
		alias: alias,
		table: "km_valor_km_vigencia",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *kmValorKmVigenciaRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo valor de KM por vigência no banco de dados.
func (r *kmValorKmVigenciaRepository) Create(ctx context.Context, valorKm *KmValorKmVigencia) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(dt_inicio, dt_fim, valor_km, observacao)
		VALUES (?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		valorKm.DtInicio,
		valorKm.DtFim,
		valorKm.ValorKm,
		valorKm.Observacao,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar valor de KM por vigência: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do valor de KM por vigência criado: %w", err)
	}

	valorKm.ID = id
	return nil
}

// FindByID busca um valor de KM por vigência por ID.
func (r *kmValorKmVigenciaRepository) FindByID(ctx context.Context, id int64) (*KmValorKmVigencia, error) {
	query := fmt.Sprintf(`
		SELECT id, dt_inicio, dt_fim, valor_km, observacao, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var valorKm KmValorKmVigencia
	var dtInicio, dtFim, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&valorKm.ID,
		&dtInicio,
		&dtFim,
		&valorKm.ValorKm,
		&valorKm.Observacao,
		&valorKm.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar valor de KM por vigência por ID: %w", err)
	}

	if dtInicio.Valid {
		valorKm.DtInicio = &dtInicio.Time
	}
	if dtFim.Valid {
		valorKm.DtFim = &dtFim.Time
	}
	if alteradoEm.Valid {
		valorKm.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		valorKm.ExcluidoEm = &excluidoEm.Time
	}

	return &valorKm, nil
}

// List retorna uma lista paginada de valores de KM por vigência.
func (r *kmValorKmVigenciaRepository) List(ctx context.Context, limit, offset int) ([]KmValorKmVigencia, error) {
	query := fmt.Sprintf(`
		SELECT id, dt_inicio, dt_fim, valor_km, observacao, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE excluido_em IS NULL
		ORDER BY dt_inicio DESC
		LIMIT ? OFFSET ?
	`, r.table)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar valores de KM por vigência: %w", err)
	}
	defer rows.Close()

	var valores []KmValorKmVigencia
	for rows.Next() {
		var valorKm KmValorKmVigencia
		var dtInicio, dtFim, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&valorKm.ID,
			&dtInicio,
			&dtFim,
			&valorKm.ValorKm,
			&valorKm.Observacao,
			&valorKm.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear valor de KM por vigência: %w", err)
		}

		if dtInicio.Valid {
			valorKm.DtInicio = &dtInicio.Time
		}
		if dtFim.Valid {
			valorKm.DtFim = &dtFim.Time
		}
		if alteradoEm.Valid {
			valorKm.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			valorKm.ExcluidoEm = &excluidoEm.Time
		}

		valores = append(valores, valorKm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar valores de KM por vigência: %w", err)
	}

	return valores, nil
}

// Update atualiza um valor de KM por vigência existente.
func (r *kmValorKmVigenciaRepository) Update(ctx context.Context, valorKm *KmValorKmVigencia) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET dt_inicio = ?, dt_fim = ?, valor_km = ?, observacao = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	valorKm.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		valorKm.DtInicio,
		valorKm.DtFim,
		valorKm.ValorKm,
		valorKm.Observacao,
		valorKm.AlteradoEm,
		valorKm.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar valor de KM por vigência: %w", err)
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

// SoftDelete marca um valor de KM por vigência como excluído (soft delete).
func (r *kmValorKmVigenciaRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir valor de KM por vigência: %w", err)
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

// kmReembolsoLoteRepository é a implementação concreta de KmReembolsoLoteRepository para MySQL.
type kmReembolsoLoteRepository struct {
	db    *sql.DB
	alias common.DBAlias
	table string // "km_reembolso_lote"
}

// NewKmReembolsoLoteRepository cria uma nova instância de kmReembolsoLoteRepository.
func NewKmReembolsoLoteRepository(db *sql.DB, alias common.DBAlias) KmReembolsoLoteRepository {
	return &kmReembolsoLoteRepository{
		db:    db,
		alias: alias,
		table: "km_reembolso_lote",
	}
}

// Ping verifica se a conexão com o banco de dados está ativa.
func (r *kmReembolsoLoteRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Create insere um novo lote de reembolso no banco de dados.
func (r *kmReembolsoLoteRepository) Create(ctx context.Context, lote *KmReembolsoLote) error {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(login_repre, dt_inicio, dt_fim, km_total, valor_km_total, valor_total, status_pagamento, pago_por, pago_em, observacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.table)

	result, err := r.db.ExecContext(ctx, query,
		lote.LoginRepre,
		lote.DtInicio,
		lote.DtFim,
		lote.KmTotal,
		lote.ValorKmTotal,
		lote.ValorTotal,
		lote.StatusPagamento,
		lote.PagoPor,
		lote.PagoEm,
		lote.Observacao,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar lote de reembolso: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do lote de reembolso criado: %w", err)
	}

	lote.ID = id
	return nil
}

// FindByID busca um lote de reembolso por ID.
func (r *kmReembolsoLoteRepository) FindByID(ctx context.Context, id int64) (*KmReembolsoLote, error) {
	query := fmt.Sprintf(`
		SELECT id, login_repre, dt_inicio, dt_fim, km_total, valor_km_total, valor_total, status_pagamento, pago_por, pago_em, observacao, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	var lote KmReembolsoLote
	var dtInicio, dtFim, pagoEm, alteradoEm, excluidoEm sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lote.ID,
		&lote.LoginRepre,
		&dtInicio,
		&dtFim,
		&lote.KmTotal,
		&lote.ValorKmTotal,
		&lote.ValorTotal,
		&lote.StatusPagamento,
		&lote.PagoPor,
		&pagoEm,
		&lote.Observacao,
		&lote.CriadoEm,
		&alteradoEm,
		&excluidoEm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar lote de reembolso por ID: %w", err)
	}

	if dtInicio.Valid {
		lote.DtInicio = &dtInicio.Time
	}
	if dtFim.Valid {
		lote.DtFim = &dtFim.Time
	}
	if pagoEm.Valid {
		lote.PagoEm = &pagoEm.Time
	}
	if alteradoEm.Valid {
		lote.AlteradoEm = &alteradoEm.Time
	}
	if excluidoEm.Valid {
		lote.ExcluidoEm = &excluidoEm.Time
	}

	return &lote, nil
}

// List retorna uma lista paginada de lotes de reembolso com filtros opcionais.
func (r *kmReembolsoLoteRepository) List(ctx context.Context, opts ListOptions) ([]KmReembolsoLote, error) {
	where := "excluido_em IS NULL"
	args := []any{}

	if opts.LoginRepre != nil {
		where += " AND login_repre = ?"
		args = append(args, *opts.LoginRepre)
	}
	if opts.StatusPagamento != nil {
		where += " AND status_pagamento = ?"
		args = append(args, *opts.StatusPagamento)
	}

	query := fmt.Sprintf(`
		SELECT id, login_repre, dt_inicio, dt_fim, km_total, valor_km_total, valor_total, status_pagamento, pago_por, pago_em, observacao, criado_em, alterado_em, excluido_em
		FROM %s
		WHERE %s
		ORDER BY criado_em DESC
		LIMIT ? OFFSET ?
	`, r.table, where)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar lotes de reembolso: %w", err)
	}
	defer rows.Close()

	var lotes []KmReembolsoLote
	for rows.Next() {
		var lote KmReembolsoLote
		var dtInicio, dtFim, pagoEm, alteradoEm, excluidoEm sql.NullTime
		err := rows.Scan(
			&lote.ID,
			&lote.LoginRepre,
			&dtInicio,
			&dtFim,
			&lote.KmTotal,
			&lote.ValorKmTotal,
			&lote.ValorTotal,
			&lote.StatusPagamento,
			&lote.PagoPor,
			&pagoEm,
			&lote.Observacao,
			&lote.CriadoEm,
			&alteradoEm,
			&excluidoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear lote de reembolso: %w", err)
		}

		if dtInicio.Valid {
			lote.DtInicio = &dtInicio.Time
		}
		if dtFim.Valid {
			lote.DtFim = &dtFim.Time
		}
		if pagoEm.Valid {
			lote.PagoEm = &pagoEm.Time
		}
		if alteradoEm.Valid {
			lote.AlteradoEm = &alteradoEm.Time
		}
		if excluidoEm.Valid {
			lote.ExcluidoEm = &excluidoEm.Time
		}

		lotes = append(lotes, lote)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar lotes de reembolso: %w", err)
	}

	return lotes, nil
}

// Update atualiza um lote de reembolso existente.
func (r *kmReembolsoLoteRepository) Update(ctx context.Context, lote *KmReembolsoLote) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET login_repre = ?, dt_inicio = ?, dt_fim = ?, km_total = ?, valor_km_total = ?, valor_total = ?, status_pagamento = ?, pago_por = ?, pago_em = ?, observacao = ?, alterado_em = ?
		WHERE id = ? AND excluido_em IS NULL
	`, r.table)

	now := time.Now().UTC()
	lote.AlteradoEm = &now

	result, err := r.db.ExecContext(ctx, query,
		lote.LoginRepre,
		lote.DtInicio,
		lote.DtFim,
		lote.KmTotal,
		lote.ValorKmTotal,
		lote.ValorTotal,
		lote.StatusPagamento,
		lote.PagoPor,
		lote.PagoEm,
		lote.Observacao,
		lote.AlteradoEm,
		lote.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar lote de reembolso: %w", err)
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

// SoftDelete marca um lote de reembolso como excluído (soft delete).
func (r *kmReembolsoLoteRepository) SoftDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET excluido_em = ?
		WHERE id = ?
	`, r.table)

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir lote de reembolso: %w", err)
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
