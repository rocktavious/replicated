package cmd

import (
	"github.com/pkg/errors"
	"github.com/replicatedhq/replicated/cli/print"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
)

func (r *runners) InitImportHelmChart(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "import-helm-chart [PATH]",
		Short: "Create a replicated release from a locally checked-out helm chart.",
		Long:  `Create a replicated release from a locally checked-out helm chart.

Requires the "helm" command be present on your machine. Will run the following helm commands as part of release generation:

    helm init --client-only 
    helm dependency update PATH
    helm template

NOTE: This command will destroy the YAML file at replicated.yaml. Ensure to back up or commit any work in replicated.yaml before running import-helm-chart
`,
	}

	parent.AddCommand(cmd)
	cmd.RunE = r.importHelmChart
}

func (r *runners) importHelmChart(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Please provide a path to a helm chart to import")
	}

	chartRoot := args[0]

	bytes, err := ioutil.ReadFile(path.Join(chartRoot, "values.yaml"))
	if err != nil {
		return errors.Wrap(err, "read values file %q")
	}

	newValues, configGroups, err := r.helmConverter.ConvertValues(string(bytes))
	if err != nil {
		return errors.Wrap(err, "convert Helm Values")
	}

	k8sYAML, err := r.helmConverter.HelmTemplate(chartRoot, newValues)
	if err != nil {
		return errors.Wrap(err, "helm template")
	}
	fullYaml, err := r.helmConverter.BuildFullReplicatedYAML(chartRoot, k8sYAML, configGroups)
	if err != nil {
		return errors.Wrap(err, "build full YAML")
	}

	return print.HelmValues(os.Stdout, fullYaml)
}
