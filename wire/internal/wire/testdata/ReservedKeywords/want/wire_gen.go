// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

// Injectors from wire.go:

func injectInterface() Interface {
	select2 := provideSelect()
	mainInterface := provideInterface(select2)
	return mainInterface
}

// wire.go:

var mainSelect = 0
