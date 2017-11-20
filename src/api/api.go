package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"errors"
	"net/http"
	"strconv"
	"forms"
	"model"
	"images"
	"github.com/sevenNt/echo-pprof"
	"file"
)

type API struct {
	router *echo.Echo
}

func New() *API {
	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	router.POST("/rotate", handleRotation)
	router.POST("/resize", handleResize)
	router.POST("/crop", handleCrop)
	router.POST("/flip", handleFlip)
	router.POST("/flop", handleFlop)
	router.GET("/images/:name", handleImage)

	echopprof.Wrap(router)

	return &API{
		router: router,
	}
}

/**
Image rotation HTTP handler
 */
func handleRotation(context echo.Context) error {
	form, err := context.MultipartForm()
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "error encoding form",
		})
	}
	angleStr, err := forms.ValueFromForm(*form, "angle")
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	angle, err := strconv.Atoi(angleStr)
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "error parsing angle",
		})
	}
	files, err := forms.FilesFromForm(*form, "file")
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	c := make(chan model.ImageProcessingResponse, len(files))
	for _, file := range files {
		go images.ProcessRotation(file, angle, c)
	}
	result := model.ImagesProcessingResponse{
		Items: make([]model.ImageProcessingResponse, len(files)),
	}
	for i := 0; i < len(files); i++ {
		result.Items[i] = <-c
	}
	return context.JSON(http.StatusOK, result)
}

/**
Image resize HTTP handler
 */
func handleResize(context echo.Context) error {
	form, err := context.MultipartForm()
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "error encoding form",
		})
	}
	widthStr, widthErr := forms.ValueFromForm(*form, "width")
	heightStr, heightErr := forms.ValueFromForm(*form, "height")
	keepRatStr, _ := forms.ValueFromForm(*form, "keepRatio")
	var keepRatio bool
	if keepRatStr == "" {
		keepRatio = true
	} else {
		if keepRatio, err = strconv.ParseBool(keepRatStr); err != nil {
			return context.JSON(http.StatusBadRequest, model.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "error parsing keepRatio",
			})
		}
	}
	if widthErr != nil && heightErr != nil {
		return errors.New(`"width" and "height" keys not found`)
	}
	files, err := forms.FilesFromForm(*form, "file")
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	c := make(chan model.ImageProcessingResponse, len(files))
	for _, file := range files {
		go images.ProcessResize(file, widthStr, heightStr, keepRatio, c)
	}
	result := model.ImagesProcessingResponse{
		Items: make([]model.ImageProcessingResponse, len(files)),
	}
	for i := 0; i < len(files); i++ {
		result.Items[i] = <-c
	}
	return context.JSON(http.StatusOK, result)
}

/**
Image crop HTTP handler
 */
func handleCrop(context echo.Context) error {
	form, err := context.MultipartForm()
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "error encoding form",
		})
	}
	widthStr, err := forms.ValueFromForm(*form, "width")
	if err != nil {
		return err
	}
	heightStr, err := forms.ValueFromForm(*form, "height")
	if err != nil {
		return err
	}
	files, err := forms.FilesFromForm(*form, "file")
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	c := make(chan model.ImageProcessingResponse, len(files))
	for _, file := range files {
		go images.ProcessCrop(file, widthStr, heightStr, c)
	}
	result := model.ImagesProcessingResponse{
		Items: make([]model.ImageProcessingResponse, len(files)),
	}
	for i := 0; i < len(files); i++ {
		result.Items[i] = <-c
	}
	return context.JSON(http.StatusOK, result)
}

/**
Image flip HTTP handler
 */
func handleFlip(context echo.Context) error {
	form, err := context.MultipartForm()
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "error encoding form",
		})
	}
	files, err := forms.FilesFromForm(*form, "file")
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	c := make(chan model.ImageProcessingResponse, len(files))
	for _, file := range files {
		go images.ProcessFlip(file, c)
	}
	result := model.ImagesProcessingResponse{
		Items: make([]model.ImageProcessingResponse, len(files)),
	}
	for i := 0; i < len(files); i++ {
		result.Items[i] = <-c
	}
	return context.JSON(http.StatusOK, result)
}

/**
Image flop HTTP handler
 */
func handleFlop(context echo.Context) error {
	form, err := context.MultipartForm()
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "error encoding form",
		})
	}
	files, err := forms.FilesFromForm(*form, "file")
	if err != nil {
		return context.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	c := make(chan model.ImageProcessingResponse, len(files))
	for _, file := range files {
		go images.ProcessFlop(file, c)
	}
	result := model.ImagesProcessingResponse{
		Items: make([]model.ImageProcessingResponse, len(files)),
	}
	for i := 0; i < len(files); i++ {
		result.Items[i] = <-c
	}
	return context.JSON(http.StatusOK, result)
}

func handleImage(context echo.Context) error {
	return context.File(file.ImageDir + "/" + context.Param("name"))
}

func (api *API) Start(port uint) {
	api.router.Start(":" + strconv.Itoa(int(port)))
}
