package km

import "context"

type KmService struct {
	kmConfigRepo          KmConfigRepository
	kmValorKmVigenciaRepo KmValorKmVigenciaRepository
	kmReembolsoLoteRepo    KmReembolsoLoteRepository
}

func NewKmService(kmConfigRepo KmConfigRepository, kmValorKmVigenciaRepo KmValorKmVigenciaRepository, kmReembolsoLoteRepo KmReembolsoLoteRepository) *KmService {
	return &KmService{
		kmConfigRepo:          kmConfigRepo,
		kmValorKmVigenciaRepo: kmValorKmVigenciaRepo,
		kmReembolsoLoteRepo:    kmReembolsoLoteRepo,
	}
}

type CreateKmConfigRequest struct {
	GpsAccuracyMaxM      *int    `json:"gps_accuracy_max_m"`
	GpsDistanciaMaxLeadM *int    `json:"gps_distancia_max_lead_m"`
	MapProvider          *string `json:"map_provider"`
	CacheEnabled         *string `json:"cache_enabled"`
}

type UpdateKmConfigRequest struct {
	GpsAccuracyMaxM      *int    `json:"gps_accuracy_max_m"`
	GpsDistanciaMaxLeadM *int    `json:"gps_distancia_max_lead_m"`
	MapProvider          *string `json:"map_provider"`
	CacheEnabled         *string `json:"cache_enabled"`
}

type CreateKmValorKmVigenciaRequest struct {
	DtInicio   *string  `json:"dt_inicio"`
	DtFim      *string  `json:"dt_fim"`
	ValorKm    *float64 `json:"valor_km"`
	Observacao *string  `json:"observacao"`
}

type UpdateKmValorKmVigenciaRequest struct {
	DtInicio   *string  `json:"dt_inicio"`
	DtFim      *string  `json:"dt_fim"`
	ValorKm    *float64 `json:"valor_km"`
	Observacao *string  `json:"observacao"`
}

type CreateKmReembolsoLoteRequest struct {
	LoginRepre     *string  `json:"login_repre"`
	DtInicio       *string  `json:"dt_inicio"`
	DtFim          *string  `json:"dt_fim"`
	KmTotal        *float64 `json:"km_total"`
	ValorKmTotal   *float64 `json:"valor_km_total"`
	ValorTotal     *float64 `json:"valor_total"`
	StatusPagamento *string  `json:"status_pagamento"`
	PagoPor        *string  `json:"pago_por"`
	PagoEm         *string  `json:"pago_em"`
	Observacao      *string  `json:"observacao"`
}

type UpdateKmReembolsoLoteRequest struct {
	LoginRepre     *string  `json:"login_repre"`
	DtInicio       *string  `json:"dt_inicio"`
	DtFim          *string  `json:"dt_fim"`
	KmTotal        *float64 `json:"km_total"`
	ValorKmTotal   *float64 `json:"valor_km_total"`
	ValorTotal     *float64 `json:"valor_total"`
	StatusPagamento *string  `json:"status_pagamento"`
	PagoPor        *string  `json:"pago_por"`
	PagoEm         *string  `json:"pago_em"`
	Observacao      *string  `json:"observacao"`
}

// KmConfig methods
func (s *KmService) CreateKmConfig(ctx context.Context, req CreateKmConfigRequest) (*KmConfig, error) {
	config := &KmConfig{
		GpsAccuracyMaxM:      req.GpsAccuracyMaxM,
		GpsDistanciaMaxLeadM: req.GpsDistanciaMaxLeadM,
		MapProvider:          req.MapProvider,
		CacheEnabled:         req.CacheEnabled,
	}

	if err := s.kmConfigRepo.Create(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *KmService) GetKmConfigByID(ctx context.Context, id int64) (*KmConfig, error) {
	config, err := s.kmConfigRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *KmService) ListKmConfigs(ctx context.Context, limit, offset int) ([]KmConfig, error) {
	configs, err := s.kmConfigRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func (s *KmService) UpdateKmConfig(ctx context.Context, id int64, req UpdateKmConfigRequest) (*KmConfig, error) {
	config, err := s.kmConfigRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.GpsAccuracyMaxM != nil {
		config.GpsAccuracyMaxM = req.GpsAccuracyMaxM
	}
	if req.GpsDistanciaMaxLeadM != nil {
		config.GpsDistanciaMaxLeadM = req.GpsDistanciaMaxLeadM
	}
	if req.MapProvider != nil {
		config.MapProvider = req.MapProvider
	}
	if req.CacheEnabled != nil {
		config.CacheEnabled = req.CacheEnabled
	}

	if err := s.kmConfigRepo.Update(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *KmService) DeleteKmConfig(ctx context.Context, id int64) error {
	return s.kmConfigRepo.SoftDelete(ctx, id)
}

// KmValorKmVigencia methods
func (s *KmService) CreateKmValorKmVigencia(ctx context.Context, req CreateKmValorKmVigenciaRequest) (*KmValorKmVigencia, error) {
	valorKm := &KmValorKmVigencia{
		ValorKm:    req.ValorKm,
		Observacao: req.Observacao,
	}

	if err := s.kmValorKmVigenciaRepo.Create(ctx, valorKm); err != nil {
		return nil, err
	}

	return valorKm, nil
}

func (s *KmService) GetKmValorKmVigenciaByID(ctx context.Context, id int64) (*KmValorKmVigencia, error) {
	valorKm, err := s.kmValorKmVigenciaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return valorKm, nil
}

func (s *KmService) ListKmValorKmVigencia(ctx context.Context, limit, offset int) ([]KmValorKmVigencia, error) {
	valores, err := s.kmValorKmVigenciaRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return valores, nil
}

func (s *KmService) UpdateKmValorKmVigencia(ctx context.Context, id int64, req UpdateKmValorKmVigenciaRequest) (*KmValorKmVigencia, error) {
	valorKm, err := s.kmValorKmVigenciaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ValorKm != nil {
		valorKm.ValorKm = req.ValorKm
	}
	if req.Observacao != nil {
		valorKm.Observacao = req.Observacao
	}

	if err := s.kmValorKmVigenciaRepo.Update(ctx, valorKm); err != nil {
		return nil, err
	}

	return valorKm, nil
}

func (s *KmService) DeleteKmValorKmVigencia(ctx context.Context, id int64) error {
	return s.kmValorKmVigenciaRepo.SoftDelete(ctx, id)
}

// KmReembolsoLote methods
func (s *KmService) CreateKmReembolsoLote(ctx context.Context, req CreateKmReembolsoLoteRequest) (*KmReembolsoLote, error) {
	lote := &KmReembolsoLote{
		LoginRepre:     req.LoginRepre,
		KmTotal:        req.KmTotal,
		ValorKmTotal:   req.ValorKmTotal,
		ValorTotal:     req.ValorTotal,
		StatusPagamento: req.StatusPagamento,
		PagoPor:        req.PagoPor,
		Observacao:      req.Observacao,
	}

	if err := s.kmReembolsoLoteRepo.Create(ctx, lote); err != nil {
		return nil, err
	}

	return lote, nil
}

func (s *KmService) GetKmReembolsoLoteByID(ctx context.Context, id int64) (*KmReembolsoLote, error) {
	lote, err := s.kmReembolsoLoteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return lote, nil
}

func (s *KmService) ListKmReembolsoLotes(ctx context.Context, limit, offset int, loginRepre, statusPagamento *string) ([]KmReembolsoLote, error) {
	opts := ListOptions{
		Limit:          limit,
		Offset:         offset,
		LoginRepre:     loginRepre,
		StatusPagamento: statusPagamento,
	}

	lotes, err := s.kmReembolsoLoteRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return lotes, nil
}

func (s *KmService) UpdateKmReembolsoLote(ctx context.Context, id int64, req UpdateKmReembolsoLoteRequest) (*KmReembolsoLote, error) {
	lote, err := s.kmReembolsoLoteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.LoginRepre != nil {
		lote.LoginRepre = req.LoginRepre
	}
	if req.KmTotal != nil {
		lote.KmTotal = req.KmTotal
	}
	if req.ValorKmTotal != nil {
		lote.ValorKmTotal = req.ValorKmTotal
	}
	if req.ValorTotal != nil {
		lote.ValorTotal = req.ValorTotal
	}
	if req.StatusPagamento != nil {
		lote.StatusPagamento = req.StatusPagamento
	}
	if req.PagoPor != nil {
		lote.PagoPor = req.PagoPor
	}
	if req.Observacao != nil {
		lote.Observacao = req.Observacao
	}

	if err := s.kmReembolsoLoteRepo.Update(ctx, lote); err != nil {
		return nil, err
	}

	return lote, nil
}

func (s *KmService) DeleteKmReembolsoLote(ctx context.Context, id int64) error {
	return s.kmReembolsoLoteRepo.SoftDelete(ctx, id)
}
