package suporte

import "time"

// Suporte representa um chamado/atendimento de suporte pós-venda.
type Suporte struct {
	ID          int64      `json:"id"`
	PedID       *int64     `json:"ped_id"`
	StatusID    *int64     `json:"status_id"`
	FoneSup     *string    `json:"fone_sup"`
	SolicitID   *int64     `json:"solicit_id"`
	DepartID    *int64     `json:"depart_id"`
	Login       *string    `json:"login"`
	AtribLogin  *string    `json:"atrib_login"`
	Prioridade  *string    `json:"prioridade"`
	DtSup       *time.Time `json:"dt_sup"`
	DtLimit     *time.Time `json:"dt_limit"`
	RelatoCli   *string    `json:"relato_cli"`
	DtResol     *time.Time `json:"dt_resol"`
	RelatoTec   *string    `json:"relato_tec"`
	ImgOrdem    *string    `json:"img_ordem"`
	CriadoEm    time.Time  `json:"criado_em"`
	AlteradoEm  *time.Time `json:"alterado_em"`
	ExcluidoEm  *time.Time `json:"excluido_em"`
}
