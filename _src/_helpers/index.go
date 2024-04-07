package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ConfigFileRead[T comparable](environment string) T {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	foundConfigFile, err := os.Open(filepath.Join(currentWorkingDirectory, "Configurations", fmt.Sprintf("%s.config.json", environment)))
	if err != nil {
		panic(err)
	}
	defer foundConfigFile.Close()

	var deserializedConfigFileData T
	if err := json.NewDecoder(foundConfigFile).Decode(&deserializedConfigFileData); err != nil {
		panic(err)
	}
	return deserializedConfigFileData
}
