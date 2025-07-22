package response

type BeginGoogleAuthLoginResponse struct {
	State       string `json:"state"`
	RedirectURL string `json:"redirect_url"`
}
