package forms

import (
	"mime/multipart"
	"fmt"
	"errors"
	"io/ioutil"
)

type File struct {
	error
	Bytes []byte
	Name  string
	Size  int64
}

func ValueFromForm(form multipart.Form, key string) (string, error) {
	values := form.Value[key]
	if len(values) == 0 {
		return "", fmt.Errorf(`"%v" key not found`, key)
	}
	return form.Value[key][0], nil
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
		buffer, err := readBytesFromFile(file)
		parsedFiles[i] = File{
			Bytes: buffer,
			error: err,
			Name:  fileHeader.Filename,
			Size:  fileHeader.Size,
		}
		file.Close()
	}
	return parsedFiles, nil
}

func readBytesFromFile(file multipart.File) ([]byte, error) {
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("error reading file")
	}
	return buffer, nil
}
