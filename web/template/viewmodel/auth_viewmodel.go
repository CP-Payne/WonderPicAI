package viewmodel

type SignupFormData struct {
	Username string
	Email    string
}

type SignupFormComponentData struct {
	Form   SignupFormData
	Errors map[string]string
	Error  string
	// CSRFToken string
}
