package log

// SystemAccessLog representa um registro de acesso/login.
type SystemAccessLog struct {
	ID              int64     `json:"id"`
	SessionID       string    `json:"sessionid"`
	Login           string    `json:"login"`
	LoginTime       string    `json:"login_time"`
	LoginYear       string    `json:"login_year"`
	LoginMonth      string    `json:"login_month"`
	LoginDay        string    `json:"login_day"`
	LogoutTime      *string   `json:"logout_time"`
	Impersonated    *bool     `json:"impersonated"`
	AccessIP        string    `json:"access_ip"`
	ImpersonatedBy  *string   `json:"impersonated_by"`
}

// SystemChangeLog representa um registro de auditoria de alteração de registros.
type SystemChangeLog struct {
	ID            int64     `json:"id"`
	LogDate       string    `json:"logdate"`
	Login         string    `json:"login"`
	TableName     string    `json:"tablename"`
	PrimaryKey    string    `json:"primarykey"`
	PKValue       string    `json:"pkvalue"`
	Operation     string    `json:"operation"`
	ColumnName    string    `json:"columnname"`
	OldValue      *string   `json:"oldvalue"`
	NewValue      *string   `json:"newvalue"`
	AccessIP      string    `json:"access_ip"`
	TransactionID string    `json:"transaction_id"`
	LogTrace      *string   `json:"log_trace"`
	SessionID     string    `json:"session_id"`
	ClassName     string    `json:"class_name"`
	PhpSapi       string    `json:"php_sapi"`
	LogYear       string    `json:"log_year"`
	LogMonth      string    `json:"log_month"`
	LogDay        string    `json:"log_day"`
}

// SystemSqlLog representa um registro de SQL executado.
type SystemSqlLog struct {
	ID            int64     `json:"id"`
	LogDate       string    `json:"logdate"`
	Login         string    `json:"login"`
	DatabaseName  string    `json:"database_name"`
	SqlCommand    string    `json:"sql_command"`
	StatementType string    `json:"statement_type"`
	AccessIP      string    `json:"access_ip"`
	TransactionID string    `json:"transaction_id"`
	LogTrace      *string   `json:"log_trace"`
	SessionID     string    `json:"session_id"`
	ClassName     string    `json:"class_name"`
	PhpSapi       string    `json:"php_sapi"`
	RequestID     string    `json:"request_id"`
	LogYear       string    `json:"log_year"`
	LogMonth      string    `json:"log_month"`
	LogDay        string    `json:"log_day"`
}

// SystemRequestLog representa um registro de requisição HTTP.
type SystemRequestLog struct {
	ID              int64     `json:"id"`
	Endpoint        string    `json:"endpoint"`
	LogDate         string    `json:"logdate"`
	LogYear         string    `json:"log_year"`
	LogMonth        string    `json:"log_month"`
	LogDay          string    `json:"log_day"`
	SessionID       string    `json:"session_id"`
	Login           string    `json:"login"`
	AccessIP        string    `json:"access_ip"`
	ClassName       string    `json:"class_name"`
	ClassMethod     string    `json:"class_method"`
	HttpHost        string    `json:"http_host"`
	ServerPort      string    `json:"server_port"`
	RequestURI      string    `json:"request_uri"`
	RequestMethod   string    `json:"request_method"`
	QueryString     *string   `json:"query_string"`
	RequestHeaders  *string   `json:"request_headers"`
	RequestBody     *string   `json:"request_body"`
	RequestDuration int       `json:"request_duration"`
}
