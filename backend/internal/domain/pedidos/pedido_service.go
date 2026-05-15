package pedidos

import "context"

type PedidoService struct {
	pedidoRepo PedidoRepository
}

func NewPedidoService(pedidoRepo PedidoRepository) *PedidoService {
	return &PedidoService{
		pedidoRepo: pedidoRepo,
	}
}

type CreatePedidoRequest struct {
	LeadID        int64    `json:"lead_id"`
	NPed          *string  `json:"n_ped"`
	DtPed         *string  `json:"dt_ped"`
	DtPrev        *string  `json:"dt_prev"`
	DtQuit        *string  `json:"dt_quit"`
	DtEntr        *string  `json:"dt_entr"`
	StatusID      *int64   `json:"status_id"`
	CanalID       *int64   `json:"canal_id"`
	FpgtoID       *int64   `json:"fpgto_id"`
	CpgtoID       *int64   `json:"cpgto_id"`
	LoginRepre    *string  `json:"login_repre"`
	Login         *string  `json:"login"`
	TotalPed      *float64 `json:"total_ped"`
	EntradaPed    *float64 `json:"entrada_ped"`
	TaxaFinanceira *float64 `json:"taxa_financeira"`
	ValorLiquido  *float64 `json:"valor_liquido"`
	ObsPed        *string  `json:"obs_ped"`
	ObsPedGer     *string  `json:"obs_ped_ger"`
	ImgPed        *string  `json:"img_ped"`
}

type UpdatePedidoRequest struct {
	LeadID        *int64   `json:"lead_id"`
	NPed          *string  `json:"n_ped"`
	DtPed         *string  `json:"dt_ped"`
	DtPrev        *string  `json:"dt_prev"`
	DtQuit        *string  `json:"dt_quit"`
	DtEntr        *string  `json:"dt_entr"`
	StatusID      *int64   `json:"status_id"`
	CanalID       *int64   `json:"canal_id"`
	FpgtoID       *int64   `json:"fpgto_id"`
	CpgtoID       *int64   `json:"cpgto_id"`
	LoginRepre    *string  `json:"login_repre"`
	Login         *string  `json:"login"`
	TotalPed      *float64 `json:"total_ped"`
	EntradaPed    *float64 `json:"entrada_ped"`
	TaxaFinanceira *float64 `json:"taxa_financeira"`
	ValorLiquido  *float64 `json:"valor_liquido"`
	ObsPed        *string  `json:"obs_ped"`
	ObsPedGer     *string  `json:"obs_ped_ger"`
	ImgPed        *string  `json:"img_ped"`
}

func (s *PedidoService) Create(ctx context.Context, req CreatePedidoRequest) (*Pedido, error) {
	pedido := &Pedido{
		LeadID:        req.LeadID,
		NPed:          req.NPed,
		StatusID:      req.StatusID,
		CanalID:       req.CanalID,
		FpgtoID:       req.FpgtoID,
		CpgtoID:       req.CpgtoID,
		LoginRepre:    req.LoginRepre,
		Login:         req.Login,
		TotalPed:      req.TotalPed,
		EntradaPed:    req.EntradaPed,
		TaxaFinanceira: req.TaxaFinanceira,
		ValorLiquido:  req.ValorLiquido,
		ObsPed:        req.ObsPed,
		ObsPedGer:     req.ObsPedGer,
		ImgPed:        req.ImgPed,
	}

	if err := s.pedidoRepo.Create(ctx, pedido); err != nil {
		return nil, err
	}

	return pedido, nil
}

func (s *PedidoService) GetByID(ctx context.Context, id int64) (*Pedido, error) {
	pedido, err := s.pedidoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return pedido, nil
}

func (s *PedidoService) List(ctx context.Context, limit, offset int, leadID, statusID *int64, loginRepre, dtPed *string) ([]Pedido, error) {
	opts := ListOptions{
		Limit:      limit,
		Offset:     offset,
		LeadID:     leadID,
		StatusID:   statusID,
		LoginRepre: loginRepre,
		DtPed:      dtPed,
	}

	pedidos, err := s.pedidoRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return pedidos, nil
}

func (s *PedidoService) Update(ctx context.Context, id int64, req UpdatePedidoRequest) (*Pedido, error) {
	pedido, err := s.pedidoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.LeadID != nil {
		pedido.LeadID = *req.LeadID
	}
	if req.NPed != nil {
		pedido.NPed = req.NPed
	}
	if req.StatusID != nil {
		pedido.StatusID = req.StatusID
	}
	if req.CanalID != nil {
		pedido.CanalID = req.CanalID
	}
	if req.FpgtoID != nil {
		pedido.FpgtoID = req.FpgtoID
	}
	if req.CpgtoID != nil {
		pedido.CpgtoID = req.CpgtoID
	}
	if req.LoginRepre != nil {
		pedido.LoginRepre = req.LoginRepre
	}
	if req.Login != nil {
		pedido.Login = req.Login
	}
	if req.TotalPed != nil {
		pedido.TotalPed = req.TotalPed
	}
	if req.EntradaPed != nil {
		pedido.EntradaPed = req.EntradaPed
	}
	if req.TaxaFinanceira != nil {
		pedido.TaxaFinanceira = req.TaxaFinanceira
	}
	if req.ValorLiquido != nil {
		pedido.ValorLiquido = req.ValorLiquido
	}
	if req.ObsPed != nil {
		pedido.ObsPed = req.ObsPed
	}
	if req.ObsPedGer != nil {
		pedido.ObsPedGer = req.ObsPedGer
	}
	if req.ImgPed != nil {
		pedido.ImgPed = req.ImgPed
	}

	if err := s.pedidoRepo.Update(ctx, pedido); err != nil {
		return nil, err
	}

	return pedido, nil
}

func (s *PedidoService) Delete(ctx context.Context, id int64) error {
	return s.pedidoRepo.SoftDelete(ctx, id)
}
