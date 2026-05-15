package produtos

import "time"

// Produto representa um produto base vendido ou orçado.
type Produto struct {
	ID            int64      `json:"id"`
	NomeProd      string     `json:"nome_prod"`
	CategID       *int64     `json:"categ_id"`
	FornecID      *int64     `json:"fornec_id"`
	MedID         *int64     `json:"med_id"`
	ModeloID      *int64     `json:"modelo_id"`
	VlrProdCompra *float64   `json:"vlr_prod_compra"`
	VlrProdVenda  *float64   `json:"vlr_prod_venda"`
	EstoqMin      *int       `json:"estoq_min"`
	EstoqMax      *int       `json:"estoq_max"`
	Ativo         string     `json:"ativo"`
	Televendas    string     `json:"televendas"`
	CriadoEm      time.Time  `json:"criado_em"`
	AlteradoEm    *time.Time `json:"alterado_em"`
	ExcluidoEm    *time.Time `json:"excluido_em"`
}

// ProdCateg representa uma categoria de produto.
type ProdCateg struct {
	ID         int64      `json:"id"`
	Categoria  string     `json:"categoria"`
	Cor        *string    `json:"cor"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// ProdMedidas representa uma unidade/medida de produto.
type ProdMedidas struct {
	ID         int64      `json:"id"`
	Medida     string     `json:"medida"`
	Sigla      *string    `json:"sigla"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}

// ProdModelos representa um modelo/variante de produto.
type ProdModelos struct {
	ID         int64      `json:"id"`
	Modelo     string     `json:"modelo"`
	CriadoEm   time.Time  `json:"criado_em"`
	AlteradoEm *time.Time `json:"alterado_em"`
	ExcluidoEm *time.Time `json:"excluido_em"`
}
