package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"errors"
	"gopkg.in/h2non/bimg.v1"
	"net/http"
	"strconv"
	"bufio"
	"mime/multipart"
	"fmt"
)

type API struct {
	router *echo.Echo
}

func New() *API {
	router := echo.New()

	router.Use(middleware.Logger())

	router.POST("/rotate", handleRotation)

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
		return errors.New("error encoding form")
	}
	angleStr, err := valueFromForm(form, "angle")
	if err != nil {
		return err
	}
	angle, err := strconv.Atoi(angleStr)
	if err != nil {
		return errors.New("error parsing angle")
	}
	file, size, err := fileFromForm(form, "file")
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	var buffer = make([]byte, size)
	if bytesRead, err := reader.Read(buffer); err != nil || bytesRead == 0 {
		return errors.New("error reading file")
	}
	image := bimg.NewImage(buffer)
	buffer, err = image.Rotate(bimg.Angle(angle))
	if err != nil {
		return errors.New("error rotating image")
	}
	image = bimg.NewImage(buffer)
	return context.Blob(http.StatusOK, "image/*", image.Image())
}

func valueFromForm(form *multipart.Form, key string) (string, error) {
	values := form.Value[key]
	if len(values) == 0 {
		return "", fmt.Errorf(`"%v" key not found`, key)
	}
	return form.Value[key][0], nil
}

func fileFromForm(form *multipart.Form, key string) (multipart.File, int64, error) {
	files := form.File[key]
	if len(files) == 0 {
		return nil, 0, fmt.Errorf(`"%v" key not found`, key)
	}
	fileHeader := form.File[key][0]
	file, err := fileHeader.Open()
	if err != nil {
		return nil, 0, errors.New("error opening file")
	}
	return file, fileHeader.Size, nil
}

func (api *API) Start(port uint) {
	api.router.Start(":" + strconv.Itoa(int(port)))
}
