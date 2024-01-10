package accountinputs

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateInput struct {
	Email           string `json:"email" validate:"omitempty,email"`
	Name            string `json:"name" validate:"omitempty"`
	Avatar          string `json:"avatar" validate:"omitempty,url"`
	OldPassword     string `json:"old_password" validate:"omitempty,min=8"`
	Password        string `json:"password" validate:"omitempty,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"omitempty,min=8"`
}

type LoginInput struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required,min=8"`
	Password        string `json:"password" validate:"required,min=8"`
}
