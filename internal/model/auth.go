package model

type CreateAuthTokenRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}
type CreateAuthTokenResponse struct {
	Token string `json:"token"`
}

type NewAPIUserRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	Website   string `json:"website,omitempty"`
}
