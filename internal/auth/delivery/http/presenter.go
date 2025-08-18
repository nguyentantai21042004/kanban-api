package http

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/auth"
)

type respObj struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Alias string `json:"alias,omitempty"`
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (req loginReq) toInput() auth.LoginInput {
	return auth.LoginInput{
		Username: req.Username,
		Password: req.Password,
	}
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (req refreshTokenReq) toInput() auth.RefreshTokenInput {
	return auth.RefreshTokenInput{
		RfrToken: req.RefreshToken,
	}
}

type loginResp struct {
	AccessToken string   `json:"access_token"`
	User        userInfo `json:"user"`
}

type userInfo struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	FullName string  `json:"full_name"`
	Role     respObj `json:"role"`
}

type refreshTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h handler) newLoginResp(o auth.LoginOutput) loginResp {
	return loginResp{
		AccessToken: o.AssToken,
		User: userInfo{
			ID:       o.User.ID,
			Username: o.User.Username,
			FullName: o.User.FullName,
			Role: respObj{
				ID:    o.Role.ID,
				Name:  o.Role.Name,
				Alias: o.Role.Alias,
			},
		},
	}
}

func (h handler) newRefreshTokenResp(o auth.RefreshTokenOutput) refreshTokenResp {
	return refreshTokenResp{
		AccessToken:  o.AssToken,
		RefreshToken: o.RfrToken,
	}
}
