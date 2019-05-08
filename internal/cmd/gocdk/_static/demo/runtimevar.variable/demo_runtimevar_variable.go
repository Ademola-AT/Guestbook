package main

import (
	"context"
	"html/template"
	"net/http"
	"os"

	"gocloud.dev/runtimevar"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
)

// TODO(rvangent): This file is user-visible, add many comments explaining
// how it works.

func init() {
	http.HandleFunc("/demo/runtimevar.variable/", runtimevarVariableHandler)
}

var variableURL string
var variable *runtimevar.Variable
var variableErr error

func init() {
	variableURL = os.Getenv("RUNTIMEVAR_VARIABLE_URL")
	if variableURL == "" {
		variableURL = "constant://?val=my-variable&decoder=string"
	}
	variable, variableErr = runtimevar.OpenVariable(context.Background(), variableURL)
}

type runtimevarVariableData struct {
	URL      string
	Err      error
	Snapshot *runtimevar.Snapshot
}

// Input: *runtimevarVariableData.
const runtimevarVariableTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>runtimevar.Variable demo</title>
</head>
<body>
  <p>
    This page demonstrates the use of a Go CDK runtimevar.Variable.
  </p>
  <p>
    It is currently using a runtimevar.Variable based on the URL "{{ .URL }}", which
    can be configured via the environment variable "RUNTIMEVAR_VARIABLE_URL".
  </p>
  <p>
  See <a href="https://gocloud.dev/concepts/urls/">here</a> for more
  information about URLs in Go CDK APIs.
  </p>
  {{if .Err}}
    <p><strong>{{ .Err }}</strong></p>
  {{end}}
  {{if .Snapshot}}
    <div>The current value of the variable is:</div>
    <textarea rows="5" cols="40" readonly="true">{{ .Snapshot.Value }}</textarea>
    <div>It was last modified at: {{ .Snapshot.UpdateTime }}.</div>
  {{end}}
</body>
</html>
`

var runtimeVariableTmpl = template.Must(template.New("runtimevar.Variable").Parse(runtimevarVariableTemplate))

func runtimevarVariableHandler(w http.ResponseWriter, req *http.Request) {
	input := &runtimevarVariableData{
		URL: variableURL,
	}
	defer func() {
		if err := runtimeVariableTmpl.Execute(w, input); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	if variableErr != nil {
		input.Err = variableErr
		return
	}
	snapshot, err := variable.Latest(req.Context())
	if err != nil {
		input.Err = err
		return
	}
	input.Snapshot = &snapshot
}