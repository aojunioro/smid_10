package comissoes

import "time"

// Comissao representa um registro agregado da comissão de um pedido.
type Comissao struct {
	ID          int64      `json:"id"`
	PedID       *int64     `json:"ped_id"`
	VlrComissao *float64   `json:"vlr_comissao"`
	TotalPago   *float64   `json:"total_pago"`
	VlrSaldo    *float64   `json:"vlr_saldo"`
	DtPrevista  *time.Time `json:"dt_prevista"`
	SttComis    *string    `json:"stt_comis"`
	Observacao  *string    `json:"observacao"`
	CriadoEm    time.Time  `json:"criado_em"`
	AlteradoEm  *time.Time `json:"alterado_em"`
	ExcluidoEm  *time.Time `json:"excluido_em"`
}

// ComissItem representa um pagamento/parcial de comissão.
type ComissItem struct {
	ID        int64      `json:"id"`
	ComisID   *int64     `json:"comis_id"`
	VlrPago   *float64   `json:"vlr_pago"`
	DtPgto    *time.Time `json:"dt_pgto"`
	ObsPgto   *string    `json:"obs_pgto"`
	CriadoEm  time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}
