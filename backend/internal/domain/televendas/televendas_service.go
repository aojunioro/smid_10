package televendas

import (
	"context"
	"fmt"
)

type TelevendasContatoService struct {
	contatoRepo TelevendasContatoRepository
}

func NewTelevendasContatoService(contatoRepo TelevendasContatoRepository) *TelevendasContatoService {
	return &TelevendasContatoService{
		contatoRepo: contatoRepo,
	}
}

type CreateContatoRequest struct {
	LeadID       *int64    `json:"lead_id"`
	Login        *string   `json:"login"`
	StatusID     *int64    `json:"status_id"`
	OrcamID      *int64    `json:"orcam_id"`
	PedID        *int64    `json:"ped_id"`
	VisitaID     *int64    `json:"visita_id"`
	HistoricoID  *int64    `json:"historico_id"`
	MidiaID      *int64    `json:"midia_id"`
	UniddID      *int64    `json:"unidd_id"`
	LoginRecep   *string   `json:"login_recep"`
	LoginRepre   *string   `json:"login_repre"`
	DtVisita     *string   `json:"dt_visita"`
	HrVisita     *string   `json:"hr_visita"`
	Cidade       *string   `json:"cidade"`
	Bairro       *string   `json:"bairro"`
	StatusVisitaID *int64  `json:"status_visita_id"`
	MotivoID     *int64    `json:"motivo_id"`
	Observacao    *string   `json:"observacao"`
}

type UpdateContatoRequest struct {
	LeadID       *int64    `json:"lead_id"`
	Login        *string   `json:"login"`
	StatusID     *int64    `json:"status_id"`
	OrcamID      *int64    `json:"orcam_id"`
	PedID        *int64    `json:"ped_id"`
	VisitaID     *int64    `json:"visita_id"`
	HistoricoID  *int64    `json:"historico_id"`
	MidiaID      *int64    `json:"midia_id"`
	UniddID      *int64    `json:"unidd_id"`
	LoginRecep   *string   `json:"login_recep"`
	LoginRepre   *string   `json:"login_repre"`
	DtVisita     *string   `json:"dt_visita"`
	HrVisita     *string   `json:"hr_visita"`
	Cidade       *string   `json:"cidade"`
	Bairro       *string   `json:"bairro"`
	StatusVisitaID *int64  `json:"status_visita_id"`
	MotivoID     *int64    `json:"motivo_id"`
	Observacao    *string   `json:"observacao"`
}

func (s *TelevendasContatoService) Create(ctx context.Context, req CreateContatoRequest) (*TelevendasContato, error) {
	contato := &TelevendasContato{
		LeadID:       req.LeadID,
		Login:        req.Login,
		StatusID:     req.StatusID,
		OrcamID:      req.OrcamID,
		PedID:        req.PedID,
		VisitaID:     req.VisitaID,
		HistoricoID:  req.HistoricoID,
		MidiaID:      req.MidiaID,
		UniddID:      req.UniddID,
		LoginRecep:   req.LoginRecep,
		LoginRepre:   req.LoginRepre,
		Cidade:       req.Cidade,
		Bairro:       req.Bairro,
		StatusVisitaID: req.StatusVisitaID,
		MotivoID:     req.MotivoID,
		Observacao:   req.Observacao,
	}

	if err := s.contatoRepo.Create(ctx, contato); err != nil {
		return nil, err
	}

	return contato, nil
}

func (s *TelevendasContatoService) GetByID(ctx context.Context, id int64) (*TelevendasContato, error) {
	contato, err := s.contatoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return contato, nil
}

func (s *TelevendasContatoService) List(ctx context.Context, limit, offset int, leadID *int64, login, statusID *string) ([]TelevendasContato, error) {
	var leadIDInt *int64
	if leadID != nil {
		leadIDInt = leadID
	}

	var loginStr *string
	if login != nil {
		loginStr = login
	}

	var statusIDInt *int64
	if statusID != nil {
		// Converter string para int64
		var id int64
		if _, err := fmt.Sscanf(*statusID, "%d", &id); err == nil {
			statusIDInt = &id
		}
	}

	opts := ListOptions{
		Limit:    limit,
		Offset:   offset,
		LeadID:   leadIDInt,
		Login:    loginStr,
		StatusID: statusIDInt,
	}

	contatos, err := s.contatoRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return contatos, nil
}

func (s *TelevendasContatoService) Update(ctx context.Context, id int64, req UpdateContatoRequest) (*TelevendasContato, error) {
	contato, err := s.contatoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.LeadID != nil {
		contato.LeadID = req.LeadID
	}
	if req.Login != nil {
		contato.Login = req.Login
	}
	if req.StatusID != nil {
		contato.StatusID = req.StatusID
	}
	if req.OrcamID != nil {
		contato.OrcamID = req.OrcamID
	}
	if req.PedID != nil {
		contato.PedID = req.PedID
	}
	if req.VisitaID != nil {
		contato.VisitaID = req.VisitaID
	}
	if req.HistoricoID != nil {
		contato.HistoricoID = req.HistoricoID
	}
	if req.MidiaID != nil {
		contato.MidiaID = req.MidiaID
	}
	if req.UniddID != nil {
		contato.UniddID = req.UniddID
	}
	if req.LoginRecep != nil {
		contato.LoginRecep = req.LoginRecep
	}
	if req.LoginRepre != nil {
		contato.LoginRepre = req.LoginRepre
	}
	if req.Cidade != nil {
		contato.Cidade = req.Cidade
	}
	if req.Bairro != nil {
		contato.Bairro = req.Bairro
	}
	if req.StatusVisitaID != nil {
		contato.StatusVisitaID = req.StatusVisitaID
	}
	if req.MotivoID != nil {
		contato.MotivoID = req.MotivoID
	}
	if req.Observacao != nil {
		contato.Observacao = req.Observacao
	}

	if err := s.contatoRepo.Update(ctx, contato); err != nil {
		return nil, err
	}

	return contato, nil
}

func (s *TelevendasContatoService) Delete(ctx context.Context, id int64) error {
	return s.contatoRepo.SoftDelete(ctx, id)
}
