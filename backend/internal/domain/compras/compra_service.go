package compras

import "context"

type CompraService struct {
	compraRepo CompraRepository
}

func NewCompraService(compraRepo CompraRepository) *CompraService {
	return &CompraService{
		compraRepo: compraRepo,
	}
}

type CreateCompraRequest struct {
	PedID     *int64   `json:"ped_id"`
	FornecID  *int64   `json:"fornec_id"`
	StatusID  *int64   `json:"status_id"`
	TranspID  *int64   `json:"transp_id"`
	Login     *string  `json:"login"`
	DtCompr   *string  `json:"dt_compr"`
	DtColeta  *string  `json:"dt_coleta"`
	DtChegada *string  `json:"dt_chegada"`
	Frete     *float64 `json:"frete"`
	VlrCompr  *float64 `json:"vlr_compr"`
	NF        *string  `json:"n_nf"`
	NParcelas *int     `json:"n_parcelas"`
	DtPgto    *string  `json:"dt_pgto"`
}

type UpdateCompraRequest struct {
	PedID     *int64   `json:"ped_id"`
	FornecID  *int64   `json:"fornec_id"`
	StatusID  *int64   `json:"status_id"`
	TranspID  *int64   `json:"transp_id"`
	Login     *string  `json:"login"`
	DtCompr   *string  `json:"dt_compr"`
	DtColeta  *string  `json:"dt_coleta"`
	DtChegada *string  `json:"dt_chegada"`
	Frete     *float64 `json:"frete"`
	VlrCompr  *float64 `json:"vlr_compr"`
	NF        *string  `json:"n_nf"`
	NParcelas *int     `json:"n_parcelas"`
	DtPgto    *string  `json:"dt_pgto"`
}

func (s *CompraService) Create(ctx context.Context, req CreateCompraRequest) (*Compra, error) {
	compra := &Compra{
		PedID:     req.PedID,
		FornecID:  req.FornecID,
		StatusID:  req.StatusID,
		TranspID:  req.TranspID,
		Login:     req.Login,
		Frete:     req.Frete,
		VlrCompr:  req.VlrCompr,
		NF:        req.NF,
		NParcelas: req.NParcelas,
	}

	if err := s.compraRepo.Create(ctx, compra); err != nil {
		return nil, err
	}

	return compra, nil
}

func (s *CompraService) GetByID(ctx context.Context, id int64) (*Compra, error) {
	compra, err := s.compraRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return compra, nil
}

func (s *CompraService) List(ctx context.Context, limit, offset int, pedID, fornecID, statusID, transpID *int64, login *string) ([]Compra, error) {
	opts := ListOptions{
		Limit:    limit,
		Offset:   offset,
		PedID:    pedID,
		FornecID: fornecID,
		StatusID: statusID,
		TranspID: transpID,
		Login:    login,
	}

	compras, err := s.compraRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return compras, nil
}

func (s *CompraService) Update(ctx context.Context, id int64, req UpdateCompraRequest) (*Compra, error) {
	compra, err := s.compraRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.PedID != nil {
		compra.PedID = req.PedID
	}
	if req.FornecID != nil {
		compra.FornecID = req.FornecID
	}
	if req.StatusID != nil {
		compra.StatusID = req.StatusID
	}
	if req.TranspID != nil {
		compra.TranspID = req.TranspID
	}
	if req.Login != nil {
		compra.Login = req.Login
	}
	if req.Frete != nil {
		compra.Frete = req.Frete
	}
	if req.VlrCompr != nil {
		compra.VlrCompr = req.VlrCompr
	}
	if req.NF != nil {
		compra.NF = req.NF
	}
	if req.NParcelas != nil {
		compra.NParcelas = req.NParcelas
	}

	if err := s.compraRepo.Update(ctx, compra); err != nil {
		return nil, err
	}

	return compra, nil
}

func (s *CompraService) Delete(ctx context.Context, id int64) error {
	return s.compraRepo.SoftDelete(ctx, id)
}
