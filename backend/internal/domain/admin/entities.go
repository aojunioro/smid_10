package admin

import "time"

// SystemUser representa um usuário do sistema.
type SystemUser struct {
	ID             int64     `json:"id"`
	Login          string    `json:"login"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"` // nunca serializar
	SystemUnitID   *int64    `json:"system_unit_id"`
	Active         bool      `json:"active"`
	FrontpageID    *int64    `json:"frontpage_id"`
	EquipeID       *int64    `json:"equipe_id"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SystemGroup representa um grupo/perfil operacional.
type SystemGroup struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	FrontpageID *int64    `json:"frontpage_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SystemRole representa um papel/permissão lógica.
type SystemRole struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SystemProgram representa um programa/tela autorizável.
type SystemProgram struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Controller  string    `json:"controller"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SystemUnit representa uma unidade/filial.
type SystemUnit struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
