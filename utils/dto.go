package utils

type BaseAuthDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,ContainSpecialChar"`
}

type RegistrationDto struct {
	BaseAuthDto
	Name string `json:"name" validate:"required,min=1,max=50"`
}

type SignInDto struct {
	BaseAuthDto
}

type CreateNoteDto struct {
	Title   string `json:"title" validate:"required,min=1,max=100"`
	Content string `json:"content,omitempty"`
}

type UpdateNoteDto struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Content *string `json:"content,omitempty" validate:"omitempty"`
}
