package historicos

import "time"

// HistoricoRepre representa um registro de evento histórico de visita.
type HistoricoRepre struct {
	ID          int64      `json:"id"`
	VisID       *int64     `json:"vis_id"`
	LeadID      *int64     `json:"lead_id"`
	MotivoID    *int64     `json:"motivo_id"`
	OcorridoID  *int64     `json:"ocorrido_id"`
	Hist        *string    `json:"hist"`
	FotoHist    *string    `json:"foto_hist"`
	Login       *string    `json:"login"`
	CriadoEm    time.Time  `json:"criado_em"`
	AlteradoEm  *time.Time `json:"alterado_em"`
	ExcluidoEm  *time.Time `json:"excluido_em"`
}

// HistoricoMotivo representa um catálogo de motivos.
type HistoricoMotivo struct {
	ID         int64      `json:"id"`
	Motivo     string     `json:"motivo"`
	Cor        *string    `json:"cor"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// HistoricoOcorrido representa um catálogo de ocorridos.
type HistoricoOcorrido struct {
	ID         int64      `json:"id"`
	Ocorrido   string     `json:"ocorrido"`
	Cor        *string    `json:"cor"`
	Ordem      *int       `json:"ordem"`
	Responsavel *string   `json:"responsavel"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}
