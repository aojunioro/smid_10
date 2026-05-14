package visitas

import "time"

// Visita representa uma visita agendada para um lead.
type Visita struct {
	ID            int64      `json:"id"`
	LeadID        int64      `json:"lead_id"`
	StatusID      *int64     `json:"status_id"`
	LoginRecep    *string    `json:"login_recep"`
	LoginRepre    *string    `json:"login_repre"`
	DtVisita      *string    `json:"dt_visita"` // YYYY-MM-DD
	HrVisita      *string    `json:"hr_visita"` // HH:MM:SS
	Confirm       *string    `json:"confirm"`   // S/N
	LoginConf     *string    `json:"login_conf"`
	DtConfirm     *time.Time `json:"dt_confirm"`
	Interesse     *string    `json:"interesse"`
	HistFeito     *string    `json:"hist_feito"`
	PosFeito      *string    `json:"pos_feito"`
	SttsLead      *int64     `json:"stts_lead"`
	CriadoEm      time.Time  `json:"criado_em"`
	AlteradoEm    *time.Time `json:"alterado_em"`
	ExcluidoEm    *time.Time `json:"excluido_em"`
}

// VisitaCheckinEvent representa um evento GPS de check-in/check-out.
type VisitaCheckinEvent struct {
	ID          int64      `json:"id"`
	VisID       int64      `json:"vis_id"`
	LoginRepre  *string    `json:"login_repre"`
	StatusID    *int64     `json:"status_id"`
	Tipo        *string    `json:"tipo"` // CHECK_IN|CHECK_OUT
	Lat         *float64   `json:"lat"`
	Lng         *float64   `json:"lng"`
	AccuracyM   *float64   `json:"accuracy_m"`
	Permission  *string    `json:"permission"` // S/N
	ErroCodigo  *string    `json:"erro_codigo"`
	ErroMsg     *string    `json:"erro_msg"`
	CapturadoEm *time.Time `json:"capturado_em"`
	CriadoEm    time.Time  `json:"criado_em"`
	AlteradoEm  *time.Time `json:"alterado_em"`
}

// PosVisita representa uma avaliação pós-visita.
type PosVisita struct {
	ID         int64      `json:"id"`
	VisID      int64      `json:"vis_id"`
	LeadID     *int64     `json:"lead_id"`
	Login      *string    `json:"login"`
	NotaRepre  *int       `json:"nota_repre"`
	NotaProd   *int       `json:"nota_prod"`
	NotaEmpre  *int       `json:"nota_empre"`
	Visitado   *string    `json:"visitado"` // S/N
	Pontual    *string    `json:"pontual"`  // S/N
	Jaleco     *string    `json:"jaleco"`    // S/N
	Adquiriu   *string    `json:"adquiriu"`  // S/N
	Obs        *string    `json:"obs"`
	CriadoEm   time.Time  `json:"criado_em"`
}
