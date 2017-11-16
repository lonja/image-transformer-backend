package images

import (
	"forms"
	"model"
	"gopkg.in/h2non/bimg.v1"
	"strings"
	"strconv"
	"errors"
	"fmt"
	fs "file"
)

const backendURL = "http://localhost:8080"

func ProcessRotation(file forms.File, angle int, c chan model.ImageProcessingResponse) {
	buffer, err := bimg.Resize(file.Bytes, bimg.Options{
		Rotate:  bimg.Angle(angle),
		Quality: 100,
	})
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error rotating image "%v"`, file.Name),
		}
		return
	}
	url, err := fs.WriteImage(file.Name, buffer)
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error writing image "%v"`, file.Name),
		}
		return
	}
	c <- model.ImageProcessingResponse{
		Image: model.Image{
			Name: file.Name,
			URL: backendURL + url,
		},
	}
}

func ProcessResize(file forms.File, width, height string, c chan model.ImageProcessingResponse) {
	image := bimg.NewImage(file.Bytes)
	var newSize bimg.ImageSize
	curSize, err := image.Size()
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: err.Error(),
		}
		return
	}
	newSize, err = parseSize(width, height, curSize)
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: err.Error(),
		}
		return
	}
	bytes, err := image.Process(bimg.Options{
		Width:   newSize.Width,
		Height:  newSize.Height,
		Force:   true,
		Quality: 100,
	})
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error resizing image "%v"`, file.Name),
		}
		return
	}
	url, err := fs.WriteImage(file.Name, bytes)
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error writing image "%v"`, file.Name),
		}
		return
	}
	c <- model.ImageProcessingResponse{
		Image: model.Image{
			Name: file.Name,
			URL: backendURL + url,
		},
	}
}

func ProcessCrop(file forms.File, width, height string, c chan model.ImageProcessingResponse) {
	image := bimg.NewImage(file.Bytes)
	var newSize bimg.ImageSize
	curSize, err := image.Size()
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: err.Error(),
		}
	}
	newSize, err = parseSize(width, height, curSize)
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: err.Error(),
		}
	}
	bytes, err := image.Process(bimg.Options{
		AreaWidth:  newSize.Width,
		AreaHeight: newSize.Height,
		Quality:    100,
	})
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error cropping image "%v"`, file.Name),
		}
	}
	url, err := fs.WriteImage(file.Name, bytes)
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error writing image "%v"`, file.Name),
		}
		return
	}
	c <- model.ImageProcessingResponse{
		Image: model.Image{
			Name: file.Name,
			URL: backendURL + url,
		},
	}
}

func ProcessFlip(file forms.File, c chan model.ImageProcessingResponse) {
	image := bimg.NewImage(file.Bytes)
	bytes, err := image.Process(bimg.Options{
		Flip:    true,
		Quality: 100,
	})
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error flipping image "%v"`, file.Name),
		}
		return
	}
	url, err := fs.WriteImage(file.Name, bytes)
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error writing image "%v"`, file.Name),
		}
		return
	}
	c <- model.ImageProcessingResponse{
		Image: model.Image{
			Name: file.Name,
			URL: backendURL + url,
		},
	}
}

func ProcessFlop(file forms.File, c chan model.ImageProcessingResponse) {
	bytes, err := bimg.Resize(file.Bytes, bimg.Options{
		Flop:    true,
		Quality: 100,
	})
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error flopping image "%v"`, file.Name),
		}
		return
	}
	url, err := fs.WriteImage(file.Name, bytes)
	if err != nil {
		c <- model.ImageProcessingResponse{
			Error: fmt.Sprintf(`error writing image "%v"`, file.Name),
		}
		return
	}
	c <- model.ImageProcessingResponse{
		Image: model.Image{
			Name: file.Name,
			URL: backendURL + url,
		},
	}
}

//TODO: Refactor this shit. This function is too long
func parseSize(width, height string, oldSize bimg.ImageSize) (bimg.ImageSize, error) {
	if width != "" && height == "" {
		if strings.Contains(width, "px") {
			widthStr := strings.Split(width, "px")[0]
			dim, err := strconv.Atoi(widthStr)
			ratio := float32(dim) / float32(oldSize.Width)
			return bimg.ImageSize{
				Width:  dim,
				Height: int(ratio * float32(oldSize.Height)),
			}, err
		} else if strings.Contains(width, "%") {
			heightStr := strings.Split(width, "%")[0]
			dim, err := strconv.Atoi(heightStr)
			ratio := float32(dim) / 100
			return bimg.ImageSize{
				Width:  int(ratio * float32(oldSize.Width)),
				Height: int(ratio * float32(oldSize.Height)),
			}, err
		}
		return bimg.ImageSize{}, errors.New("unknown unit")
	} else if width == "" && height != "" {
		if strings.Contains(height, "px") {
			heightStr := strings.Split(height, "px")[0]
			dim, err := strconv.Atoi(heightStr)
			ratio := float32(dim) / float32(oldSize.Height)
			return bimg.ImageSize{
				Width:  int(ratio * float32(oldSize.Width)),
				Height: dim,
			}, err
		} else if strings.Contains(width, "%") {
			heightStr := strings.Split(width, "%")[0]
			dim, err := strconv.Atoi(heightStr)
			ratio := float32(dim) / 100
			return bimg.ImageSize{
				Width:  int(ratio * float32(oldSize.Width)),
				Height: int(ratio * float32(oldSize.Height)),
			}, err
		}
		return bimg.ImageSize{}, errors.New("unknown unit")
	} else if width != "" && height != "" {
		if strings.Contains(width, "px") {
			widthStr := strings.Split(width, "px")[0]
			widthDim, err := strconv.Atoi(widthStr)
			heightStr := strings.Split(height, "px")[0]
			heightDim, err := strconv.Atoi(heightStr)
			widthRatio := float32(widthDim) / float32(oldSize.Width)
			heightRatio := float32(heightDim) / float32(oldSize.Height)
			return bimg.ImageSize{
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
			return bimg.ImageSize{
				Width:  int(widthRatio * float32(oldSize.Width)),
				Height: int(heightRatio * float32(oldSize.Height)),
			}, err
		}
		return bimg.ImageSize{}, errors.New("unknown unit")
	}
	return bimg.ImageSize{}, errors.New("empty dimensions")
}
