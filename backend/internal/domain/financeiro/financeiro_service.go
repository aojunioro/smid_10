package financeiro

import "context"

type FinanceiroService struct {
	contaPagarRepo    FinContaPagarRepository
	contaReceberRepo  FinContaReceberRepository
}

func NewFinanceiroService(contaPagarRepo FinContaPagarRepository, contaReceberRepo FinContaReceberRepository) *FinanceiroService {
	return &FinanceiroService{
		contaPagarRepo:   contaPagarRepo,
		contaReceberRepo: contaReceberRepo,
	}
}

type CreateContaPagarRequest struct {
	CategoriaID      *int64   `json:"categoria_id"`
	PedidoID         *int64   `json:"pedido_id"`
	Descricao        *string  `json:"descricao"`
	Valor            *float64 `json:"valor"`
	DtVencimento     *string  `json:"dt_vencimento"`
	Status           *string  `json:"status"`
	Observacao       *string  `json:"observacao"`
	LancamentoAuto   *string  `json:"lancamento_automatico"`
	RecorrenciaID    *int64   `json:"recorrencia_id"`
	DtVencOrig       *string  `json:"dt_venc_orig"`
}

type UpdateContaPagarRequest struct {
	CategoriaID      *int64   `json:"categoria_id"`
	PedidoID         *int64   `json:"pedido_id"`
	Descricao        *string  `json:"descricao"`
	Valor            *float64 `json:"valor"`
	DtVencimento     *string  `json:"dt_vencimento"`
	DtPagamento      *string  `json:"dt_pagamento"`
	Status           *string  `json:"status"`
	Observacao       *string  `json:"observacao"`
	LancamentoAuto   *string  `json:"lancamento_automatico"`
	RecorrenciaID    *int64   `json:"recorrencia_id"`
	DtVencOrig       *string  `json:"dt_venc_orig"`
}

type CreateContaReceberRequest struct {
	CategoriaID      *int64   `json:"categoria_id"`
	PedidoID         *int64   `json:"pedido_id"`
	ClienteNome      *string  `json:"cliente_nome"`
	ClienteCPF       *string  `json:"cliente_cpf"`
	Descricao        *string  `json:"descricao"`
	Valor            *float64 `json:"valor"`
	DtVencimento     *string  `json:"dt_vencimento"`
	Status           *string  `json:"status"`
	Observacao       *string  `json:"observacao"`
	LancamentoAuto   *string  `json:"lancamento_automatico"`
	RecorrenciaID    *int64   `json:"recorrencia_id"`
	DtVencOrig       *string  `json:"dt_venc_orig"`
}

type UpdateContaReceberRequest struct {
	CategoriaID      *int64   `json:"categoria_id"`
	PedidoID         *int64   `json:"pedido_id"`
	ClienteNome      *string  `json:"cliente_nome"`
	ClienteCPF       *string  `json:"cliente_cpf"`
	Descricao        *string  `json:"descricao"`
	Valor            *float64 `json:"valor"`
	DtVencimento     *string  `json:"dt_vencimento"`
	DtRecebimento    *string  `json:"dt_recebimento"`
	Status           *string  `json:"status"`
	Observacao       *string  `json:"observacao"`
	LancamentoAuto   *string  `json:"lancamento_automatico"`
	RecorrenciaID    *int64   `json:"recorrencia_id"`
	DtVencOrig       *string  `json:"dt_venc_orig"`
}

// ContaPagar methods
func (s *FinanceiroService) CreateContaPagar(ctx context.Context, req CreateContaPagarRequest) (*FinContaPagar, error) {
	conta := &FinContaPagar{
		CategoriaID:    req.CategoriaID,
		PedidoID:       req.PedidoID,
		Descricao:      req.Descricao,
		Valor:          req.Valor,
		Status:         req.Status,
		Observacao:     req.Observacao,
		LancamentoAuto: req.LancamentoAuto,
		RecorrenciaID:  req.RecorrenciaID,
	}

	if err := s.contaPagarRepo.Create(ctx, conta); err != nil {
		return nil, err
	}

	return conta, nil
}

func (s *FinanceiroService) GetContaPagarByID(ctx context.Context, id int64) (*FinContaPagar, error) {
	conta, err := s.contaPagarRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return conta, nil
}

func (s *FinanceiroService) ListContasPagar(ctx context.Context, limit, offset int, categoriaID, pedidoID *int64, status *string) ([]FinContaPagar, error) {
	opts := ListOptions{
		Limit:      limit,
		Offset:     offset,
		CategoriaID: categoriaID,
		PedidoID:   pedidoID,
		Status:     status,
	}

	contas, err := s.contaPagarRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return contas, nil
}

func (s *FinanceiroService) UpdateContaPagar(ctx context.Context, id int64, req UpdateContaPagarRequest) (*FinContaPagar, error) {
	conta, err := s.contaPagarRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CategoriaID != nil {
		conta.CategoriaID = req.CategoriaID
	}
	if req.PedidoID != nil {
		conta.PedidoID = req.PedidoID
	}
	if req.Descricao != nil {
		conta.Descricao = req.Descricao
	}
	if req.Valor != nil {
		conta.Valor = req.Valor
	}
	if req.Status != nil {
		conta.Status = req.Status
	}
	if req.Observacao != nil {
		conta.Observacao = req.Observacao
	}
	if req.LancamentoAuto != nil {
		conta.LancamentoAuto = req.LancamentoAuto
	}
	if req.RecorrenciaID != nil {
		conta.RecorrenciaID = req.RecorrenciaID
	}

	if err := s.contaPagarRepo.Update(ctx, conta); err != nil {
		return nil, err
	}

	return conta, nil
}

func (s *FinanceiroService) DeleteContaPagar(ctx context.Context, id int64) error {
	return s.contaPagarRepo.SoftDelete(ctx, id)
}

// ContaReceber methods
func (s *FinanceiroService) CreateContaReceber(ctx context.Context, req CreateContaReceberRequest) (*FinContaReceber, error) {
	conta := &FinContaReceber{
		CategoriaID:    req.CategoriaID,
		PedidoID:       req.PedidoID,
		ClienteNome:    req.ClienteNome,
		ClienteCPF:     req.ClienteCPF,
		Descricao:      req.Descricao,
		Valor:          req.Valor,
		Status:         req.Status,
		Observacao:     req.Observacao,
		LancamentoAuto: req.LancamentoAuto,
		RecorrenciaID:  req.RecorrenciaID,
	}

	if err := s.contaReceberRepo.Create(ctx, conta); err != nil {
		return nil, err
	}

	return conta, nil
}

func (s *FinanceiroService) GetContaReceberByID(ctx context.Context, id int64) (*FinContaReceber, error) {
	conta, err := s.contaReceberRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return conta, nil
}

func (s *FinanceiroService) ListContasReceber(ctx context.Context, limit, offset int, categoriaID, pedidoID *int64, status *string) ([]FinContaReceber, error) {
	opts := ListOptions{
		Limit:      limit,
		Offset:     offset,
		CategoriaID: categoriaID,
		PedidoID:   pedidoID,
		Status:     status,
	}

	contas, err := s.contaReceberRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return contas, nil
}

func (s *FinanceiroService) UpdateContaReceber(ctx context.Context, id int64, req UpdateContaReceberRequest) (*FinContaReceber, error) {
	conta, err := s.contaReceberRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CategoriaID != nil {
		conta.CategoriaID = req.CategoriaID
	}
	if req.PedidoID != nil {
		conta.PedidoID = req.PedidoID
	}
	if req.ClienteNome != nil {
		conta.ClienteNome = req.ClienteNome
	}
	if req.ClienteCPF != nil {
		conta.ClienteCPF = req.ClienteCPF
	}
	if req.Descricao != nil {
		conta.Descricao = req.Descricao
	}
	if req.Valor != nil {
		conta.Valor = req.Valor
	}
	if req.Status != nil {
		conta.Status = req.Status
	}
	if req.Observacao != nil {
		conta.Observacao = req.Observacao
	}
	if req.LancamentoAuto != nil {
		conta.LancamentoAuto = req.LancamentoAuto
	}
	if req.RecorrenciaID != nil {
		conta.RecorrenciaID = req.RecorrenciaID
	}

	if err := s.contaReceberRepo.Update(ctx, conta); err != nil {
		return nil, err
	}

	return conta, nil
}

func (s *FinanceiroService) DeleteContaReceber(ctx context.Context, id int64) error {
	return s.contaReceberRepo.SoftDelete(ctx, id)
}
