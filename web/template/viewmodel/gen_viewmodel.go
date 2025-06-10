package viewmodel

type GenFormData struct {
	Prompt string
	Number int
}

type GenFormComponentData struct {
	Form   GenFormData
	Errors map[string]string
	Error  string
}

type Image struct {
	ID     string
	Data   string
	Status string
}

type GalleryComponentData struct {
	Images []Image
}

type GenPageData struct {
	GalleryData GalleryComponentData
	GenFormData GenFormComponentData
}
