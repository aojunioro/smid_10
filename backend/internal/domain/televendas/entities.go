package televendas

import "time"

// TelevendasContato representa um registro de contato realizado pelo operador.
type TelevendasContato struct {
	ID           int64      `json:"id"`
	LeadID       *int64     `json:"lead_id"`
	Login        *string    `json:"login"`
	StatusID     *int64     `json:"status_id"`
	OrcamID      *int64     `json:"orcam_id"`
	PedID        *int64     `json:"ped_id"`
	VisitaID     *int64     `json:"visita_id"`
	HistoricoID  *int64     `json:"historico_id"`
	MidiaID      *int64     `json:"midia_id"`
	UniddID      *int64     `json:"unidd_id"`
	LoginRecep   *string    `json:"login_recep"`
	LoginRepre   *string    `json:"login_repre"`
	DtVisita     *time.Time `json:"dt_visita"`
	HrVisita     *time.Time `json:"hr_visita"`
	Cidade       *string    `json:"cidade"`
	Bairro       *string    `json:"bairro"`
	StatusVisitaID *int64   `json:"status_visita_id"`
	MotivoID     *int64     `json:"motivo_id"`
	Observacao    *string    `json:"observacao"`
	CriadoEm     time.Time  `json:"criado_em"`
	AlteradoEm   *time.Time `json:"alterado_em"`
	ExcluidoEm   *time.Time `json:"excluido_em"`
}

// TelevendasStatus representa um status do contato Televendas.
type TelevendasStatus struct {
	ID         int64      `json:"id"`
	SttsTele   string     `json:"stts_tele"`
	Ordem      *int       `json:"ordem"`
	Cor        *string    `json:"cor"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// TeleCoeficiente representa coeficientes comerciais/financeiros do canal.
type TeleCoeficiente struct {
	ID          int64      `json:"id"`
	Nome        string     `json:"nome"`
	Coeficiente *float64   `json:"coeficiente"`
	CriadoEm    time.Time  `json:"criado_em"`
	AlteradoEm  *time.Time `json:"alterado_em"`
	ExcluidoEm  *time.Time `json:"excluido_em"`
}

// Orcam representa uma proposta comercial Televendas.
type Orcam struct {
	ID          int64      `json:"id"`
	LeadID      *int64     `json:"lead_id"`
	Login       *string    `json:"login"`
	StatusID    *int64     `json:"status_id"`
	Temperatura *string    `json:"temperatura"`
	DtOrcam     *time.Time `json:"dt_orcam"`
	VlrTotal    *float64   `json:"vlr_total"`
	ObsOrcam    *string    `json:"obs_orcam"`
	PedID       *int64     `json:"ped_id"`
	CriadoEm    time.Time  `json:"criado_em"`
	AlteradoEm  *time.Time `json:"alterado_em"`
	ExcluidoEm  *time.Time `json:"excluido_em"`
}

// OrcamStatus representa um status de orçamento.
type OrcamStatus struct {
	ID         int64      `json:"id"`
	Status     string     `json:"status"`
	Ordem      *int       `json:"ordem"`
	Cor        *string    `json:"cor"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// OrcamProdItem representa itens do orçamento.
type OrcamProdItem struct {
	ID            int64      `json:"id"`
	OrcamID       int64      `json:"orcam_id"`
	ProdOrcamID   *int64     `json:"prod_orcam_id"`
	MedID         *int64     `json:"med_id"`
	QtddItem      *int       `json:"qtdd_item"`
	VlrItem       *float64   `json:"vlr_item"`
	VlrTotalItem  *float64   `json:"vlr_total_item"`
	CriadoEm      time.Time  `json:"criado_em"`
	AlteradoEm    *time.Time `json:"alterado_em"`
	ExcluidoEm    *time.Time `json:"excluido_em"`
}
