package suporte

import (
	"context"
	"fmt"
)

type SuporteService struct {
	suporteRepo SuporteRepository
}

func NewSuporteService(suporteRepo SuporteRepository) *SuporteService {
	return &SuporteService{
		suporteRepo: suporteRepo,
	}
}

type CreateSuporteRequest struct {
	PedID      *int64  `json:"ped_id"`
	StatusID   *int64  `json:"status_id"`
	FoneSup    *string `json:"fone_sup"`
	SolicitID  *int64  `json:"solicit_id"`
	DepartID   *int64  `json:"depart_id"`
	Login      *string `json:"login"`
	AtribLogin *string `json:"atrib_login"`
	Prioridade *string `json:"prioridade"`
	DtSup      *string `json:"dt_sup"`
	DtLimit    *string `json:"dt_limit"`
	RelatoCli  *string `json:"relato_cli"`
	DtResol    *string `json:"dt_resol"`
	RelatoTec  *string `json:"relato_tec"`
	ImgOrdem   *string `json:"img_ordem"`
}

type UpdateSuporteRequest struct {
	PedID      *int64  `json:"ped_id"`
	StatusID   *int64  `json:"status_id"`
	FoneSup    *string `json:"fone_sup"`
	SolicitID  *int64  `json:"solicit_id"`
	DepartID   *int64  `json:"depart_id"`
	Login      *string `json:"login"`
	AtribLogin *string `json:"atrib_login"`
	Prioridade *string `json:"prioridade"`
	DtSup      *string `json:"dt_sup"`
	DtLimit    *string `json:"dt_limit"`
	RelatoCli  *string `json:"relato_cli"`
	DtResol    *string `json:"dt_resol"`
	RelatoTec  *string `json:"relato_tec"`
	ImgOrdem   *string `json:"img_ordem"`
}

func (s *SuporteService) Create(ctx context.Context, req CreateSuporteRequest) (*Suporte, error) {
	suporte := &Suporte{
		PedID:      req.PedID,
		StatusID:   req.StatusID,
		FoneSup:    req.FoneSup,
		SolicitID:  req.SolicitID,
		DepartID:   req.DepartID,
		Login:      req.Login,
		AtribLogin: req.AtribLogin,
		Prioridade: req.Prioridade,
		RelatoCli:  req.RelatoCli,
		RelatoTec:  req.RelatoTec,
		ImgOrdem:   req.ImgOrdem,
	}

	if err := s.suporteRepo.Create(ctx, suporte); err != nil {
		return nil, err
	}

	return suporte, nil
}

func (s *SuporteService) GetByID(ctx context.Context, id int64) (*Suporte, error) {
	suporte, err := s.suporteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return suporte, nil
}

func (s *SuporteService) List(ctx context.Context, limit, offset int, pedID *int64, login, atribLogin, statusID *string) ([]Suporte, error) {
	var pedIDInt *int64
	if pedID != nil {
		pedIDInt = pedID
	}

	var loginStr *string
	if login != nil {
		loginStr = login
	}

	var atribLoginStr *string
	if atribLogin != nil {
		atribLoginStr = atribLogin
	}

	var statusIDInt *int64
	if statusID != nil {
		var id int64
		if _, err := fmt.Sscanf(*statusID, "%d", &id); err == nil {
			statusIDInt = &id
		}
	}

	opts := ListOptions{
		Limit:      limit,
		Offset:     offset,
		PedID:      pedIDInt,
		Login:      loginStr,
		AtribLogin: atribLoginStr,
		StatusID:   statusIDInt,
	}

	suportes, err := s.suporteRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return suportes, nil
}

func (s *SuporteService) Update(ctx context.Context, id int64, req UpdateSuporteRequest) (*Suporte, error) {
	suporte, err := s.suporteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.PedID != nil {
		suporte.PedID = req.PedID
	}
	if req.StatusID != nil {
		suporte.StatusID = req.StatusID
	}
	if req.FoneSup != nil {
		suporte.FoneSup = req.FoneSup
	}
	if req.SolicitID != nil {
		suporte.SolicitID = req.SolicitID
	}
	if req.DepartID != nil {
		suporte.DepartID = req.DepartID
	}
	if req.Login != nil {
		suporte.Login = req.Login
	}
	if req.AtribLogin != nil {
		suporte.AtribLogin = req.AtribLogin
	}
	if req.Prioridade != nil {
		suporte.Prioridade = req.Prioridade
	}
	if req.RelatoCli != nil {
		suporte.RelatoCli = req.RelatoCli
	}
	if req.RelatoTec != nil {
		suporte.RelatoTec = req.RelatoTec
	}
	if req.ImgOrdem != nil {
		suporte.ImgOrdem = req.ImgOrdem
	}

	if err := s.suporteRepo.Update(ctx, suporte); err != nil {
		return nil, err
	}

	return suporte, nil
}

func (s *SuporteService) Delete(ctx context.Context, id int64) error {
	return s.suporteRepo.SoftDelete(ctx, id)
}
