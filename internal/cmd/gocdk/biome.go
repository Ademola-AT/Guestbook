// Copyright 2019 The Go Cloud Development Kit Authors
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/xerrors"
)

const biomeConfigFileName = "biome.json"

// biomeConfig is the parsed configuration from a biome.json file.
type biomeConfig struct {
	ServeEnabled *bool   `json:"serve_enabled,omitempty"`
	Launcher     *string `json:"launcher,omitempty"`
}

// findBiomeDir returns the path to the named biome.
func findBiomeDir(moduleRoot, name string) string {
	return filepath.Join(moduleRoot, "biomes", name)
}

// readBiomeConfig reads and parses the biome configuration from the filesystem.
// If the configuration file could not be found, readBiomeConfig returns an
// error for which xerrors.As(err, new(*biomeNotFoundError)) returns true.
func readBiomeConfig(moduleRoot, biome string) (*biomeConfig, error) {
	configPath := filepath.Join(findBiomeDir(moduleRoot, biome), biomeConfigFileName)
	data, err := ioutil.ReadFile(configPath)
	if os.IsNotExist(err) {
		// TODO(light): Wrap error for formatting chain but not unwrap chain.
		notFound := &biomeNotFoundError{
			moduleRoot: moduleRoot,
			biome:      biome,
			frame:      xerrors.Caller(0),
			detail:     err,
		}
		return nil, xerrors.Errorf("read biome %s configuration: %w", biome, notFound)
	}
	if err != nil {
		return nil, xerrors.Errorf("read biome %s configuration: %w", err)
	}
	config := new(biomeConfig)
	if err := json.Unmarshal(data, config); err != nil {
		return nil, xerrors.Errorf("read biome %s configuration: %w", err)
	}
	return config, nil
}

// tfReadOutput runs `terraform output` on the given directory and returns
// the parsed result.
func tfReadOutput(ctx context.Context, dir string, env []string) (map[string]*tfOutput, error) {
	c := exec.CommandContext(ctx, "terraform", "output", "-json")
	c.Dir = dir
	c.Env = env
	data, err := c.Output()
	if err != nil {
		return nil, xerrors.Errorf("read terraform output: %w", err)
	}
	var parsed map[string]*tfOutput
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, xerrors.Errorf("read terraform output: %w", err)
	}
	return parsed, nil
}

// tfOutput describes a single output value.
type tfOutput struct {
	Type      string      `json:"type"` // one of "string", "list", or "map"
	Sensitive bool        `json:"sensitive"`
	Value     interface{} `json:"value"`
}

// stringValue returns the output's value if it is a string.
func (out *tfOutput) stringValue() string {
	if out == nil {
		return ""
	}
	v, _ := out.Value.(string)
	return v
}

// mapValue returns the output's value if it is a map.
func (out *tfOutput) mapValue() map[string]interface{} {
	if out == nil {
		return nil
	}
	v, _ := out.Value.(map[string]interface{})
	return v
}

// biomeNotFoundError is an error returned when a biome cannot be found.
type biomeNotFoundError struct {
	moduleRoot string
	biome      string
	frame      xerrors.Frame
	detail     error
}

func (e *biomeNotFoundError) Error() string {
	return fmt.Sprintf("biome %s not found", e.biome)
}

func (e *biomeNotFoundError) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())
	if !p.Detail() {
		return nil
	}
	p.Printf("biome = %q", findBiomeDir(e.moduleRoot, e.biome))
	e.frame.Format(p)
	return e.detail
}

func (e *biomeNotFoundError) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
}
