// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"database/sql"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/go-sql-driver/mysql"
	"go.opencensus.io/trace"
	"gocloud.dev/aws"
	"gocloud.dev/aws/rds"
	"gocloud.dev/blob"
	"gocloud.dev/blob/azureblob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/blob/s3blob"
	"gocloud.dev/gcp"
	"gocloud.dev/gcp/cloudsql"
	"gocloud.dev/mysql/cloudmysql"
	"gocloud.dev/mysql/rdsmysql"
	"gocloud.dev/requestlog"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/awsparamstore"
	"gocloud.dev/runtimevar/blobvar"
	"gocloud.dev/runtimevar/filevar"
	"gocloud.dev/runtimevar/gcpruntimeconfig"
	"gocloud.dev/server"
	"gocloud.dev/server/sdserver"
	"gocloud.dev/server/xrayserver"
	"google.golang.org/genproto/googleapis/cloud/runtimeconfig/v1beta1"
	"net/http"
)

// Injectors from inject_aws.go:

func setupAWS(ctx context.Context, flags *cliFlags) (*server.Server, func(), error) {
	client := _wireClientValue
	certFetcher := &rds.CertFetcher{
		Client: client,
	}
	params := awsSQLParams(flags)
	db, cleanup, err := rdsmysql.Open(ctx, certFetcher, params)
	if err != nil {
		return nil, nil, err
	}
	session, err := aws.NewDefaultSession()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	bucket, cleanup2, err := awsBucket(ctx, session, flags)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	variable, err := awsMOTDVar(ctx, session, flags)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	mainApplication := newApplication(db, bucket, variable)
	router := newRouter(mainApplication)
	ncsaLogger := xrayserver.NewRequestLogger()
	v, cleanup3 := appHealthChecks(db)
	xRay := xrayserver.NewXRayClient(session)
	exporter, cleanup4, err := xrayserver.NewExporter(xRay)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	sampler := trace.AlwaysSample()
	defaultDriver := _wireDefaultDriverValue
	options := &server.Options{
		RequestLogger:         ncsaLogger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		DefaultSamplingPolicy: sampler,
		Driver:                defaultDriver,
	}
	serverServer := server.New(router, options)
	return serverServer, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireClientValue        = http.DefaultClient
	_wireDefaultDriverValue = &server.DefaultDriver{}
)

// Injectors from inject_azure.go:

func setupAzure(ctx context.Context, flags *cliFlags) (*server.Server, func(), error) {
	db, err := dialLocalSQL(flags)
	if err != nil {
		return nil, nil, err
	}
	accountName, err := azureblob.DefaultAccountName()
	if err != nil {
		return nil, nil, err
	}
	accountKey, err := azureblob.DefaultAccountKey()
	if err != nil {
		return nil, nil, err
	}
	sharedKeyCredential, err := azureblob.NewCredential(accountName, accountKey)
	if err != nil {
		return nil, nil, err
	}
	pipelineOptions := _wirePipelineOptionsValue
	pipeline := azureblob.NewPipeline(sharedKeyCredential, pipelineOptions)
	bucket, cleanup, err := azureBucket(ctx, pipeline, accountName, flags)
	if err != nil {
		return nil, nil, err
	}
	variable, err := azureMOTDVar(ctx, bucket, flags)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	mainApplication := newApplication(db, bucket, variable)
	router := newRouter(mainApplication)
	logger := _wireLoggerValue
	v, cleanup2 := appHealthChecks(db)
	exporter := _wireExporterValue
	sampler := trace.AlwaysSample()
	defaultDriver := _wireDefaultDriverValue
	options := &server.Options{
		RequestLogger:         logger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		DefaultSamplingPolicy: sampler,
		Driver:                defaultDriver,
	}
	serverServer := server.New(router, options)
	return serverServer, func() {
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wirePipelineOptionsValue = azblob.PipelineOptions{}
	_wireLoggerValue          = requestlog.Logger(nil)
	_wireExporterValue        = trace.Exporter(nil)
)

// Injectors from inject_gcp.go:

func setupGCP(ctx context.Context, flags *cliFlags) (*server.Server, func(), error) {
	roundTripper := gcp.DefaultTransport()
	credentials, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, nil, err
	}
	tokenSource := gcp.CredentialsTokenSource(credentials)
	httpClient, err := gcp.NewHTTPClient(roundTripper, tokenSource)
	if err != nil {
		return nil, nil, err
	}
	remoteCertSource := cloudsql.NewCertSource(httpClient)
	projectID, err := gcp.DefaultProjectID(credentials)
	if err != nil {
		return nil, nil, err
	}
	params := gcpSQLParams(projectID, flags)
	db, cleanup, err := cloudmysql.Open(ctx, remoteCertSource, params)
	if err != nil {
		return nil, nil, err
	}
	bucket, cleanup2, err := gcpBucket(ctx, flags, httpClient)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	runtimeConfigManagerClient, cleanup3, err := gcpruntimeconfig.Dial(ctx, tokenSource)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	variable, cleanup4, err := gcpMOTDVar(ctx, runtimeConfigManagerClient, projectID, flags)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	mainApplication := newApplication(db, bucket, variable)
	router := newRouter(mainApplication)
	stackdriverLogger := sdserver.NewRequestLogger()
	v, cleanup5 := appHealthChecks(db)
	monitoredresourceInterface := monitoredresource.Autodetect()
	exporter, cleanup6, err := sdserver.NewExporter(projectID, tokenSource, monitoredresourceInterface)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	sampler := trace.AlwaysSample()
	defaultDriver := _wireDefaultDriverValue
	options := &server.Options{
		RequestLogger:         stackdriverLogger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		DefaultSamplingPolicy: sampler,
		Driver:                defaultDriver,
	}
	serverServer := server.New(router, options)
	return serverServer, func() {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// Injectors from inject_local.go:

func setupLocal(ctx context.Context, flags *cliFlags) (*server.Server, func(), error) {
	db, err := dialLocalSQL(flags)
	if err != nil {
		return nil, nil, err
	}
	bucket, err := localBucket(flags)
	if err != nil {
		return nil, nil, err
	}
	variable, cleanup, err := localRuntimeVar(flags)
	if err != nil {
		return nil, nil, err
	}
	mainApplication := newApplication(db, bucket, variable)
	router := newRouter(mainApplication)
	logger := _wireRequestlogLoggerValue
	v, cleanup2 := appHealthChecks(db)
	exporter := _wireTraceExporterValue
	sampler := trace.AlwaysSample()
	defaultDriver := _wireDefaultDriverValue
	options := &server.Options{
		RequestLogger:         logger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		DefaultSamplingPolicy: sampler,
		Driver:                defaultDriver,
	}
	serverServer := server.New(router, options)
	return serverServer, func() {
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireRequestlogLoggerValue = requestlog.Logger(nil)
	_wireTraceExporterValue    = trace.Exporter(nil)
)

// inject_aws.go:

// awsBucket is a Wire provider function that returns the S3 bucket based on the
// command-line flags.
func awsBucket(ctx context.Context, cp client.ConfigProvider, flags *cliFlags) (*blob.Bucket, func(), error) {
	b, err := s3blob.OpenBucket(ctx, cp, flags.bucket, nil)
	if err != nil {
		return nil, nil, err
	}
	return b, func() { b.Close() }, nil
}

// awsSQLParams is a Wire provider function that returns the RDS SQL connection
// parameters based on the command-line flags. Other providers inside
// awscloud.AWS use the parameters to construct a *sql.DB.
func awsSQLParams(flags *cliFlags) *rdsmysql.Params {
	return &rdsmysql.Params{
		Endpoint: flags.dbHost,
		Database: flags.dbName,
		User:     flags.dbUser,
		Password: flags.dbPassword,
	}
}

// awsMOTDVar is a Wire provider function that returns the Message of the Day
// variable from SSM Parameter Store.
func awsMOTDVar(ctx context.Context, sess client.ConfigProvider, flags *cliFlags) (*runtimevar.Variable, error) {
	return awsparamstore.OpenVariable(sess, flags.motdVar, runtimevar.StringDecoder, &awsparamstore.Options{
		WaitDuration: flags.motdVarWaitTime,
	})
}

// inject_azure.go:

// azureBucket is a Wire provider function that returns the Azure bucket based
// on the command-line flags.
func azureBucket(ctx context.Context, p pipeline.Pipeline, accountName azureblob.AccountName, flags *cliFlags) (*blob.Bucket, func(), error) {
	b, err := azureblob.OpenBucket(ctx, p, accountName, flags.bucket, nil)
	if err != nil {
		return nil, nil, err
	}
	return b, func() { b.Close() }, nil
}

// azureMOTDVar is a Wire provider function that returns the Message of the Day
// variable read from a blob stored in Azure.
func azureMOTDVar(ctx context.Context, b *blob.Bucket, flags *cliFlags) (*runtimevar.Variable, error) {
	return blobvar.OpenVariable(b, flags.motdVar, runtimevar.StringDecoder, &blobvar.Options{
		WaitDuration: flags.motdVarWaitTime,
	})
}

// inject_gcp.go:

// gcpBucket is a Wire provider function that returns the GCS bucket based on
// the command-line flags.
func gcpBucket(ctx context.Context, flags *cliFlags, client2 *gcp.HTTPClient) (*blob.Bucket, func(), error) {
	b, err := gcsblob.OpenBucket(ctx, client2, flags.bucket, nil)
	if err != nil {
		return nil, nil, err
	}
	return b, func() { b.Close() }, nil
}

// gcpSQLParams is a Wire provider function that returns the Cloud SQL
// connection parameters based on the command-line flags. Other providers inside
// gcpcloud.GCP use the parameters to construct a *sql.DB.
func gcpSQLParams(id gcp.ProjectID, flags *cliFlags) *cloudmysql.Params {
	return &cloudmysql.Params{
		ProjectID: string(id),
		Region:    flags.cloudSQLRegion,
		Instance:  flags.dbHost,
		Database:  flags.dbName,
		User:      flags.dbUser,
		Password:  flags.dbPassword,
	}
}

// gcpMOTDVar is a Wire provider function that returns the Message of the Day
// variable from Runtime Configurator.
func gcpMOTDVar(ctx context.Context, client2 runtimeconfig.RuntimeConfigManagerClient, project gcp.ProjectID, flags *cliFlags) (*runtimevar.Variable, func(), error) {
	variableKey := gcpruntimeconfig.VariableKey(project, flags.runtimeConfigName, flags.motdVar)
	v, err := gcpruntimeconfig.OpenVariable(client2, variableKey, runtimevar.StringDecoder, &gcpruntimeconfig.Options{
		WaitDuration: flags.motdVarWaitTime,
	})
	if err != nil {
		return nil, nil, err
	}
	return v, func() { v.Close() }, nil
}

// inject_local.go:

// localBucket is a Wire provider function that returns a directory-based bucket
// based on the command-line flags.
func localBucket(flags *cliFlags) (*blob.Bucket, error) {
	return fileblob.OpenBucket(flags.bucket, nil)
}

// dialLocalSQL is a Wire provider function that connects to a MySQL database
// (usually on localhost).
func dialLocalSQL(flags *cliFlags) (*sql.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 flags.dbHost,
		DBName:               flags.dbName,
		User:                 flags.dbUser,
		Passwd:               flags.dbPassword,
		AllowNativePasswords: true,
	}
	return sql.Open("mysql", cfg.FormatDSN())
}

// localRuntimeVar is a Wire provider function that returns the Message of the
// Day variable based on a local file.
func localRuntimeVar(flags *cliFlags) (*runtimevar.Variable, func(), error) {
	v, err := filevar.OpenVariable(flags.motdVar, runtimevar.StringDecoder, &filevar.Options{
		WaitDuration: flags.motdVarWaitTime,
	})
	if err != nil {
		return nil, nil, err
	}
	return v, func() { v.Close() }, nil
}
