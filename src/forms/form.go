package forms

import (
	"mime/multipart"
	"fmt"
	"errors"
	"bufio"
)

type File struct {
	multipart.File
	error
	Name string
	Size int64
}

func ValueFromForm(form multipart.Form, key string) (string, error) {
	values := form.Value[key]
	if len(values) == 0 {
		return "", fmt.Errorf(`"%v" key not found`, key)
	}
	return form.Value[key][0], nil
}

func FileFromForm(form multipart.Form, key string) (multipart.File, int64, error) {
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

func FilesFromForm(form multipart.Form, key string) ([]File, error) {
	files := form.File[key]
	if len(files) == 0 {
		return nil, fmt.Errorf(`"%v" key not found`, key)
	}
	parsedFiles := make([]File, len(files))
	for i, file := range files {
		fileHeader := file
		file, err := fileHeader.Open()
		parsedFiles[i] = File{
			File:  file,
			error: err,
			Name:  fileHeader.Filename,
			Size:  fileHeader.Size,
		}
	}
	return parsedFiles, nil
}

func ReadBytesFromFile(file File) ([]byte, error) {
	reader := bufio.NewReader(file)
	var buffer = make([]byte, file.Size)
	if bytesRead, err := reader.Read(buffer); err != nil || bytesRead == 0 {
		return nil, errors.New("error reading file")
	}
	return buffer, nil
}
