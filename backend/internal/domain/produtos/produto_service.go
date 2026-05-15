package produtos

import "context"

type ProdutoService struct {
	produtoRepo ProdutoRepository
}

func NewProdutoService(produtoRepo ProdutoRepository) *ProdutoService {
	return &ProdutoService{
		produtoRepo: produtoRepo,
	}
}

type CreateProdutoRequest struct {
	NomeProd      string   `json:"nome_prod"`
	CategID       *int64   `json:"categ_id"`
	FornecID      *int64   `json:"fornec_id"`
	MedID         *int64   `json:"med_id"`
	ModeloID      *int64   `json:"modelo_id"`
	VlrProdCompra *float64 `json:"vlr_prod_compra"`
	VlrProdVenda  *float64 `json:"vlr_prod_venda"`
	EstoqMin      *int     `json:"estoq_min"`
	EstoqMax      *int     `json:"estoq_max"`
	Ativo         string   `json:"ativo"`
	Televendas    string   `json:"televendas"`
}

type UpdateProdutoRequest struct {
	NomeProd      *string  `json:"nome_prod"`
	CategID       *int64   `json:"categ_id"`
	FornecID      *int64   `json:"fornec_id"`
	MedID         *int64   `json:"med_id"`
	ModeloID      *int64   `json:"modelo_id"`
	VlrProdCompra *float64 `json:"vlr_prod_compra"`
	VlrProdVenda  *float64 `json:"vlr_prod_venda"`
	EstoqMin      *int     `json:"estoq_min"`
	EstoqMax      *int     `json:"estoq_max"`
	Ativo         *string  `json:"ativo"`
	Televendas    *string  `json:"televendas"`
}

func (s *ProdutoService) Create(ctx context.Context, req CreateProdutoRequest) (*Produto, error) {
	produto := &Produto{
		NomeProd:      req.NomeProd,
		CategID:       req.CategID,
		FornecID:      req.FornecID,
		MedID:         req.MedID,
		ModeloID:      req.ModeloID,
		VlrProdCompra: req.VlrProdCompra,
		VlrProdVenda:  req.VlrProdVenda,
		EstoqMin:      req.EstoqMin,
		EstoqMax:      req.EstoqMax,
		Ativo:         req.Ativo,
		Televendas:    req.Televendas,
	}

	if err := s.produtoRepo.Create(ctx, produto); err != nil {
		return nil, err
	}

	return produto, nil
}

func (s *ProdutoService) GetByID(ctx context.Context, id int64) (*Produto, error) {
	produto, err := s.produtoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return produto, nil
}

func (s *ProdutoService) List(ctx context.Context, limit, offset int, categID, medID *int64, ativo, televendas *string) ([]Produto, error) {
	opts := ListOptions{
		Limit:      limit,
		Offset:     offset,
		CategID:    categID,
		MedID:      medID,
		Ativo:      ativo,
		Televendas: televendas,
	}

	produtos, err := s.produtoRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return produtos, nil
}

func (s *ProdutoService) Update(ctx context.Context, id int64, req UpdateProdutoRequest) (*Produto, error) {
	produto, err := s.produtoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.NomeProd != nil {
		produto.NomeProd = *req.NomeProd
	}
	if req.CategID != nil {
		produto.CategID = req.CategID
	}
	if req.FornecID != nil {
		produto.FornecID = req.FornecID
	}
	if req.MedID != nil {
		produto.MedID = req.MedID
	}
	if req.ModeloID != nil {
		produto.ModeloID = req.ModeloID
	}
	if req.VlrProdCompra != nil {
		produto.VlrProdCompra = req.VlrProdCompra
	}
	if req.VlrProdVenda != nil {
		produto.VlrProdVenda = req.VlrProdVenda
	}
	if req.EstoqMin != nil {
		produto.EstoqMin = req.EstoqMin
	}
	if req.EstoqMax != nil {
		produto.EstoqMax = req.EstoqMax
	}
	if req.Ativo != nil {
		produto.Ativo = *req.Ativo
	}
	if req.Televendas != nil {
		produto.Televendas = *req.Televendas
	}

	if err := s.produtoRepo.Update(ctx, produto); err != nil {
		return nil, err
	}

	return produto, nil
}

func (s *ProdutoService) Delete(ctx context.Context, id int64) error {
	return s.produtoRepo.SoftDelete(ctx, id)
}
