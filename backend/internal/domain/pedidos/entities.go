package pedidos

import "time"

// Pedido representa um pedido de venda (PedidosRepre).
type Pedido struct {
	ID            int64      `json:"id"`
	LeadID        int64      `json:"lead_id"`
	NPed          *string    `json:"n_ped"`
	DtPed         *time.Time `json:"dt_ped"`
	DtPrev        *time.Time `json:"dt_prev"`
	DtQuit        *time.Time `json:"dt_quit"`
	DtEntr        *time.Time `json:"dt_entr"`
	StatusID      *int64     `json:"status_id"`
	CanalID       *int64     `json:"canal_id"`
	FpgtoID       *int64     `json:"fpgto_id"`
	CpgtoID       *int64     `json:"cpgto_id"`
	LoginRepre    *string    `json:"login_repre"`
	Login         *string    `json:"login"`
	TotalPed      *float64   `json:"total_ped"`
	EntradaPed    *float64   `json:"entrada_ped"`
	TaxaFinanceira *float64   `json:"taxa_financeira"`
	ValorLiquido  *float64   `json:"valor_liquido"`
	ObsPed        *string    `json:"obs_ped"`
	ObsPedGer     *string    `json:"obs_ped_ger"`
	ImgPed        *string    `json:"img_ped"`
	CriadoEm      time.Time  `json:"criado_em"`
	AlteradoEm    *time.Time `json:"alterado_em"`
	ExcluidoEm    *time.Time `json:"excluido_em"`
}

// PedProdItem representa um item do pedido (PedProdItem).
type PedProdItem struct {
	ID            int64      `json:"id"`
	PedID         int64      `json:"ped_id"`
	Produto       *string    `json:"produto"`
	Quantidade    *int       `json:"quantidade"`
	ValorUnitario *float64   `json:"valor_unitario"`
	Desconto      *float64   `json:"desconto"`
	CriadoEm      time.Time  `json:"criado_em"`
	AlteradoEm    *time.Time `json:"alterado_em"`
	ExcluidoEm    *time.Time `json:"excluido_em"`
}

// PedStatus representa um status de pedido.
type PedStatus struct {
	ID         int64      `json:"id"`
	Status     string     `json:"status"`
	CategoriaID *int64    `json:"categoria_id"`
	Sistema    *string    `json:"sistema"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// PedStatusCategoria representa uma categoria funcional de status de pedido.
type PedStatusCategoria struct {
	ID         int64      `json:"id"`
	Codigo     string     `json:"codigo"`
	Nome       string     `json:"nome"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// PedFormaPagamento representa uma forma de pagamento.
type PedFormaPagamento struct {
	ID         int64      `json:"id"`
	Fpgto      string     `json:"fpgto"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// PedCondicaoPagamento representa uma condição de pagamento.
type PedCondicaoPagamento struct {
	ID         int64      `json:"id"`
	Cpgto      string     `json:"cpgto"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// PedCanal representa um canal de venda.
type PedCanal struct {
	ID         int64      `json:"id"`
	Canal      string     `json:"canal"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}
