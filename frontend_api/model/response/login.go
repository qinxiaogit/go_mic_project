package response

type LoginResponse struct {
	Token string `json:"token" binding:"required"`
}
