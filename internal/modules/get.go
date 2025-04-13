package modules

import (
	"github.com/bndrmrtn/flare/internal/modules/sqlmodule"
	"github.com/bndrmrtn/flare/internal/modules/thread"
	"github.com/bndrmrtn/flare/internal/modules/zruntime"
	"github.com/bndrmrtn/flare/internal/modules/ztime"
	"github.com/bndrmrtn/flare/lang"
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
		NewCrypto(),
		sqlmodule.New(),
		zruntime.New(),
		thread.New(),
		ztime.New(),
	}
}
