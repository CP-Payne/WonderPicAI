package viewmodel

var (
	SignUpTitle = "Sign Up"
	LoginTitle  = "Login"
)

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

type LoginFormData struct {
	Email string
}

type LoginFormComponentData struct {
	Form   LoginFormData
	Errors map[string]string
	Error  string
	// CSRFToken string
}
