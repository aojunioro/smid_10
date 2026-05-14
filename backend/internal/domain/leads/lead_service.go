package leads

import "context"

type LeadService struct {
	leadRepo LeadRepository
}

func NewLeadService(leadRepo LeadRepository) *LeadService {
	return &LeadService{
		leadRepo: leadRepo,
	}
}

type CreateLeadRequest struct {
	Fone1         string     `json:"fone1"`
	StartTime     *string    `json:"starttime"`
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
}

type UpdateLeadRequest struct {
	Fone1         *string    `json:"fone1"`
	StartTime     *string    `json:"starttime"`
	Fone2         *string    `json:"fone2"`
	Nome          *string    `json:"nome"`
	Profissao     *string    `json:"profissao"`
	Idade         *int       `json:"idade"`
	Patologia     *string    `json:"patologia"`
	NomeAcomp     *string    `json:"nome_acomp"`
	ProfisAcomp   *string    `json:"profis_acomp"`
	IddAcomp      *int       `json:"idd_acomp"`
	PatoAcomp     *string    `json:"pato_acomp"`
	MidiaID       *int64     `json:"midia_id"`
	TentID        *int64     `json:"tent_id"`
	ContatoOK     *string    `json:"contato_ok"`
	StatusID      *int64     `json:"status_id"`
	UniddID       *int64     `json:"unidd_id"`
	MeioID        *int64     `json:"meio_id"`
	MotPendID     *int64     `json:"mot_pend_id"`
	MotPerdID     *int64     `json:"mot_perd_id"`
	Email         *string    `json:"email"`
	ObsCurtaLead  *string    `json:"obs_curta_lead"`
	Login         *string    `json:"login"`
	LoginRecep    *string    `json:"login_recep"`
	LoginSuper    *string    `json:"login_super"`
}

func (s *LeadService) Create(ctx context.Context, req CreateLeadRequest) (*Lead, error) {
	lead := &Lead{
		Fone1:         req.Fone1,
		Fone2:         req.Fone2,
		Nome:          req.Nome,
		Profissao:     req.Profissao,
		Idade:         req.Idade,
		Patologia:     req.Patologia,
		NomeAcomp:     req.NomeAcomp,
		ProfisAcomp:   req.ProfisAcomp,
		IddAcomp:      req.IddAcomp,
		PatoAcomp:     req.PatoAcomp,
		MidiaID:       req.MidiaID,
		TentID:        req.TentID,
		ContatoOK:     req.ContatoOK,
		StatusID:      req.StatusID,
		UniddID:       req.UniddID,
		MeioID:        req.MeioID,
		MotPendID:     req.MotPendID,
		MotPerdID:     req.MotPerdID,
		Email:         req.Email,
		ObsCurtaLead:  req.ObsCurtaLead,
		Login:         req.Login,
		LoginRecep:    req.LoginRecep,
		LoginSuper:    req.LoginSuper,
	}

	if err := s.leadRepo.Create(ctx, lead); err != nil {
		return nil, err
	}

	return lead, nil
}

func (s *LeadService) GetByID(ctx context.Context, id int64) (*Lead, error) {
	lead, err := s.leadRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return lead, nil
}

func (s *LeadService) List(ctx context.Context, limit, offset int) ([]Lead, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	leads, err := s.leadRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return leads, nil
}

func (s *LeadService) Update(ctx context.Context, id int64, req UpdateLeadRequest) (*Lead, error) {
	lead, err := s.leadRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Fone1 != nil {
		lead.Fone1 = *req.Fone1
	}
	if req.Fone2 != nil {
		lead.Fone2 = req.Fone2
	}
	if req.Nome != nil {
		lead.Nome = *req.Nome
	}
	if req.Profissao != nil {
		lead.Profissao = req.Profissao
	}
	if req.Idade != nil {
		lead.Idade = req.Idade
	}
	if req.Patologia != nil {
		lead.Patologia = req.Patologia
	}
	if req.NomeAcomp != nil {
		lead.NomeAcomp = req.NomeAcomp
	}
	if req.ProfisAcomp != nil {
		lead.ProfisAcomp = req.ProfisAcomp
	}
	if req.IddAcomp != nil {
		lead.IddAcomp = req.IddAcomp
	}
	if req.PatoAcomp != nil {
		lead.PatoAcomp = req.PatoAcomp
	}
	if req.MidiaID != nil {
		lead.MidiaID = req.MidiaID
	}
	if req.TentID != nil {
		lead.TentID = req.TentID
	}
	if req.ContatoOK != nil {
		lead.ContatoOK = *req.ContatoOK
	}
	if req.StatusID != nil {
		lead.StatusID = *req.StatusID
	}
	if req.UniddID != nil {
		lead.UniddID = req.UniddID
	}
	if req.MeioID != nil {
		lead.MeioID = req.MeioID
	}
	if req.MotPendID != nil {
		lead.MotPendID = req.MotPendID
	}
	if req.MotPerdID != nil {
		lead.MotPerdID = req.MotPerdID
	}
	if req.Email != nil {
		lead.Email = req.Email
	}
	if req.ObsCurtaLead != nil {
		lead.ObsCurtaLead = req.ObsCurtaLead
	}
	if req.Login != nil {
		lead.Login = req.Login
	}
	if req.LoginRecep != nil {
		lead.LoginRecep = req.LoginRecep
	}
	if req.LoginSuper != nil {
		lead.LoginSuper = req.LoginSuper
	}

	if err := s.leadRepo.Update(ctx, lead); err != nil {
		return nil, err
	}

	return lead, nil
}
