package model

type ImageProcessingResponse struct {
	Image []byte `json:"image,omitempty"`
	Error string `json:"error,omitempty"`
}

type ImagesProcessingResponse struct {
	Items []ImageProcessingResponse `json:"items"`
}

type ErrorResponse struct {
	Code int `json:"status_code"`
	Message string `json:"message"`
}
