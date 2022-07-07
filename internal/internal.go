package internal

import (
	"github.com/kiru997/go-ex/internal/controller"
	"github.com/kiru997/go-ex/internal/service/demo"
	"go.uber.org/fx"
)

var (
	Module = fx.Options(
		ModuleController,
		ModuleService,
	)

	ModuleController = fx.Options(
		fx.Invoke(
			controller.RegisterDemoController,
		),
	)

	ModuleService = fx.Options(
		fx.Provide(
			demo.NewService,
		),
	)
)
