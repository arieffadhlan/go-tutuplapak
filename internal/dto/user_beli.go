package dto

type (
	AuthUserBeliRequest struct {
		Username string `json:"username" validate:"required,min=5,max=30"`
		Password string `json:"password" validate:"required,min=5,max=30"`
		Email    string `json:"email" validate:"required,email"`
	}
	AuthLoginUserBeliRequest struct {
		Username string `json:"username" validate:"required,min=5,max=30"`
		Password string `json:"password" validate:"required,min=5,max=30"`
	}
	AuthLoginBeliResponse struct {
		Token string `json:"token"`
	}
)
