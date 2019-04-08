// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/wire"
	"go.opencensus.io/trace"
	"gocloud.dev/health"
	"gocloud.dev/server"
	"gocloud.dev/server/xrayserver"
)

// Injectors from inject.go:

func initialize(ctx context.Context, cfg *appConfig) (*server.Server, func(), error) {
	serveMux := NewRouter()
	ncsaLogger := xrayserver.NewRequestLogger()
	v := _wireValue
	options := awsOptions(cfg)
	sessionSession, err := session.NewSessionWithOptions(options)
	if err != nil {
		return nil, nil, err
	}
	xRay := xrayserver.NewXRayClient(sessionSession)
	exporter, cleanup, err := xrayserver.NewExporter(xRay)
	if err != nil {
		return nil, nil, err
	}
	sampler := trace.AlwaysSample()
	defaultDriver := _wireDefaultDriverValue
	serverOptions := &server.Options{
		RequestLogger:         ncsaLogger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		DefaultSamplingPolicy: sampler,
		Driver:                defaultDriver,
	}
	serverServer := server.New(serveMux, serverOptions)
	return serverServer, func() {
		cleanup()
	}, nil
}

var (
	_wireValue              = []health.Checker{connection}
	_wireDefaultDriverValue = &server.DefaultDriver{}
)

// inject.go:

var awsSession = wire.NewSet(session.NewSessionWithOptions, awsOptions, wire.Bind((*client.ConfigProvider)(nil), (*session.Session)(nil)), configConfidentials)

func awsOptions(cfg *appConfig) session.Options {
	return session.Options{
		Config: aws.Config{
			Region: aws.String(cfg.region),
		},
	}
}

func configConfidentials(cfg *aws.Config) *credentials.Credentials {
	return cfg.Credentials
}
