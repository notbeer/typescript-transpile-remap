package tools

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/notbeer/typescript-transpile-remap/internal"
	"github.com/notbeer/typescript-transpile-remap/logger"
)

func walkDirectory(dir string) ([]string, error) {
	paths := make([]string, 0)

	err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".js") {
			paths = append(paths, path)
		}
		return err
	})
	logger.Error("Error occured while walking through directory recursively - ", err)

	return paths, nil
}

var tsconfigData = internal.JSONUnmarshal("./tsconfig.json")
var tsconfigCompilerOptions = tsconfigData["compilerOptions"].(map[string]interface{})
var packagePaths = tsconfigCompilerOptions["paths"].(map[string]interface{})
var OutDir = tsconfigCompilerOptions["outDir"].(string)

func remapFile(filePath string) {
	fileText, err := os.ReadFile(filePath)
	logger.Error("Error occured while reading "+filePath+" - ", err)

	var newFile string
	remapped := false

	importRegex := regexp.MustCompile("import.+[\"'](.+)[\"']")

	for _, fileLine := range strings.Split(string(fileText), "\r\n") {
		var newPackageName string
		if importRegex.MatchString(fileLine) {
			packageMatches := importRegex.FindStringSubmatch(fileLine)
			for key, value := range packagePaths {
				packagePath := value.([]interface{})[0].(string)
				if key == packageMatches[1] {
					rel, _ := filepath.Rel(filePath, OutDir+"/"+packagePath)
					packageRoute := strings.Replace(strings.ReplaceAll(rel, "\\", "/"), "../", "./", 1)
					newPackageName = strings.Replace(fileLine, packageMatches[1], packageRoute, 1)
					remapped = true
				}
			}
		}
		if len(newPackageName) > 0 {
			newFile += newPackageName + "\r\n"
		} else {
			newFile += fileLine + "\r\n"
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	logger.Error("Error occured while opening "+filePath+" - ", err)

	defer file.Close()

	_, err = file.WriteAt([]byte(newFile), 0)
	logger.Error("Error occured while editing file "+filePath+" - ", err)

	if remapped {
		logger.Success("Successfully remapped import(s) for " + strings.ReplaceAll(filePath, "\\", "/"))
	}
}

func ImportRemap() {
	distPaths, _ := walkDirectory("./" + OutDir)
	for _, distPath := range distPaths {
		remapFile(distPath)
	}
}
