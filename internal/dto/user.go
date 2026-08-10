package dto

type RegisterRequest struct {
	Username string `json:"username" valid:"stringlength(8|64),matches(^[A-Za-z0-9]+$)" validate:"min=8,max=64"`
	Password string `json:"password" valid:"stringlength(8|64)" validate:"min=8,max=64"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"min=8,max=64"`
	Password string `json:"password" validate:"min=8,max=64"`
}

type AccessTokenResponse struct {
	Token string `json:"access_token"`
}

type UserDataResponse struct {
	Username string `json:"username"`
}
