package file

import (
	"os"
	"io/ioutil"
	"github.com/pkg/errors"
	"github.com/labstack/gommon/log"
)

const ImageDir = "/usr/local/share/images"

func WriteImage(fileName string, content []byte) (string, error) {
	if err := os.MkdirAll(ImageDir, 0777); err != nil {
		return "", err
	}
	if _, err := os.Stat(ImageDir + "/" + fileName); os.IsNotExist(err) {
		if err := ioutil.WriteFile(ImageDir+"/"+fileName, content, os.ModePerm); err != nil {
			return "", err
		}
		return "/images/" + fileName, nil
	}
	file, err := os.OpenFile(ImageDir+"/"+fileName, os.O_WRONLY, os.ModeAppend)
	defer file.Close()
	if err != nil {
		return "", err
	}
	if bW, err := file.Write(content); err != nil || bW == 0 {
		log.Infof("error and bytes %v %v", err, bW)
		return "", errors.Errorf(`Cannot write file "%v"`, fileName)
	}
	return "/images/" + fileName, nil
}
