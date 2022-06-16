package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/notbeer/typescript-transpile-remap/internal"
	"github.com/notbeer/typescript-transpile-remap/logger"
	"github.com/notbeer/typescript-transpile-remap/tools"
)

func main() {
	buildData := internal.JSONUnmarshal("./configs/build.json")
	deleteScriptsDir := buildData["buildOptions"].(map[string]interface{})["deleteScriptsDir"].(bool)

	if deleteScriptsDir {
		os.RemoveAll(tools.TsconfigCompilerOptions["outDir"].(string))
	}

	log.Println("Transpiling source folder...")
	output, err := exec.Command("tsc").Output()
	logger.Error(string(output), err)
	logger.Success("Transpiling complete")

	log.Println("Starting to remap files...")
	tools.ImportRemap()
	logger.Success("Mapping complete")

	fmt.Println("Press any key to exit...")
	var input string
	fmt.Scanf(input, "%s")
}
