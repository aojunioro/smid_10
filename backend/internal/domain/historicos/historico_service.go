package historicos

import "context"

type HistoricoService struct {
	historicoRepo HistoricoRepository
}

func NewHistoricoService(historicoRepo HistoricoRepository) *HistoricoService {
	return &HistoricoService{
		historicoRepo: historicoRepo,
	}
}

type CreateHistoricoRequest struct {
	VisID      *int64  `json:"vis_id"`
	LeadID     *int64  `json:"lead_id"`
	MotivoID   *int64  `json:"motivo_id"`
	OcorridoID *int64  `json:"ocorrido_id"`
	Hist       *string `json:"hist"`
	FotoHist   *string `json:"foto_hist"`
	Login      *string `json:"login"`
}

type UpdateHistoricoRequest struct {
	VisID      *int64  `json:"vis_id"`
	LeadID     *int64  `json:"lead_id"`
	MotivoID   *int64  `json:"motivo_id"`
	OcorridoID *int64  `json:"ocorrido_id"`
	Hist       *string `json:"hist"`
	FotoHist   *string `json:"foto_hist"`
	Login      *string `json:"login"`
}

func (s *HistoricoService) Create(ctx context.Context, req CreateHistoricoRequest) (*HistoricoRepre, error) {
	historico := &HistoricoRepre{
		VisID:      req.VisID,
		LeadID:     req.LeadID,
		MotivoID:   req.MotivoID,
		OcorridoID: req.OcorridoID,
		Hist:       req.Hist,
		FotoHist:   req.FotoHist,
		Login:      req.Login,
	}

	if err := s.historicoRepo.Create(ctx, historico); err != nil {
		return nil, err
	}

	return historico, nil
}

func (s *HistoricoService) GetByID(ctx context.Context, id int64) (*HistoricoRepre, error) {
	historico, err := s.historicoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return historico, nil
}

func (s *HistoricoService) List(ctx context.Context, limit, offset int, visID, leadID *int64, login *string, motivoID, ocorridoID *int64) ([]HistoricoRepre, error) {
	opts := ListOptions{
		Limit:      limit,
		Offset:     offset,
		VisID:      visID,
		LeadID:     leadID,
		Login:      login,
		MotivoID:   motivoID,
		OcorridoID: ocorridoID,
	}

	historicos, err := s.historicoRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return historicos, nil
}

func (s *HistoricoService) Update(ctx context.Context, id int64, req UpdateHistoricoRequest) (*HistoricoRepre, error) {
	historico, err := s.historicoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.VisID != nil {
		historico.VisID = req.VisID
	}
	if req.LeadID != nil {
		historico.LeadID = req.LeadID
	}
	if req.MotivoID != nil {
		historico.MotivoID = req.MotivoID
	}
	if req.OcorridoID != nil {
		historico.OcorridoID = req.OcorridoID
	}
	if req.Hist != nil {
		historico.Hist = req.Hist
	}
	if req.FotoHist != nil {
		historico.FotoHist = req.FotoHist
	}
	if req.Login != nil {
		historico.Login = req.Login
	}

	if err := s.historicoRepo.Update(ctx, historico); err != nil {
		return nil, err
	}

	return historico, nil
}

func (s *HistoricoService) Delete(ctx context.Context, id int64) error {
	return s.historicoRepo.SoftDelete(ctx, id)
}
