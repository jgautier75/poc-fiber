package model

type Claims struct {
	Email            string `json:"email"`
	Verified         bool   `json:"email_verified"`
	Name             string `json:"name"`
	GivenName        string `json:"given_name"`
	PreferedUserName string `json:"preferred_username"`
}

type TokenPayload struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
