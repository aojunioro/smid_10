package financeiro

import "time"

// FinContaPagar representa uma conta a pagar (obrigação financeira da empresa).
type FinContaPagar struct {
	ID               int64      `json:"id"`
	CategoriaID      *int64     `json:"categoria_id"`
	PedidoID         *int64     `json:"pedido_id"`
	Descricao        *string    `json:"descricao"`
	Valor            *float64   `json:"valor"`
	DtVencimento     *time.Time `json:"dt_vencimento"`
	DtPagamento      *time.Time `json:"dt_pagamento"`
	Status           *string    `json:"status"`
	Observacao       *string    `json:"observacao"`
	LancamentoAuto   *string    `json:"lancamento_automatico"`
	RecorrenciaID    *int64     `json:"recorrencia_id"`
	DtVencOrig       *time.Time `json:"dt_venc_orig"`
	CriadoEm         time.Time  `json:"criado_em"`
	AlteradoEm       *time.Time `json:"alterado_em"`
	ExcluidoEm       *time.Time `json:"excluido_em"`
}

// FinContaReceber representa uma conta a receber (valor que a empresa tem a receber).
type FinContaReceber struct {
	ID               int64      `json:"id"`
	CategoriaID      *int64     `json:"categoria_id"`
	PedidoID         *int64     `json:"pedido_id"`
	ClienteNome      *string    `json:"cliente_nome"`
	ClienteCPF       *string    `json:"cliente_cpf"`
	Descricao        *string    `json:"descricao"`
	Valor            *float64   `json:"valor"`
	DtVencimento     *time.Time `json:"dt_vencimento"`
	DtRecebimento    *time.Time `json:"dt_recebimento"`
	Status           *string    `json:"status"`
	Observacao       *string    `json:"observacao"`
	LancamentoAuto   *string    `json:"lancamento_automatico"`
	RecorrenciaID    *int64     `json:"recorrencia_id"`
	DtVencOrig       *time.Time `json:"dt_venc_orig"`
	CriadoEm         time.Time  `json:"criado_em"`
	AlteradoEm       *time.Time `json:"alterado_em"`
	ExcluidoEm       *time.Time `json:"excluido_em"`
}
