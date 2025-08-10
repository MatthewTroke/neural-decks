package validation

type ValidationError struct {
	Code    string
	Message string
	Field   string
}
