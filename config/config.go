package config

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

type decoder func(file *os.File, destination interface{}) error

var decoderTypes map[string]decoder = map[string]decoder{
	"json": decodeJSON,
	"yaml": decodeYAML,
}

// FromFile inititalize destination config structure from file with extension  located at path
func FromFile(path string, destination interface{}) error {
	extension, err := getExtension(path)
	if err != nil {
		return errors.Wrap(err, "failed to get file extension")
	}

	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}

	decoder, ok := decoderTypes[extension]
	if !ok {
		return errors.Errorf("file must have .json/.yaml extension")
	}

	return decoder(file, destination)
}

func getExtension(path string) (string, error) {
	pathSplitted := strings.Split(path, ".")
	pathSplittedLen := len(pathSplitted)
	if pathSplittedLen < 2 {
		return "", errors.Errorf("invalid path: %s", path)
	}

	return pathSplitted[pathSplittedLen-1], nil
}

func decodeJSON(source *os.File, destination interface{}) error {
	decoder := json.NewDecoder(source)
	decoder.DisallowUnknownFields()

	return decoder.Decode(destination)
}

func decodeYAML(source *os.File, destination interface{}) error {
	decoder := yaml.NewDecoder(source)
	return decoder.Decode(destination)
}
