package leads

import "time"

// Lead representa um potencial cliente.
type Lead struct {
	ID            int64      `json:"id"`
	Fone1         string     `json:"fone1"`
	StartTime     *time.Time `json:"starttime"`
	Fone2         *string    `json:"fone2"`
	Nome          string     `json:"nome"`
	Profissao     *string    `json:"profissao"`
	Idade         *int       `json:"idade"`
	Patologia     *string    `json:"patologia"`
	NomeAcomp     *string    `json:"nome_acomp"`
	ProfisAcomp   *string    `json:"profis_acomp"`
	IddAcomp      *int       `json:"idd_acomp"`
	PatoAcomp     *string    `json:"pato_acomp"`
	MidiaID       *int64     `json:"midia_id"`
	TentID        *int64     `json:"tent_id"`
	ContatoOK     string     `json:"contato_ok"`
	StatusID      int64      `json:"status_id"`
	UniddID       *int64     `json:"unidd_id"`
	MeioID        *int64     `json:"meio_id"`
	MotPendID     *int64     `json:"mot_pend_id"`
	MotPerdID     *int64     `json:"mot_perd_id"`
	Email         *string    `json:"email"`
	ObsCurtaLead  *string    `json:"obs_curta_lead"`
	Login         *string    `json:"login"`
	LoginRecep    *string    `json:"login_recep"`
	LoginSuper    *string    `json:"login_super"`
	CriadoEm      time.Time  `json:"criado_em"`
	AlteradoEm    time.Time  `json:"alterado_em"`
	ExcluidoEm    *time.Time `json:"excluido_em"`
}

// LeadStatus representa um status do lead no funil.
type LeadStatus struct {
	ID         int64  `json:"id"`
	SttLead    string `json:"stt_lead"`
	Cor        *string `json:"cor"`
	Ordem      *int   `json:"ordem"`
	Kanban     *string `json:"kanban"`
	SttInicial *string `json:"stt_inicial"`
	SttFinal   *string `json:"stt_final"`
	PermEdit   *string `json:"perm_edit"`
	PermDel    *string `json:"perm_del"`
}

// LeadMeio representa um canal técnico de entrada.
type LeadMeio struct {
	ID    int64  `json:"id"`
	Meio  string `json:"meio"`
	Ativo string `json:"ativo"`
}

// Midia representa uma campanha ou veículo publicitário.
type Midia struct {
	ID        int64     `json:"id"`
	UniddID   *string   `json:"unidd_id"`
	Midia     string    `json:"midia"`
	TipoID    *int64    `json:"tipo_id"`
	HoraIni   *string   `json:"hora_ini"`
	HoraFim   *string   `json:"hora_fim"`
	Ativa     string    `json:"ativa"`
	CreatedAt *string   `json:"created_at"`
}

// Endereco representa o endereço de um lead.
type Endereco struct {
	ID          int64  `json:"id"`
	LeadID      int64  `json:"lead_id"`
	CEP         *string `json:"CEP"`
	Rua         string `json:"rua"`
	Numero      string `json:"numero"`
	Complemento *string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Cidade      string `json:"cidade"`
	UF          string `json:"uf"`
	Referencias *string `json:"referencias"`
}
