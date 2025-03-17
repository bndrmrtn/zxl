package modules

import (
	"github.com/bndrmrtn/zxl/internal/modules/sqlmodule"
	"github.com/bndrmrtn/zxl/internal/modules/thread"
	"github.com/bndrmrtn/zxl/internal/modules/zruntime"
	"github.com/bndrmrtn/zxl/internal/modules/ztime"
	"github.com/bndrmrtn/zxl/lang"
)

// Get returns a list of all available modules.
func Get() []lang.Module {
	return []lang.Module{
		NewRandModule(),
		NewIOModule(),
		NewHttpModule(),
		NewJSONModule(),
		NewEnv(),
		NewConvert(),
		sqlmodule.New(),
		zruntime.New(),
		thread.New(),
		ztime.New(),
	}
}
