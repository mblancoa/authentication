package controllers

type LoginRequest struct {
	UserId   string `json:"user_id" mapper:"Id"`
	Password string `json:"password" mapper:"Password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}
