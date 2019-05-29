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
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/spf13/cobra"
	"gocloud.dev/gcp"
	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"
)

// generate_static converts the files in _static/ into constants in a new
// file,
//go:generate go run generate_static.go -- static.go

func main() {
	pctx, err := osProcessContext()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)

	}
	ctx, done := withInterrupt(context.Background())
	err = run(ctx, pctx, os.Args[1:])
	done()
	if err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context, pctx *processContext, args []string) error {
	var rootCmd = &cobra.Command{
		Use:   "gocdk",
		Short: "TODO gocdk",
		Long:  "TODO more about gocdk",
		// We want to print usage for any command/flag parsing errors, but
		// suppress usage for random errors returned by a command's RunE.
		// This function gets called right before RunE, so suppress from now on.
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceUsage = true
		},
	}

	registerInitCmd(ctx, pctx, rootCmd)
	registerServeCmd(ctx, pctx, rootCmd)
	registerDemoCmd(ctx, pctx, rootCmd)
	registerBiomeCmd(ctx, pctx, rootCmd)
	registerDeployCmd(ctx, pctx, rootCmd)
	registerBuildCmd(ctx, pctx, rootCmd)
	registerApplyCmd(ctx, pctx, rootCmd)
	registerLaunchCmd(ctx, pctx, rootCmd)

	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

// processContext is the state that gocdk uses to run. It is collected in
// this struct to avoid obtaining this from globals for simpler testing.
type processContext struct {
	workdir string
	env     []string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

// osProcessContext returns the default process context from global variables.
func osProcessContext() (*processContext, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &processContext{
		workdir: workdir,
		env:     os.Environ(),
		stdin:   os.Stdin,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
	}, nil
}

// overrideEnv returns a copy of env that has vars appended to the end.
// It will not modify env's backing array.
func overrideEnv(env []string, vars ...string) []string {
	// Setting the slice's capacity to length ensures that a new backing array
	// is allocated if len(vars) > 0.
	return append(env[:len(env):len(env)], vars...)
}

// resolve joins path with pctx.workdir if path is relative. Otherwise,
// it returns path.
func (pctx *processContext) resolve(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(pctx.workdir, path)
}

// gcpCredentials returns the credentials to use for this process context.
func (pctx *processContext) gcpCredentials(ctx context.Context) (*google.Credentials, error) {
	// TODO(light): google.DefaultCredentials uses Getenv directly, so it is
	// difficult to disentangle to use processContext.
	return gcp.DefaultCredentials(ctx)
}

// findModuleRoot searches the given directory and those above it for the Go
// module root.
func findModuleRoot(ctx context.Context, dir string) (string, error) {
	c := exec.CommandContext(ctx, "go", "list", "-m", "-f", "{{.Dir}}")
	c.Dir = dir
	output, err := c.Output()
	if err != nil {
		return "", xerrors.Errorf("find module root for %s: %w", dir, err)
	}
	output = bytes.TrimSuffix(output, []byte("\n"))
	if len(output) == 0 {
		return "", xerrors.Errorf("find module root for %s: no module found", dir, err)
	}
	return string(output), nil
}

// withInterrupt returns a copy of parent with a new Done channel. The returned
// context's Done channel will be closed when the process receives an interrupt
// signal, the parent context's Done channel is closed, or the stop function is
// called, whichever comes first.
//
// The stop function releases resources and stops listening for signals, so code
// should call it as soon as the operation using the context completes.
func withInterrupt(parent context.Context) (_ context.Context, stop func()) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, interruptSignals()...)
	ctx, cancel := context.WithCancel(parent)
	done := make(chan struct{})
	go func() {
		select {
		case <-sig:
			cancel()
		case <-done:
		}
	}()
	return ctx, func() {
		cancel()
		signal.Stop(sig)
		close(done)
	}
}
