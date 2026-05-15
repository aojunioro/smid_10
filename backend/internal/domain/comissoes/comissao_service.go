package comissoes

import "context"

type ComissaoService struct {
	comissaoRepo  ComissaoRepository
	comissItemRepo ComissItemRepository
}

func NewComissaoService(comissaoRepo ComissaoRepository, comissItemRepo ComissItemRepository) *ComissaoService {
	return &ComissaoService{
		comissaoRepo:  comissaoRepo,
		comissItemRepo: comissItemRepo,
	}
}

type CreateComissaoRequest struct {
	PedID      *int64   `json:"ped_id"`
	VlrComissao *float64 `json:"vlr_comissao"`
	TotalPago  *float64 `json:"total_pago"`
	VlrSaldo   *float64 `json:"vlr_saldo"`
	DtPrevista *string  `json:"dt_prevista"`
	SttComis   *string  `json:"stt_comis"`
	Observacao *string  `json:"observacao"`
}

type UpdateComissaoRequest struct {
	PedID      *int64   `json:"ped_id"`
	VlrComissao *float64 `json:"vlr_comissao"`
	TotalPago  *float64 `json:"total_pago"`
	VlrSaldo   *float64 `json:"vlr_saldo"`
	DtPrevista *string  `json:"dt_prevista"`
	SttComis   *string  `json:"stt_comis"`
	Observacao *string  `json:"observacao"`
}

type CreateComissItemRequest struct {
	ComisID *int64  `json:"comis_id"`
	VlrPago *float64 `json:"vlr_pago"`
	DtPgto  *string `json:"dt_pgto"`
	ObsPgto *string `json:"obs_pgto"`
}

type UpdateComissItemRequest struct {
	ComisID *int64  `json:"comis_id"`
	VlrPago *float64 `json:"vlr_pago"`
	DtPgto  *string `json:"dt_pgto"`
	ObsPgto *string `json:"obs_pgto"`
}

// Comissao methods
func (s *ComissaoService) CreateComissao(ctx context.Context, req CreateComissaoRequest) (*Comissao, error) {
	comissao := &Comissao{
		PedID:      req.PedID,
		VlrComissao: req.VlrComissao,
		TotalPago:  req.TotalPago,
		VlrSaldo:   req.VlrSaldo,
		SttComis:   req.SttComis,
		Observacao: req.Observacao,
	}

	if err := s.comissaoRepo.Create(ctx, comissao); err != nil {
		return nil, err
	}

	return comissao, nil
}

func (s *ComissaoService) GetComissaoByID(ctx context.Context, id int64) (*Comissao, error) {
	comissao, err := s.comissaoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return comissao, nil
}

func (s *ComissaoService) ListComissoes(ctx context.Context, limit, offset int, pedID *int64, sttComis *string) ([]Comissao, error) {
	opts := ListOptions{
		Limit:    limit,
		Offset:   offset,
		PedID:    pedID,
		SttComis: sttComis,
	}

	comissoes, err := s.comissaoRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return comissoes, nil
}

func (s *ComissaoService) UpdateComissao(ctx context.Context, id int64, req UpdateComissaoRequest) (*Comissao, error) {
	comissao, err := s.comissaoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.PedID != nil {
		comissao.PedID = req.PedID
	}
	if req.VlrComissao != nil {
		comissao.VlrComissao = req.VlrComissao
	}
	if req.TotalPago != nil {
		comissao.TotalPago = req.TotalPago
	}
	if req.VlrSaldo != nil {
		comissao.VlrSaldo = req.VlrSaldo
	}
	if req.SttComis != nil {
		comissao.SttComis = req.SttComis
	}
	if req.Observacao != nil {
		comissao.Observacao = req.Observacao
	}

	if err := s.comissaoRepo.Update(ctx, comissao); err != nil {
		return nil, err
	}

	return comissao, nil
}

func (s *ComissaoService) DeleteComissao(ctx context.Context, id int64) error {
	return s.comissaoRepo.SoftDelete(ctx, id)
}

// ComissItem methods
func (s *ComissaoService) CreateComissItem(ctx context.Context, req CreateComissItemRequest) (*ComissItem, error) {
	item := &ComissItem{
		ComisID: req.ComisID,
		VlrPago: req.VlrPago,
		ObsPgto: req.ObsPgto,
	}

	if err := s.comissItemRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ComissaoService) GetComissItemByID(ctx context.Context, id int64) (*ComissItem, error) {
	item, err := s.comissItemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ComissaoService) ListComissItems(ctx context.Context, comisID int64, limit, offset int) ([]ComissItem, error) {
	itens, err := s.comissItemRepo.List(ctx, comisID, limit, offset)
	if err != nil {
		return nil, err
	}

	return itens, nil
}

func (s *ComissaoService) UpdateComissItem(ctx context.Context, id int64, req UpdateComissItemRequest) (*ComissItem, error) {
	item, err := s.comissItemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ComisID != nil {
		item.ComisID = req.ComisID
	}
	if req.VlrPago != nil {
		item.VlrPago = req.VlrPago
	}
	if req.ObsPgto != nil {
		item.ObsPgto = req.ObsPgto
	}

	if err := s.comissItemRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ComissaoService) DeleteComissItem(ctx context.Context, id int64) error {
	return s.comissItemRepo.SoftDelete(ctx, id)
}
