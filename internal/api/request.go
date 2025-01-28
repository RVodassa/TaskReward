package api

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
