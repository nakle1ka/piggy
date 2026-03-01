package handler

// general

type ErrorResponse struct {
	Error string `json:"error"`
}

// user & auth

type CreateUserDTO struct {
	Username string `json:"username" binding:"required,max=35"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type TokenDTO struct {
	AccessToken string `json:"access_token"`
}

// piggy

type CreatePiggyDTO struct {
	Title  string `json:"title" binding:"required,max=35"`
	Amount int64  `json:"amount" binding:"required,min=1"`
}

type UpdatePiggyDTO struct {
	Title string `json:"title" binding:"required,max=35"`
}

type AmountPiggyDTO struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}
