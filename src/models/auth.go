package models

type SignInResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpResult struct {
	IsCreated bool `json:"is_created"`
}
