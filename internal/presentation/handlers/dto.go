package handlers

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	AccessLevel int    `json:"access_level" binding:"required"`
}

type RegisterResponse struct {
	Username string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}
