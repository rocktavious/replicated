package print

import (
	"github.com/replicatedhq/kots-lint/pkg/kots"
	"text/tabwriter"
	"text/template"
)

var lintTmplSrc = `RULE	TYPE	FILENAME	LINE	MESSAGE
{{ range . -}}
{{ .Rule }}	{{ .Type }}	{{ .Path }}	{{with .Positions}}{{ (index . 0).Start.Line }}{{else}}	{{end}}	{{ .Message}}	
{{ end }}`

var lintTmpl = template.Must(template.New("lint").Parse(lintTmplSrc))

func LintErrors(w *tabwriter.Writer, lintErrors []kots.LintExpression) error {
	if err := lintTmpl.Execute(w, lintErrors); err != nil {
		return err
	}
	return w.Flush()
}
