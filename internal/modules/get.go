package modules

import (
	"github.com/flarelang/flare/internal/modules/sqlmodule"
	"github.com/flarelang/flare/internal/modules/thread"
	"github.com/flarelang/flare/internal/modules/zruntime"
	"github.com/flarelang/flare/internal/modules/ztime"
	"github.com/flarelang/flare/lang"
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
