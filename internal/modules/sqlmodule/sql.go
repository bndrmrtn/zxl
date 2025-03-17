package sqlmodule

import (
	"github.com/bndrmrtn/zxl/lang"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type SQL struct{}

func New() *SQL {
	return &SQL{}
}

func (*SQL) Namespace() string {
	return "sql"
}

func (s *SQL) Objects() map[string]lang.Object {
	return nil
}

func (s *SQL) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"open": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			return NewDB(args[0].Value().(string), args[1].Value().(string))
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "driver", Type: lang.TString}, lang.TypeSafeArg{Name: "dsn", Type: lang.TString}),
	}
}
