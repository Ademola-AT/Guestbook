// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The deploy program deploys the Guestbook app to GAE.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("gae/deploy: ")
	guestbookDir := flag.String("guestbook_dir", ".", "directory containing the guestbook example")
	tfStatePath := flag.String("tfstate", "terraform.tfstate", "path to terraform state file")
	flag.Parse()
	if err := deploy(*guestbookDir, *tfStatePath); err != nil {
		log.Fatal(err)
	}
}

func deploy(guestbookDir, tfStatePath string) error {
	// Extract some info from the tfstate file.
	type tfItem struct {
		Sensitive bool
		Type      string
		Value     string
	}
	type state struct {
		Project          tfItem
		Bucket           tfItem
		DatabaseInstance tfItem `json:"database_instance"`
		DatabaseRegion   tfItem `json:"database_region"`
		MotdVarConfig    tfItem `json:"motd_var_config"`
		MotdVarName      tfItem `json:"motd_var_name"`
	}
	tfStateb, err := runb("terraform", "output", "-state", tfStatePath, "-json")
	if err != nil {
		return err
	}
	var tfState state
	if err := json.Unmarshal(tfStateb, &tfState); err != nil {
		return fmt.Errorf("parsing terraform state JSON: %v", err)
	}

	// Fill out the params for app.yaml.
	var p Params
	p.Instance = fmt.Sprintf("%s:%s:%s", state.Project.Value, state.DatabaseRegion.Value, state.DatabaseInstance.Value)
	p.Password = fmt.Sprintf("%d", rand.Int())
	if err = setDbPassword(p.Password); err != nil {
		return err
	}
	p.Bucket = state.Bucket.Value

	// Write the app.yaml configuration file.
	t, err := template.New("appyaml").Parse(appYAML)
	if err != nil {
		return fmt.Errorf("parsing app.yaml template: %v", err)
	}
	f, err := os.Create("app.yaml")
	if err != nil {
		return err
	}
	defer f.Close()
	t.Execute(f, p)

	// Deploy the app to GAE.
	cmd := exec.Command("gcloud", "app", "deploy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running gcloud app deploy: %v", err)
	}

	return nil
}

type Params struct {
	Instance string
	Password string
	Bucket   string
}

const appYAML = `runtime: go111
env_variables:
  DB_USER: guestbook
  DB_INSTANCE: {{.Instance}}
  DB_DATABASE: guestbook
  DB_PASSWORD: {{.Password}}
  GUESTBOOK_BUCKET: {{.Bucket}}

`

type gcloud struct {
	projectID string
}

func (gcp *gcloud) cmd(args ...string) *exec.Cmd {
	args = append([]string{"--quiet", "--project", gcp.projectID}, args...)
	cmd := exec.Command("gcloud", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Stderr = os.Stderr
	return cmd
}

func run(args ...string) (stdout string, err error) {
	stdoutb, err := runb(args...)
	return strings.TrimSpace(string(stdoutb)), err
}

func runb(args ...string) (stdout []byte, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return stdoutb, nil
}

func setDbPassword(pw string) error {
	instances, err := getSQLInstances()
	if err != nil {
		return fmt.Errorf("getting GCP SQL instances: %v", err)
	}
	if len(instances) != 1 {
		return fmt.Errorf("got %d instances, want exactly one", len(instances))
	}
	user := "guestbook"
	if err := setPassword(instance, user, pw); err != nil {
		return fmt.Errorf("setting guestbook db user password: %v", err)
	}
	return nil
}

func getSQLInstances() ([]string, error) {
	out, err := run("gcloud", "sql", "instances", "list")
	if err != nil {
		return nil, fmt.Errorf("calling gcloud sql instances list: %v", err)
	}
	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return nil, errors.New("no lines returned from gcloud sql instances list")
	}
	var insts []string
	for _, l := range lines[1:] {
		parts := strings.Split(l, " ")
		if len(parts) == 0 {
			break
		}
		inst := parts[0]
		insts = append(insts, inst)
	}
	return insts, nil
}
