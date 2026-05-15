package km

import "time"

// KmConfig representa as configurações globais de KM/GPS.
type KmConfig struct {
	ID                   int64      `json:"id"`
	GpsAccuracyMaxM      *int       `json:"gps_accuracy_max_m"`
	GpsDistanciaMaxLeadM *int       `json:"gps_distancia_max_lead_m"`
	MapProvider          *string    `json:"map_provider"`
	CacheEnabled         *string    `json:"cache_enabled"`
	CriadoEm             time.Time  `json:"criado_em"`
	AlteradoEm           *time.Time `json:"alterado_em"`
	ExcluidoEm           *time.Time `json:"excluido_em"`
}

// KmValorKmVigencia representa o valor do KM por período de vigência.
type KmValorKmVigencia struct {
	ID         int64      `json:"id"`
	DtInicio   *time.Time `json:"dt_inicio"`
	DtFim      *time.Time `json:"dt_fim"`
	ValorKm    *float64   `json:"valor_km"`
	Observacao *string    `json:"observacao"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// KmReembolsoLote representa um lote de reembolso por representante/período.
type KmReembolsoLote struct {
	ID               int64      `json:"id"`
	LoginRepre       *string    `json:"login_repre"`
	DtInicio         *time.Time `json:"dt_inicio"`
	DtFim            *time.Time `json:"dt_fim"`
	KmTotal          *float64   `json:"km_total"`
	ValorKmTotal     *float64   `json:"valor_km_total"`
	ValorTotal       *float64   `json:"valor_total"`
	StatusPagamento   *string    `json:"status_pagamento"`
	PagoPor          *string    `json:"pago_por"`
	PagoEm           *time.Time `json:"pago_em"`
	Observacao       *string    `json:"observacao"`
	CriadoEm         time.Time  `json:"criado_em"`
	AlteradoEm       *time.Time `json:"alterado_em"`
	ExcluidoEm       *time.Time `json:"excluido_em"`
}
