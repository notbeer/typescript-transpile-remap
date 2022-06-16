package internal

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/notbeer/typescript-transpile-remap/logger"
)

func JSONUnmarshal(path string) map[string]interface{} {
	file, err := ioutil.ReadFile(path)
	fileName := filepath.Base(path)
	logger.Error("Error occured while reading "+fileName+" - ", err)

	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	logger.Error("Error occured during json.Unmarshal() with "+fileName+" - ", err)

	return data
}
