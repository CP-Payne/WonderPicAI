package viewmodel

type ErrorPageData struct {
	StatusCode  int
	Title       string
	Message     string
	ShowDetails bool
	ErrorID     string
}
