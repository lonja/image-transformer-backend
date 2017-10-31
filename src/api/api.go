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
	"strings"
)

type API struct {
	router *echo.Echo
}

func New() *API {
	router := echo.New()

	router.Use(middleware.Logger())

	router.POST("/rotate", handleRotation)
	router.POST("/resize", handleResize)

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

/**
Image resize HTTP handler
 */
func handleResize(context echo.Context) error {
	form, err := context.MultipartForm()
	if err != nil {
		return errors.New("error encoding form")
	}
	widthStr, widthErr := valueFromForm(form, "width")
	heightStr, heightErr := valueFromForm(form, "height")
	if widthErr != nil && heightErr != nil {
		return errors.New(`"width" and "height" keys not found`)
	}
	file, size, err := fileFromForm(form, "file")
	if err != nil {
		return errors.New("error opening file")
	}
	reader := bufio.NewReader(file)
	var buffer = make([]byte, size)
	if bytesRead, err := reader.Read(buffer); err != nil || bytesRead == 0 {
		return errors.New("error reading file")
	}
	image := bimg.NewImage(buffer)
	var newSize *bimg.ImageSize
	curSize, err := image.Size()
	if err != nil {
		return err
	}
	newSize, err = parseSize(widthStr, heightStr, curSize)
	if err != nil {
		return err
	}
	buffer, err = image.ForceResize(newSize.Width, newSize.Height)
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

//TODO: Refactor this shit. This function is too long
func parseSize(width, height string, oldSize bimg.ImageSize) (*bimg.ImageSize, error) {
	if width != "" && height == "" {
		if strings.Contains(width, "px") {
			widthStr := strings.Split(width, "px")[0]
			dim, err := strconv.Atoi(widthStr)
			ratio := float32(dim) / float32(oldSize.Width)
			return &bimg.ImageSize{
				Width:  dim,
				Height: int(ratio * float32(oldSize.Height)),
			}, err
		} else if strings.Contains(width, "%") {
			heightStr := strings.Split(width, "%")[0]
			dim, err := strconv.Atoi(heightStr)
			ratio := float32(dim) / 100
			return &bimg.ImageSize{
				Width:  int(ratio * float32(oldSize.Width)),
				Height: int(ratio * float32(oldSize.Height)),
			}, err
		}
		return nil, errors.New("unknown unit")
	} else if width == "" && height != "" {
		if strings.Contains(height, "px") {
			heightStr := strings.Split(height, "px")[0]
			dim, err := strconv.Atoi(heightStr)
			ratio := float32(dim) / float32(oldSize.Height)
			return &bimg.ImageSize{
				Width:  int(ratio * float32(oldSize.Width)),
				Height: dim,
			}, err
		} else if strings.Contains(width, "%") {
			heightStr := strings.Split(width, "%")[0]
			dim, err := strconv.Atoi(heightStr)
			ratio := float32(dim) / 100
			return &bimg.ImageSize{
				Width:  int(ratio * float32(oldSize.Width)),
				Height: int(ratio * float32(oldSize.Height)),
			}, err
		}
		return nil, errors.New("unknown unit")
	} else if width != "" && height != "" {
		if strings.Contains(width, "px") {
			widthStr := strings.Split(width, "px")[0]
			widthDim, err := strconv.Atoi(widthStr)
			heightStr := strings.Split(height, "px")[0]
			heightDim, err := strconv.Atoi(heightStr)
			widthRatio := float32(widthDim) / float32(oldSize.Width)
			heightRatio := float32(heightDim) / float32(oldSize.Height)
			return &bimg.ImageSize{
				Width:  int(widthRatio * float32(oldSize.Width)),
				Height: int(heightRatio * float32(oldSize.Height)),
			}, err
		} else if strings.Contains(width, "%") {
			widthStr := strings.Split(width, "%")[0]
			widthDim, err := strconv.Atoi(widthStr)
			heightStr := strings.Split(width, "%")[0]
			heightDim, err := strconv.Atoi(heightStr)
			widthRatio := float32(widthDim) / 100
			heightRatio := float32(heightDim) / 100
			return &bimg.ImageSize{
				Width:  int(widthRatio * float32(oldSize.Width)),
				Height: int(heightRatio * float32(oldSize.Height)),
			}, err
		}
		return nil, errors.New("unknown unit")
	}
	return nil, errors.New("empty dimensions")
}

func (api *API) Start(port uint) {
	api.router.Start(":" + strconv.Itoa(int(port)))
}
