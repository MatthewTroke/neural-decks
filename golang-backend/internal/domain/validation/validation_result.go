package validation

type ValidationResult struct {
	IsValid  bool
	Errors   []ValidationError
	Warnings []ValidationWarning
}
