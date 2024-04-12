package enums

import (
	"bytes"
	"errors"
	"fmt"
)

type Environment uint8

const (
	DevelopmentEnvironment Environment = iota
	StagingEnvironment
	ProductionEnvironment
)

func (env Environment) String() string {
	switch env {
	case DevelopmentEnvironment:
		return "development"
	case StagingEnvironment:
		return "staging"
	case ProductionEnvironment:
		return "production"
	default:
		return "development"
	}
}
func (env *Environment) Set(valueToSet string) error {
	if env == nil {
		return errors.New("enum cannot be nil")
	}
	valueAsBytes := []byte(valueToSet)
	if !env.unmarshalText(valueAsBytes) && !env.unmarshalText(bytes.ToLower(valueAsBytes)) {
		return fmt.Errorf(`unrecognized enum: %q`, valueAsBytes)
	}
	return nil
}
func (env *Environment) unmarshalText(text []byte) bool {
	switch string(text) {
	case "development", "DEVELOPMENT", "":
		*env = DevelopmentEnvironment
	case "staging", "STAGING":
		*env = StagingEnvironment
	case "production", "PRODUCTION":
		*env = ProductionEnvironment
	default:
		return false
	}
	return true
}
