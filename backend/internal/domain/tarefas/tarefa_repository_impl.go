package tarefas

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aojunioro/smid_10/backend/internal/domain/common"
)

type tarefaRepository struct {
	db     *sql.DB
	alias common.DBAlias
}

func NewTarefaRepository(db *sql.DB, alias common.DBAlias) TarefaRepository {
	return &tarefaRepository{
		db:     db,
		alias: alias,
	}
}

func (r *tarefaRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *tarefaRepository) Create(ctx context.Context, tarefa *Tarefa) error {
	query := `
		INSERT INTO tarefas (tarefa, dt_tarefa, hr_tarefa, status, login, lead_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		tarefa.Tarefa,
		tarefa.DtTarefa,
		tarefa.HrTarefa,
		tarefa.Status,
		tarefa.Login,
		tarefa.LeadID,
	)
	if err != nil {
		return fmt.Errorf("failed to create tarefa: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	tarefa.ID = id
	return nil
}

func (r *tarefaRepository) Update(ctx context.Context, tarefa *Tarefa) error {
	query := `
		UPDATE tarefas
		SET tarefa = ?, dt_tarefa = ?, hr_tarefa = ?, status = ?, login = ?, lead_id = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		tarefa.Tarefa,
		tarefa.DtTarefa,
		tarefa.HrTarefa,
		tarefa.Status,
		tarefa.Login,
		tarefa.LeadID,
		tarefa.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update tarefa: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tarefa not found")
	}

	return nil
}

func (r *tarefaRepository) FindByID(ctx context.Context, id int64) (*Tarefa, error) {
	query := `
		SELECT id, tarefa, dt_tarefa, hr_tarefa, status, login, lead_id
		FROM tarefas
		WHERE id = ?
	`
	
	var tarefa Tarefa
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tarefa.ID,
		&tarefa.Tarefa,
		&tarefa.DtTarefa,
		&tarefa.HrTarefa,
		&tarefa.Status,
		&tarefa.Login,
		&tarefa.LeadID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tarefa not found")
		}
		return nil, fmt.Errorf("failed to find tarefa: %w", err)
	}

	return &tarefa, nil
}

func (r *tarefaRepository) List(ctx context.Context, opts ListOptions) ([]Tarefa, error) {
	query := `
		SELECT id, tarefa, dt_tarefa, hr_tarefa, status, login, lead_id
		FROM tarefas
		ORDER BY id DESC
	`
	
	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tarefas: %w", err)
	}
	defer rows.Close()

	var tarefas []Tarefa
	for rows.Next() {
		var tarefa Tarefa
		err := rows.Scan(
			&tarefa.ID,
			&tarefa.Tarefa,
			&tarefa.DtTarefa,
			&tarefa.HrTarefa,
			&tarefa.Status,
			&tarefa.Login,
			&tarefa.LeadID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tarefa: %w", err)
		}
		tarefas = append(tarefas, tarefa)
	}

	return tarefas, nil
}
