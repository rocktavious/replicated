package util

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/replicatedhq/libyaml"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

type HelmConverter interface {
	ConvertValues(helmValues string) (string, []libyaml.ConfigGroup, error)
}

func NewHelmConverter() HelmConverter {
	return &helmConverter{}
}

var _ HelmConverter = &helmConverter{}

type helmConverter struct {
}

func (c helmConverter) ConvertValues(in string) (string, []libyaml.ConfigGroup, error) {
	var values yaml.MapSlice
	err := yaml.Unmarshal([]byte(in), &values)
	if err != nil {
		return "", nil, errors.Wrap(err, "unmarshal values")
	}

	var path []string

	result, configItems := c.convertValuesRec(values, path)

	marshalled, err := yaml.Marshal(result)
	if err != nil {
		return "", nil, errors.Wrap(err, "marshal result")
	}
	return string(marshalled), []libyaml.ConfigGroup{{
		Items: configItems,
		Name:  "values",
	}}, nil
}

func (c helmConverter) convertValuesRec(in yaml.MapSlice, path []string) (yaml.MapSlice, []*libyaml.ConfigItem) {
	var valuesYAMLAcc yaml.MapSlice
	var configItems []*libyaml.ConfigItem

	for _, item := range in {
		key, ok := item.Key.(string);
		if !ok {
			// skip non-string keys (log me?)
			continue
		}

		newPath := append(path, key)

		configItemName := strings.Join(newPath, ".")
		configItemTitle := configItemName

		appendScalar := func() {
			valuesYAMLAcc = append(valuesYAMLAcc, yaml.MapItem{Key: key,
				Value: fmt.Sprintf(
					"{{repl ConfigOption %q }}",
					configItemName,
				),
			})
		}

		// todo support for more types
		switch typedValue := item.Value.(type) {
		case int:
			appendScalar()
			configItems = append(configItems, &libyaml.ConfigItem{
				Name:    configItemName,
				Title:   configItemTitle,
				Default: strconv.Itoa(typedValue),
				Type:    "text",
			})
		case string:
			appendScalar()
			configItems = append(configItems, &libyaml.ConfigItem{
				Name:    configItemName,
				Title:   configItemTitle,
				Default: typedValue,
				Type:    "text",
			})
		case bool:
			appendScalar()
			configItems = append(configItems, &libyaml.ConfigItem{
				Name:    configItemName,
				Title:   configItemTitle,
				Default: fmt.Sprintf("%v", typedValue),
				Type:    "bool",
			})
		case yaml.MapSlice:
			value, items := c.convertValuesRec(typedValue, newPath)
			valuesYAMLAcc = append(valuesYAMLAcc, yaml.MapItem{
				Key:   key,
				Value: value,
			})
			configItems = append(configItems, items...)
		case []interface{}:
		default:
			// todo need a real logger here
			fmt.Printf("Unsupported value type \"%T\" at %q, using default.", configItemName, typedValue)
			valuesYAMLAcc = append(valuesYAMLAcc, item)
		}
	}

	return valuesYAMLAcc, configItems
}
