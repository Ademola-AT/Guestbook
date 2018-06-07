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

package runtimeconfigurator

import (
	"context"
	"flag"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cloud/gcp"
	"github.com/google/go-cloud/runtimevar"
	"github.com/google/go-cloud/runtimevar/driver"
	"github.com/google/go-cloud/testing/replay"
	"github.com/google/go-cmp/cmp"
	pb "google.golang.org/genproto/googleapis/cloud/runtimeconfig/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"
)

const (
	// config is the runtimeconfig high-level config that variables sit under.
	config      = "runtimeconfigurator_test"
	description = "Config for test variables created by runtimeconfigurator_test.go"
)

var projectID = flag.String("project", "", "GCP project ID (string, not project number) to run tests against")

// Ensure that watcher implements driver.Watcher.
var _ driver.Watcher = &watcher{}

// fakeServer partially implements runtimevarManagerServer for Client to connect to.  Prefill
// responses field with the ordered list of responses to GetVariable calls.
type fakeServer struct {
	pb.RuntimeConfigManagerServer
	responses []response
	index     int
}

type response struct {
	vrbl *pb.Variable
	err  error
}

func (s *fakeServer) GetVariable(context.Context, *pb.GetVariableRequest) (*pb.Variable, error) {
	if len(s.responses) == 0 {
		return nil, fmt.Errorf("fakeClient missing responses")
	}
	resp := s.responses[s.index]
	// Adjust index to next response for next call till it gets to last one, then keep using the
	// last one.
	if s.index < len(s.responses)-1 {
		s.index++
	}
	return resp.vrbl, resp.err
}

func setUp(t *testing.T, fs *fakeServer) (*Client, func()) {
	t.Helper()
	// Set up gRPC server.
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("tcp listen failed: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRuntimeConfigManagerServer(s, fs)
	// Run gRPC server on a background goroutine.
	go s.Serve(lis)

	// Set up client.
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	client := NewClient(pb.NewRuntimeConfigManagerClient(conn))
	return client, func() {
		conn.Close()
		s.Stop()
	}
}

type jsonData struct {
	Hello string `json:"hello"`
}

var (
	// Set wait timeout used for tests.
	watchOpt = &WatchOptions{
		WaitTime: 100 * time.Millisecond,
	}
	resourceName = ResourceName{
		ProjectID: "ID42",
		Config:    "config",
		Variable:  "greetings",
	}
	startTime = time.Now().Unix()
	jsonVar1  = &pb.Variable{
		Name:       "greetings",
		Contents:   &pb.Variable_Text{Text: `{"hello": "hello"}`},
		UpdateTime: &tspb.Timestamp{Seconds: startTime},
	}
	jsonVar2 = &pb.Variable{
		Name:       "greetings",
		Contents:   &pb.Variable_Value{Value: []byte(`{"hello": "hola"}`)},
		UpdateTime: &tspb.Timestamp{Seconds: startTime + 100},
	}
	jsonDataPtr *jsonData
)

func TestInitialStringWatch(t *testing.T) {
	ctx := context.Background()

	client, done, err := newConfigClient(ctx, t.Logf, "initial-string-watch.replay")
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	rn := ResourceName{
		ProjectID: *projectID,
		Config:    config,
		desc:      description,
		Variable:  "TestStringWatch",
	}

	want := "facepalm: 🤦"
	_, err = createStringVariable(ctx, client.client, rn, want)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := deleteConfig(ctx, client.client, rn); err != nil {
			t.Fatalf("delete config failed, possible human cleanup required: %v", err)
		}
	}()

	variable, err := client.NewVariable(ctx, rn, runtimevar.StringDecoder, nil)
	if err != nil {
		t.Fatal(err)
	}

	got, err := variable.Watch(ctx)
	if err != nil {
		t.Fatalf("got error %v; want nil", err)
	}
	if diff := cmp.Diff(got.Value, want); diff != "" {
		t.Errorf("got diff %v; want nil", diff)
	}
}

func TestInitialJSONWatch(t *testing.T) {
	ctx := context.Background()

	client, done, err := newConfigClient(ctx, t.Logf, "initial-json-watch.replay")
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	rn := ResourceName{
		ProjectID: *projectID,
		Config:    config,
		desc:      description,
		Variable:  "TestJSONWatch",
	}

	type home struct {
		Person string `json:"Person"`
		Home   string `json:"Home"`
	}
	var jsonDataPtr *home
	want := &home{"Batman", "Gotham"}
	_, err = createByteVariable(ctx, client.client, rn, []byte(`{"Person": "Batman", "Home": "Gotham"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := deleteConfig(ctx, client.client, rn); err != nil {
			t.Logf("delete config failed, possible human cleanup required: %v", err)
		}
	}()

	variable, err := client.NewVariable(ctx, rn, runtimevar.NewDecoder(jsonDataPtr, runtimevar.JSONDecode), nil)
	if err != nil {
		t.Fatal(err)
	}

	got, err := variable.Watch(ctx)
	if err != nil {
		t.Fatalf("got error %v; want nil", err)
	}
	if diff := cmp.Diff(got.Value.(*home), want); diff != "" {
		t.Errorf("got diff %v; want nil", diff)
	}
}

func TestWatch(t *testing.T) {
	client, cleanUp := setUp(t, &fakeServer{
		responses: []response{
			{vrbl: jsonVar1},
			{vrbl: jsonVar2},
		},
	})
	defer cleanUp()

	ctx := context.Background()
	variable, err := client.NewVariable(ctx, resourceName, runtimevar.NewDecoder(jsonDataPtr, runtimevar.JSONDecode), watchOpt)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}

	got1, err := variable.Watch(ctx)
	if err != nil {
		t.Fatalf("Variable.Watch returned error: %v", err)
	}
	if diff := cmp.Diff(got1.Value.(*jsonData), &jsonData{"hello"}); diff != "" {
		t.Errorf("Snapshot.Value: %s", diff)
	}

	got2, err := variable.Watch(ctx)
	if err != nil {
		t.Fatalf("Variable.Watch returned error: %v", err)
	}
	if diff := cmp.Diff(got2.Value.(*jsonData), &jsonData{"hola"}); diff != "" {
		t.Errorf("Snapshot.Value: %s", diff)
	}
}

func TestCustomDecode(t *testing.T) {
	value := "hello world"
	strVar := &pb.Variable{
		Name:       "greetings",
		Contents:   &pb.Variable_Value{Value: []byte(value)},
		UpdateTime: &tspb.Timestamp{Seconds: startTime},
	}

	client, cleanUp := setUp(t, &fakeServer{
		responses: []response{
			{vrbl: strVar},
		},
	})
	defer cleanUp()

	ctx := context.Background()
	watchOpt := &WatchOptions{
		WaitTime: 500 * time.Millisecond,
	}
	variable, err := client.NewVariable(ctx, resourceName, runtimevar.NewDecoder("", stringDecode), watchOpt)
	if err != nil {
		t.Fatalf("Client.NewConfig returned error: %v", err)
	}

	got, err := variable.Watch(ctx)
	if err != nil {
		t.Fatalf("Variable.Watch returned error: %v", err)
	}
	if diff := cmp.Diff(got.Value.(string), value); diff != "" {
		t.Errorf("Snapshot.Value: %s", diff)
	}
}

func stringDecode(b []byte, obj interface{}) error {
	// obj is a pointer to a string.
	v := reflect.ValueOf(obj).Elem()
	v.SetString(string(b))
	return nil
}

func TestWatchCancelledBeforeFirstWatch(t *testing.T) {
	client, cleanUp := setUp(t, &fakeServer{
		responses: []response{
			{vrbl: jsonVar1},
		},
	})
	defer cleanUp()

	ctx := context.Background()
	variable, err := client.NewVariable(ctx, resourceName, runtimevar.NewDecoder(jsonDataPtr, runtimevar.JSONDecode), watchOpt)
	if err != nil {
		t.Fatalf("Client.NewConfig returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, err = variable.Watch(ctx)
	if err == nil {
		t.Fatal("Variable.Watch returned nil error, expecting an error from cancelling")
	}
}

func TestContextCancelledInBetweenWatchCalls(t *testing.T) {
	client, cleanUp := setUp(t, &fakeServer{
		responses: []response{
			{vrbl: jsonVar1},
		},
	})
	defer cleanUp()

	ctx := context.Background()
	variable, err := client.NewVariable(ctx, resourceName, runtimevar.NewDecoder(jsonDataPtr, runtimevar.JSONDecode), watchOpt)
	if err != nil {
		t.Fatalf("Client.NewConfig returned error: %v", err)
	}

	_, err = variable.Watch(ctx)
	if err != nil {
		t.Fatalf("Variable.Watch returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, err = variable.Watch(ctx)
	if err == nil {
		t.Fatal("Variable.Watch returned nil error, expecting an error from cancelling")
	}
}

func TestWatchDeletedAndReset(t *testing.T) {
	client, cleanUp := setUp(t, &fakeServer{
		responses: []response{
			{vrbl: jsonVar1},
			{err: status.Error(codes.NotFound, "deleted")},
			{vrbl: jsonVar2},
		},
	})
	defer cleanUp()

	ctx := context.Background()
	variable, err := client.NewVariable(ctx, resourceName, runtimevar.NewDecoder(jsonDataPtr, runtimevar.JSONDecode), watchOpt)
	if err != nil {
		t.Fatalf("Client.NewConfig() returned error: %v", err)
	}

	prev, err := variable.Watch(ctx)
	if err != nil {
		t.Fatalf("Variable.Watch returned error: %v", err)
	}

	// Expect deleted error.
	if _, err := variable.Watch(ctx); err == nil {
		t.Fatalf("Variable.Watch returned nil, want error")
	}

	// Calling Watch again will poll for jsonVar2.
	got, err := variable.Watch(ctx)
	if err != nil {
		t.Fatalf("Variable.Watch returned error: %v", err)
	}
	if diff := cmp.Diff(got.Value.(*jsonData), &jsonData{"hola"}); diff != "" {
		t.Errorf("Snapshot.Value: %s", diff)
	}
	if !got.UpdateTime.After(prev.UpdateTime) {
		t.Errorf("Snapshot.UpdateTime is less than or equal to previous value")
	}
}

func newConfigClient(ctx context.Context, logf func(string, ...interface{}), filepath string) (*Client, func(), error) {
	creds, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, nil, err
	}

	mode := recorder.ModeRecording
	if testing.Short() {
		mode = recorder.ModeReplaying
	}

	rOpts, done, err := replay.NewGCPDialOptions(logf, mode, filepath, scrubber)
	if err != nil {
		return nil, nil, err
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithPerRPCCredentials(oauth.TokenSource{gcp.CredentialsTokenSource(creds)}),
	}
	opts = append(opts, rOpts...)
	conn, err := grpc.DialContext(ctx, endPoint, opts...)
	if err != nil {
		return nil, nil, err
	}

	return NewClient(pb.NewRuntimeConfigManagerClient(conn)), done, nil
}

// createConfig creates a fresh config. It will always overwrite any previous configuration,
// thus it is not thread safe.
func createConfig(ctx context.Context, client pb.RuntimeConfigManagerClient, rn ResourceName) (*pb.RuntimeConfig, error) {
	// No need to handle this error; either the config doesn't exist (good) or the test
	// will fail on the create step and requires human intervention anyway.
	_ = deleteConfig(ctx, client, rn)
	return client.CreateConfig(ctx, &pb.CreateConfigRequest{
		Parent: "projects/" + rn.ProjectID,
		Config: &pb.RuntimeConfig{
			Name:        rn.configPath(),
			Description: rn.desc,
		},
	})
}

func deleteConfig(ctx context.Context, client pb.RuntimeConfigManagerClient, rn ResourceName) error {
	_, err := client.DeleteConfig(ctx, &pb.DeleteConfigRequest{
		Name: rn.configPath(),
	})

	return err
}

func createByteVariable(ctx context.Context, client pb.RuntimeConfigManagerClient, rn ResourceName, value []byte) (*pb.Variable, error) {
	if _, err := createConfig(ctx, client, rn); err != nil {
		return nil, fmt.Errorf("unable to create parent config for %+v: %v", rn, err)
	}

	return client.CreateVariable(ctx, &pb.CreateVariableRequest{
		Parent: rn.configPath(),
		Variable: &pb.Variable{
			Name:     rn.String(),
			Contents: &pb.Variable_Value{Value: value},
		},
	})
}

func createStringVariable(ctx context.Context, client pb.RuntimeConfigManagerClient, rn ResourceName, str string) (*pb.Variable, error) {
	if _, err := createConfig(ctx, client, rn); err != nil {
		return nil, fmt.Errorf("unable to create parent config for %+v: %v", rn, err)
	}

	return client.CreateVariable(ctx, &pb.CreateVariableRequest{
		Parent: rn.configPath(),
		Variable: &pb.Variable{
			Name:     rn.String(),
			Contents: &pb.Variable_Text{Text: str},
		},
	})

}

type fakeProto struct{}

func (p *fakeProto) Reset()         {}
func (p *fakeProto) String() string { return "fake" }
func (p *fakeProto) ProtoMessage()  {}

func TestScrubber(t *testing.T) {
	var tests = []struct {
		name      string
		msg, want proto.Message
		wantErr   bool
	}{
		{
			name: "Messages that match the regexp should have project IDs redacted",
			msg: &pb.DeleteConfigRequest{
				Name: "projects/project_id/name",
			},
			want: &pb.DeleteConfigRequest{
				Name: "projects/REDACTED/name",
			},
		},
		{
			name: "Messages that have nested strings where project IDs can be found should all be redacted",
			msg: &pb.CreateConfigRequest{
				Parent: "/projects/project_id/parent",
				Config: &pb.RuntimeConfig{
					Name: "projects/project_id/config/name",
				},
			},
			want: &pb.CreateConfigRequest{
				Parent: "/projects/REDACTED/parent",
				Config: &pb.RuntimeConfig{
					Name: "projects/REDACTED/config/name",
				},
			},
		},
		{
			name: "Messages that don't match the regexp should be returned unchanged",
			msg: &pb.DeleteConfigRequest{
				Name: "project_id/name",
			},
			want: &pb.DeleteConfigRequest{
				Name: "project_id/name",
			},
		},
		{
			name: "Empty messages should be returned unchanged",
			msg:  &empty.Empty{},
			want: &empty.Empty{},
		},
		{
			name:    "Unknown messages should return an error",
			msg:     &fakeProto{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := proto.Clone(tc.msg)
			err := scrubber(t.Logf, "", got)

			switch {
			case err != nil && !tc.wantErr:
				t.Fatal(err)
			case err == nil && tc.wantErr:
				t.Errorf("want error; got nil")
			case err != nil && tc.wantErr:
				// Got error as expected, test passed.
				return
			case !cmp.Equal(got, tc.want):
				t.Errorf("got %s; want %s", got, tc.want)
			}
		})
	}
}

func scrubber(logf func(string, ...interface{}), _ string, msg proto.Message) error {
	// Example matches:
	// projects/foobar
	// /projects/foobar/baz
	re := regexp.MustCompile(`(?U)(\/?projects\/)(.*)(\/|$)`)
	// Without the curly braces, Go interprets the group as named $1REDACTED which
	// doesn't match anything.
	replacePattern := "${1}REDACTED${3}"
	logf("Proto begins as %s", msg)

	switch m := msg.(type) {
	case *pb.DeleteConfigRequest:
		m.Name = re.ReplaceAllString(m.GetName(), replacePattern)
	case *pb.CreateConfigRequest:
		m.Parent = re.ReplaceAllString(m.GetParent(), replacePattern)
		m.Config.Name = re.ReplaceAllString(m.GetConfig().GetName(), replacePattern)
	case *pb.CreateVariableRequest:
		m.Parent = re.ReplaceAllString(m.GetParent(), replacePattern)
		m.Variable.Name = re.ReplaceAllString(m.GetVariable().GetName(), replacePattern)
	case *pb.GetVariableRequest:
		m.Name = re.ReplaceAllString(m.GetName(), replacePattern)
	case *pb.RuntimeConfig:
		m.Name = re.ReplaceAllString(m.GetName(), replacePattern)
	case *pb.Variable:
		m.Name = re.ReplaceAllString(m.GetName(), replacePattern)
	case *empty.Empty:
	default:
		return fmt.Errorf("unknown proto type, can't scrub: %v", reflect.TypeOf(msg))
	}

	logf("Proto ends as %s", msg)
	return nil
}
