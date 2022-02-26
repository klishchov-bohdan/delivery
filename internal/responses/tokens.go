package responses

type TokensResponse struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}
