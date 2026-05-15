package representantes

import "time"

// RepreDespesaExtra representa uma despesa extra lançada pelo representante.
type RepreDespesaExtra struct {
	ID          int64      `json:"id"`
	Login       string     `json:"login"`
	DataDespesa *time.Time `json:"data_despesa"`
	CategID     *int64     `json:"categ_id"`
	Valor       *float64   `json:"valor"`
	Descricao   *string    `json:"descricao"`
	Status      string     `json:"status"`
	CriadoEm    time.Time  `json:"criado_em"`
	AlteradoEm  *time.Time `json:"alterado_em"`
	ExcluidoEm  *time.Time `json:"excluido_em"`
}

// RepreDespesaCateg representa uma categoria de despesa extra.
type RepreDespesaCateg struct {
	ID         int64      `json:"id"`
	Categoria  string     `json:"categoria"`
	Descricao  *string    `json:"descricao"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}
