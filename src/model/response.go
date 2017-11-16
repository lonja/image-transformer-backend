package model

type Image struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name"`
}

type ImageProcessingResponse struct {
	Image        `json:"image,omitempty"`
	Error string `json:"error,omitempty"`
}

type ImagesProcessingResponse struct {
	Items []ImageProcessingResponse `json:"items"`
}

type ErrorResponse struct {
	Code    int    `json:"status_code"`
	Message string `json:"message"`
}
