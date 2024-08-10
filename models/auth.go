package models

import (
	"time"
)

type RequestRegister struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FullName       string `json:"full_name"`
	NativeLanguage string `json:"native_language"`
}

type ResponseRegister struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	FullName       string `json:"full_name"`
	NativeLanguage string `json:"native_language"`
	CreatedAt      string `json:"created_at"`
}

type Response struct {
	StatusCode  int
	Description string
	Data        interface{}
}

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseLogin struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RequestRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type ResponseRefreshToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserForLogin struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

type StoreRefreshToken struct {
	UserId       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
}

type UpdatePassword struct {
	Password string `json:"password"`
}

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

type ResetPasswordReq struct {
	VerificationCode string `json:"verification_code"`
	NewPassword      string `json:"new_password"`
}
