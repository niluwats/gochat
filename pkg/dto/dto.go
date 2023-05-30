package dto

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Total   int    `json:"total,omitempty"`
}
