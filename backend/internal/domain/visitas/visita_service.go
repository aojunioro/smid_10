package visitas

import "context"

type VisitaService struct {
	visitaRepo VisitaRepository
}

func NewVisitaService(visitaRepo VisitaRepository) *VisitaService {
	return &VisitaService{
		visitaRepo: visitaRepo,
	}
}

type CreateVisitaRequest struct {
	LeadID     int64   `json:"lead_id"`
	StatusID   *int64  `json:"status_id"`
	LoginRecep *string `json:"login_recep"`
	LoginRepre *string `json:"login_repre"`
	DtVisita   *string `json:"dt_visita"`
	HrVisita   *string `json:"hr_visita"`
	Confirm    *string `json:"confirm"`
	LoginConf  *string `json:"login_conf"`
	DtConfirm  *string `json:"dt_confirm"`
	Interesse  *string `json:"interesse"`
	HistFeito  *string `json:"hist_feito"`
	PosFeito   *string `json:"pos_feito"`
	SttsLead   *int64  `json:"stts_lead"`
}

type UpdateVisitaRequest struct {
	LeadID     *int64  `json:"lead_id"`
	StatusID   *int64  `json:"status_id"`
	LoginRecep *string `json:"login_recep"`
	LoginRepre *string `json:"login_repre"`
	DtVisita   *string `json:"dt_visita"`
	HrVisita   *string `json:"hr_visita"`
	Confirm    *string `json:"confirm"`
	LoginConf  *string `json:"login_conf"`
	DtConfirm  *string `json:"dt_confirm"`
	Interesse  *string `json:"interesse"`
	HistFeito  *string `json:"hist_feito"`
	PosFeito   *string `json:"pos_feito"`
	SttsLead   *int64  `json:"stts_lead"`
}

func (s *VisitaService) Create(ctx context.Context, req CreateVisitaRequest) (*Visita, error) {
	visita := &Visita{
		LeadID:     req.LeadID,
		StatusID:   req.StatusID,
		LoginRecep: req.LoginRecep,
		LoginRepre: req.LoginRepre,
		DtVisita:   req.DtVisita,
		HrVisita:   req.HrVisita,
		Confirm:    req.Confirm,
		LoginConf:  req.LoginConf,
		Interesse:  req.Interesse,
		HistFeito:  req.HistFeito,
		PosFeito:   req.PosFeito,
		SttsLead:   req.SttsLead,
	}

	if err := s.visitaRepo.Create(ctx, visita); err != nil {
		return nil, err
	}

	return visita, nil
}

func (s *VisitaService) GetByID(ctx context.Context, id int64) (*Visita, error) {
	visita, err := s.visitaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return visita, nil
}

func (s *VisitaService) List(ctx context.Context, limit, offset int, leadID *int64, loginRepre *string, loginRecep *string, statusID *int64, dtVisita *string) ([]Visita, error) {
	opts := ListOptions{
		Limit:      limit,
		Offset:     offset,
		LeadID:     leadID,
		LoginRepre: loginRepre,
		LoginRecep: loginRecep,
		StatusID:   statusID,
		DtVisita:   dtVisita,
	}

	visitas, err := s.visitaRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return visitas, nil
}

func (s *VisitaService) Update(ctx context.Context, id int64, req UpdateVisitaRequest) (*Visita, error) {
	visita, err := s.visitaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.LeadID != nil {
		visita.LeadID = *req.LeadID
	}
	if req.StatusID != nil {
		visita.StatusID = req.StatusID
	}
	if req.LoginRecep != nil {
		visita.LoginRecep = req.LoginRecep
	}
	if req.LoginRepre != nil {
		visita.LoginRepre = req.LoginRepre
	}
	if req.DtVisita != nil {
		visita.DtVisita = req.DtVisita
	}
	if req.HrVisita != nil {
		visita.HrVisita = req.HrVisita
	}
	if req.Confirm != nil {
		visita.Confirm = req.Confirm
	}
	if req.LoginConf != nil {
		visita.LoginConf = req.LoginConf
	}
	// DtConfirm não é atualizado no Update (campo de auditoria)
	if req.Interesse != nil {
		visita.Interesse = req.Interesse
	}
	if req.HistFeito != nil {
		visita.HistFeito = req.HistFeito
	}
	if req.PosFeito != nil {
		visita.PosFeito = req.PosFeito
	}
	if req.SttsLead != nil {
		visita.SttsLead = req.SttsLead
	}

	if err := s.visitaRepo.Update(ctx, visita); err != nil {
		return nil, err
	}

	return visita, nil
}

func (s *VisitaService) Delete(ctx context.Context, id int64) error {
	return s.visitaRepo.SoftDelete(ctx, id)
}
