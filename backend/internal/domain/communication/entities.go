package communication

// SystemNotification representa uma notificação do sistema para um usuário.
type SystemNotification struct {
	ID              int64  `json:"id"`
	SystemUserID    *int64 `json:"system_user_id"`
	SystemUserToID  *int64 `json:"system_user_to_id"`
	Subject         string `json:"subject"`
	Message         string `json:"message"`
	DtMessage       string `json:"dt_message"`
	ActionURL       *string `json:"action_url"`
	ActionLabel     *string `json:"action_label"`
	Icon            *string `json:"icon"`
	Checked         *bool  `json:"checked"`
}

// SystemMessage representa uma mensagem interna entre usuários.
type SystemMessage struct {
	ID              int64  `json:"id"`
	SystemUserID    *int64 `json:"system_user_id"`
	SystemUserToID  *int64 `json:"system_user_to_id"`
	Subject         string `json:"subject"`
	Message         string `json:"message"`
	DtMessage       string `json:"dt_message"`
	Checked         *bool  `json:"checked"`
	Removed         *bool  `json:"removed"`
	Viewed          *bool  `json:"viewed"`
	Attachments     *string `json:"attachments"`
}
