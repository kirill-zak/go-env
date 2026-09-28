package variables

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/kirill-zak/go-env/entity"
)

// ReadFromYAML reads variables from YAML reader and converts it to Variable representation.
func ReadFromYAML(r io.Reader) ([]entity.Variable, error) {
	var pld struct {
		Env map[string]yamlVariable `yaml:"env"`
	}
	if err := yaml.NewDecoder(r).Decode(&pld); err != nil {
		return nil, fmt.Errorf("read yaml file failed: %w", err)
	}
	if pld.Env == nil {
		return nil, errors.New("provided yaml is missing required fields")
	}

	res := make([]entity.Variable, 0, len(pld.Env))
	for name, e := range pld.Env {
		res = append(res, e.toVariable(name))
	}

	return res, nil
}

// ReadSectionFromYAML reads variables from YAML reader and converts them to Section representation.
func ReadSectionFromYAML(r io.Reader) ([]entity.Section, error) {
	var pld struct {
		Env yaml.Node `yaml:"env"`
	}
	if err := yaml.NewDecoder(r).Decode(&pld); err != nil {
		return nil, fmt.Errorf("read configuration file failed: %w", err)
	}
	if pld.Env.Kind != yaml.MappingNode {
		return nil, errors.New("provided yaml is missing required fields")
	}

	var (
		res                          = make([]entity.Section, 0)
		currentSection               entity.Section
		currentValue, currentComment string
	)

	for _, node := range pld.Env.Content {
		switch node.Kind { //nolint:exhaustive // yaml logic
		case yaml.ScalarNode:
			currentValue = node.Value
			currentComment = node.HeadComment
		case yaml.MappingNode:
			// decode variable
			var currentVariable yamlVariable
			if err := node.Decode(&currentVariable); err != nil {
				return nil, fmt.Errorf("decode node failed: %w", err)
			}

			// solve current section
			newSectionName := strings.Trim(currentComment, "# \n")
			if newSectionName != "" {
				if len(currentSection.Variables) > 0 {
					res = append(res, currentSection)
				}
				currentSection = entity.Section{Name: newSectionName}
			}

			// append variable to section
			currentSection.Variables = append(
				currentSection.Variables, currentVariable.toVariable(currentValue),
			)
		default:
			return nil, errors.New("invalid format")
		}
	}

	// append last section
	res = append(res, currentSection)

	return res, nil
}
