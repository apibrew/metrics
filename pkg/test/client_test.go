package test

import (
	"github.com/apibrew/apibrew/pkg/api"
	"github.com/apibrew/apibrew/pkg/service"
	"github.com/apibrew/apibrew/pkg/service/impl"
	"github.com/apibrew/apibrew/pkg/test/setup"
	"github.com/apibrew/template/pkg"
)

var container service.Container

var apiInterface api.Interface

func init() {
	container = setup.GetContainer()

	app := container.(*impl.App)

	app.RegisterModule(pkg.NewModule)

	apiInterface = api.NewInterface(container)

}
