package accountoutputs

type AccountOutput struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	IsEmailValid bool   `json:"is_email_valid"`
	Role         string `json:"role"`
	IsActive     bool   `json:"is_active"`
}

type AccountPublicOutput struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
