package tarefas

import "context"

type TarefaService struct {
	tarefaRepo TarefaRepository
}

func NewTarefaService(tarefaRepo TarefaRepository) *TarefaService {
	return &TarefaService{
		tarefaRepo: tarefaRepo,
	}
}

type CreateTarefaRequest struct {
	Tarefa   string  `json:"tarefa"`
	DtTarefa string  `json:"dt_tarefa"`
	HrTarefa string  `json:"hr_tarefa"`
	Login    string  `json:"login"`
	LeadID   *int64 `json:"lead_id"`
}

type UpdateTarefaRequest struct {
	Tarefa   *string `json:"tarefa"`
	DtTarefa *string `json:"dt_tarefa"`
	HrTarefa *string `json:"hr_tarefa"`
	Status   *string `json:"status"`
	Login    *string `json:"login"`
	LeadID   *int64  `json:"lead_id"`
}

func (s *TarefaService) Create(ctx context.Context, req CreateTarefaRequest) (*Tarefa, error) {
	tarefa := &Tarefa{
		Tarefa:   req.Tarefa,
		DtTarefa: req.DtTarefa,
		HrTarefa: req.HrTarefa,
		Status:   "N", // Status padrão: Pendente
		Login:    req.Login,
		LeadID:   req.LeadID,
	}

	if err := s.tarefaRepo.Create(ctx, tarefa); err != nil {
		return nil, err
	}

	return tarefa, nil
}

func (s *TarefaService) GetByID(ctx context.Context, id int64) (*Tarefa, error) {
	tarefa, err := s.tarefaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return tarefa, nil
}

func (s *TarefaService) List(ctx context.Context, limit, offset int) ([]Tarefa, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	tarefas, err := s.tarefaRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return tarefas, nil
}

func (s *TarefaService) Update(ctx context.Context, id int64, req UpdateTarefaRequest) (*Tarefa, error) {
	tarefa, err := s.tarefaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Tarefa != nil {
		tarefa.Tarefa = *req.Tarefa
	}
	if req.DtTarefa != nil {
		tarefa.DtTarefa = *req.DtTarefa
	}
	if req.HrTarefa != nil {
		tarefa.HrTarefa = *req.HrTarefa
	}
	if req.Status != nil {
		tarefa.Status = *req.Status
	}
	if req.Login != nil {
		tarefa.Login = *req.Login
	}
	if req.LeadID != nil {
		tarefa.LeadID = req.LeadID
	}

	if err := s.tarefaRepo.Update(ctx, tarefa); err != nil {
		return nil, err
	}

	return tarefa, nil
}
