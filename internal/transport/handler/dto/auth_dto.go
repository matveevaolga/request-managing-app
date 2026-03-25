package dto

type LoginRequest struct {
	Login    string `json:"login" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
