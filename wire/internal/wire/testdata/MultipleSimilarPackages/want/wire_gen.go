// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	bar "example.com/bar"
	baz "example.com/baz"
	foo "example.com/foo"
	fmt "fmt"
)

// Injectors from wire.go:

func newMainService(config *foo.Config, config2 *bar.Config, config3 *baz.Config) *MainService {
	service := foo.New(config)
	service2 := bar.New(config2, service)
	service3 := baz.New(config3, service2)
	mainService := &MainService{
		Foo: service,
		Bar: service2,
		Baz: service3,
	}
	return mainService
}

// wire.go:

type MainConfig struct {
	Foo *foo.Config
	Bar *bar.Config
	Baz *baz.Config
}

type MainService struct {
	Foo *foo.Service
	Bar *bar.Service
	Baz *baz.Service
}

func (m *MainService) String() string {
	return fmt.Sprintf("%d %d %d", m.Foo.Cfg.V, m.Bar.Cfg.V, m.Baz.Cfg.V)
}

func main() {
	cfg := &MainConfig{
		Foo: &foo.Config{1},
		Bar: &bar.Config{2},
		Baz: &baz.Config{3},
	}
	svc := newMainService(cfg.Foo, cfg.Bar, cfg.Baz)
	fmt.Println(svc.String())
}
