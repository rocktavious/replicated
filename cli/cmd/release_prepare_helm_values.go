package cmd

import (
	"github.com/pkg/errors"
	"github.com/replicatedhq/libyaml"
	"github.com/replicatedhq/replicated/cli/print"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func (r *runners) InitPrepareHelmValues(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "prepare-helm-values [PATH]",
		Short: "Convert a helm values.yaml to be used with a replicated release",
		Long:  `Convert a helm values.yaml to be used with a replicated release. 

PATH is optional, will default to reading a values file from a 
"values.yaml" in the current working directory.

`,
	}

	parent.AddCommand(cmd)
	cmd.RunE = r.prepareHelmValues
}

func (r *runners) prepareHelmValues(cmd *cobra.Command, args []string) error {
	valuesPath := "values.yaml"
	if len(args) == 1 {
		valuesPath = args[0]
	}

	bytes, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		return errors.Wrap(err, "read values file %q")
	}

	newValues, configGroups, err := r.helmConverter.ConvertValues(string(bytes))
	if err != nil {
		return errors.Wrap(err, "convert Helm Values")
	}
	doc := libyaml.RootConfig{
		ConfigGroups: configGroups,
	}

	marshalled, err := yaml.Marshal(doc)
	if err != nil {
		return errors.Wrap(err, "marshall replicated YAML")
	}

	return print.HelmValues(os.Stdout, `
---
# kind: replicated
` + string(marshalled) + `
---
# kind: helm-values
` + newValues)
}
