package module

import (
	"go.uber.org/fx"

	"ecosystem.garyle/service/internal/app/module/ota"
	"ecosystem.garyle/service/internal/app/module/wms"
)

// Module combines all application modules
var Module = fx.Module("app",
	ota.Module,
	wms.Module,
)
