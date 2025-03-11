package modules

import (
	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/modules/sqlmodule"
	"github.com/bndrmrtn/zxl/internal/modules/zruntime"
)

// Get returns a list of all available modules.
func Get() []lang.Module {
	return []lang.Module{
		NewRandModule(),
		NewIOModule(),
		NewHttpModule(),
		NewJSONModule(),
		NewThreadModule(),
		sqlmodule.New(),
		zruntime.New(),
	}
}
