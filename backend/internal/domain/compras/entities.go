package compras

import "time"

// Compra representa um registro de aquisição de produtos junto a fornecedor.
type Compra struct {
	ID         int64      `json:"id"`
	PedID      *int64     `json:"ped_id"`
	FornecID   *int64     `json:"fornec_id"`
	StatusID   *int64     `json:"status_id"`
	TranspID   *int64     `json:"transp_id"`
	Login      *string    `json:"login"`
	DtCompr    *time.Time `json:"dt_compr"`
	DtColeta   *time.Time `json:"dt_coleta"`
	DtChegada  *time.Time `json:"dt_chegada"`
	Frete      *float64   `json:"frete"`
	VlrCompr   *float64   `json:"vlr_compr"`
	NF         *string    `json:"n_nf"`
	NParcelas  *int       `json:"n_parcelas"`
	DtPgto     *time.Time `json:"dt_pgto"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}
