package response

type BeginAuthLoginProcess struct {
	State       string `json:"state"`
	RedirectURL string `json:"redirect_url"`
}
