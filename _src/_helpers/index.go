package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
)

func ConfigFileRead[T comparable](environment enums.Environment, applicationConfigChannel chan<- T, errorChannel chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	foundConfigFile, err := os.Open(filepath.Join(currentWorkingDirectory, "Configurations", fmt.Sprintf("%s.config.json", environment.String())))
	if err != nil {
		panic(err)
	}
	defer foundConfigFile.Close()

	var deserializedConfigFileData T
	if err := json.NewDecoder(foundConfigFile).Decode(&deserializedConfigFileData); err != nil {
		errorChannel <- err
	}
	applicationConfigChannel <- deserializedConfigFileData
	close(applicationConfigChannel)
	close(errorChannel)
}
