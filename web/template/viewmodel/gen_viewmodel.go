package viewmodel

type GenFormData struct {
	Prompt          string
	Credits         int
	MinCost         int
	MaxImagesPerGen int
	ImageCount      int
	HasFailedImages bool
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
