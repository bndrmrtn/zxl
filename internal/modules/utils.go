package modules

import "github.com/bndrmrtn/zxl/internal/lang"

func immute(obj lang.Object) lang.Object {
	obj.Immute()
	return obj
}
