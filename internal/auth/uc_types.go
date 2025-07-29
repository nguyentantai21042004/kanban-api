package auth

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginOutput struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}
