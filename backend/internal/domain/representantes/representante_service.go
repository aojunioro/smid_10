package representantes

import "context"

type RepreDespesaExtraService struct {
	despesaRepo RepreDespesaExtraRepository
}

func NewRepreDespesaExtraService(despesaRepo RepreDespesaExtraRepository) *RepreDespesaExtraService {
	return &RepreDespesaExtraService{
		despesaRepo: despesaRepo,
	}
}

type CreateDespesaExtraRequest struct {
	Login       string    `json:"login"`
	DataDespesa *string   `json:"data_despesa"`
	CategID     *int64    `json:"categ_id"`
	Valor       *float64  `json:"valor"`
	Descricao   *string   `json:"descricao"`
	Status      string    `json:"status"`
}

type UpdateDespesaExtraRequest struct {
	Login       *string   `json:"login"`
	DataDespesa *string   `json:"data_despesa"`
	CategID     *int64    `json:"categ_id"`
	Valor       *float64  `json:"valor"`
	Descricao   *string   `json:"descricao"`
	Status      *string   `json:"status"`
}

func (s *RepreDespesaExtraService) Create(ctx context.Context, req CreateDespesaExtraRequest) (*RepreDespesaExtra, error) {
	despesa := &RepreDespesaExtra{
		Login:     req.Login,
		CategID:   req.CategID,
		Valor:     req.Valor,
		Descricao: req.Descricao,
		Status:    req.Status,
	}

	if err := s.despesaRepo.Create(ctx, despesa); err != nil {
		return nil, err
	}

	return despesa, nil
}

func (s *RepreDespesaExtraService) GetByID(ctx context.Context, id int64) (*RepreDespesaExtra, error) {
	despesa, err := s.despesaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return despesa, nil
}

func (s *RepreDespesaExtraService) List(ctx context.Context, limit, offset int, login, status *string) ([]RepreDespesaExtra, error) {
	var loginStr *string
	if login != nil {
		loginStr = login
	}

	var statusStr *string
	if status != nil {
		statusStr = status
	}

	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
		Login:  loginStr,
		Status: statusStr,
	}

	despesas, err := s.despesaRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return despesas, nil
}

func (s *RepreDespesaExtraService) Update(ctx context.Context, id int64, req UpdateDespesaExtraRequest) (*RepreDespesaExtra, error) {
	despesa, err := s.despesaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Login != nil {
		despesa.Login = *req.Login
	}
	if req.CategID != nil {
		despesa.CategID = req.CategID
	}
	if req.Valor != nil {
		despesa.Valor = req.Valor
	}
	if req.Descricao != nil {
		despesa.Descricao = req.Descricao
	}
	if req.Status != nil {
		despesa.Status = *req.Status
	}

	if err := s.despesaRepo.Update(ctx, despesa); err != nil {
		return nil, err
	}

	return despesa, nil
}

func (s *RepreDespesaExtraService) Delete(ctx context.Context, id int64) error {
	return s.despesaRepo.SoftDelete(ctx, id)
}
