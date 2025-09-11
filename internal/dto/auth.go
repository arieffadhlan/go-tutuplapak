package dto

type (
	AuthEmailRequest struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	AuthPhoneRequest struct {
		Phone    string `json:"phone"    validate:"required,e164"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	AuthEmailResponse struct {
		Token string `json:"token"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	AuthPhoneResponse struct {
		Token string `json:"token"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
)
