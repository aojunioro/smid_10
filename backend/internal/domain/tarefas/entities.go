package tarefas

// Tarefa representa um lembrete ou atividade operacional.
type Tarefa struct {
	ID         int64  `json:"id"`
	Tarefa     string `json:"tarefa"`
	DtTarefa   string `json:"dt_tarefa"`
	HrTarefa   string `json:"hr_tarefa"`
	Status     string `json:"status"` // N = Pendente, S = Resolvido
	Login      string `json:"login"`
	LeadID     *int64 `json:"lead_id"`
}
