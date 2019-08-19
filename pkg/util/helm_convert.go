package util

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/replicatedhq/libyaml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type HelmConverter interface {
	ConvertValues(helmValues string) (string, []*libyaml.ConfigGroup, error)
	HelmTemplate(chartRoot, valuesYAML string) (string, error)
	BuildFullReplicatedYAML(chartRoot, k8sYAML string, configGroups []*libyaml.ConfigGroup) (string, error)
}

func NewHelmConverter() HelmConverter {
	return &helmConverter{}
}

var _ HelmConverter = &helmConverter{}

type helmConverter struct {
}

func (c helmConverter) HelmTemplate(chartRoot, valuesYAML string) (string, error) {

	err := c.helmInit()
	if err != nil {
		return "", errors.Wrap(err, "helm init")
	}

	err = c.helmDependencyUpdate(chartRoot)
	if err != nil {
		return "", errors.Wrap(err, "helm dependency update")
	}

	k8sYAML, err := c.helmTemplate(chartRoot, valuesYAML)
	if err != nil {
		return "", errors.Wrap(err, "helm template")
	}

	return k8sYAML, nil
}

func (c helmConverter) BuildFullReplicatedYAML(chartRoot, k8sYAML string, configGroups []*libyaml.ConfigGroup) (string, error) {
	k8sYAMLWithComments := strings.Join(strings.Split(k8sYAML, "\n---\n"), "\n---\n# kind: scheduler-kubernetes\n")
	chartYaml, err := ioutil.ReadFile(path.Join(chartRoot, "Chart.yaml"))
	if err != nil {
		return "", errors.Wrapf(err, "read Chart.yaml from %q", chartRoot)
	}

	chartInfo := make(map[string]interface{})
	err = yaml.Unmarshal(chartYaml, &chartInfo)
	if err != nil {
		return "", errors.Wrapf(err, "unmarshal Chart.yaml from %q", chartRoot)
	}

	doc := libyaml.RootConfig{
		APIVersion:   "2.38.0",
		ConfigGroups: configGroups,
		Name:         fmt.Sprintf("%v", chartInfo["name"]),
		Version:      fmt.Sprintf("%v", chartInfo["version"]),
		Properties: libyaml.Properties{
			ConsoleTitle: fmt.Sprintf("%v", chartInfo["name"]),
			LogoUrl:      fmt.Sprintf("%v", chartInfo["icon"]),
		},
	}
	serialized, err := yaml.Marshal(doc)
	if err != nil {
		return "", errors.Wrap(err, "serialize replicated yaml")
	}
	return string(serialized) + `\n\n---\n# kind: scheduler-kubernetes` + k8sYAMLWithComments, nil
}

func (c helmConverter) ConvertValues(in string) (string, []*libyaml.ConfigGroup, error) {
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
	return string(marshalled), []*libyaml.ConfigGroup{{
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
			if typedValue == nil {
				appendScalar()
				configItems = append(configItems, &libyaml.ConfigItem{
					Name:    configItemName,
					Title:   configItemTitle,
					Default: fmt.Sprintf("%v", typedValue),
					Type:    "text",
				})
			} else {
				// todo need a real logger here
				fmt.Fprint(os.Stderr, fmt.Sprintf("Unsupported value type \"%T\" at %q, using default.\n", configItemName, typedValue))
				valuesYAMLAcc = append(valuesYAMLAcc, item)
			}
		}
	}

	return valuesYAMLAcc, configItems
}

func (c helmConverter) helmInit() error {
	var stderr bytes.Buffer
	cmd := exec.Command("helm")
	cmd.Args = append(cmd.Args, "init", "--client-only")
	cmd.Stderr = &stderr
	stdOut, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, string(stdOut)+"\n"+stderr.String())
	}
	return nil
}

func (c helmConverter) helmDependencyUpdate(chartRoot string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("helm")
	cmd.Args = append(cmd.Args, "dependency", "update", chartRoot)
	cmd.Stderr = &stderr
	stdOut, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, string(stdOut)+"\n"+stderr.String())
	}
	return nil
}

func (c helmConverter) helmTemplate(chartRoot, valuesYAML string) (string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("helm")
	cmd.Args = append(cmd.Args, "template", "-f", "-", chartRoot)
	cmd.Stderr = &stderr
	cmd.Stdin = strings.NewReader(valuesYAML)
	stdOut, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, string(stdOut)+"\n"+stderr.String())
	}
	return string(stdOut), nil
}
