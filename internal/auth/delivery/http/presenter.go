package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/auth"
)

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (req loginReq) toInput() auth.LoginInput {
	return auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (req refreshTokenReq) toInput() auth.RefreshTokenInput {
	return auth.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	}
}

type loginResp struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         userInfo `json:"user"`
}

type userInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type refreshTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h handler) newLoginResp(o auth.LoginOutput) loginResp {
	return loginResp{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
		User: userInfo{
			ID:       o.User.ID,
			Email:    o.User.Email,
			FullName: o.User.FullName,
			Role:     o.User.Role,
		},
	}
}

func (h handler) newRefreshTokenResp(o auth.RefreshTokenOutput) refreshTokenResp {
	return refreshTokenResp{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
	}
}
