package accountinputs

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
